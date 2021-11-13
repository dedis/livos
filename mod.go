package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dedis/livos/storage/bbolt"
	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"github.com/dedis/livos/web/controller"
)

type key int

//go:embed web/index.html
var content embed.FS

//go:embed web/homepage.html
var contenthomepage embed.FS

//go:embed web/static
var static embed.FS

//go:embed web/views
var views embed.FS

//go:embed web/images
var image embed.FS

const (
	requestIDKey key = 0
)

func main() {

	listenAddr := ":9000"
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:  tracing(nextRequestID)(logging(logger)(mux)),
		ErrorLog: logger,
	}

	//creation of the database
	db, err := bbolt.New("database.db")
	if err != nil {
		logger.Fatalf("failed to create the database : %v", err)
		return
	}

	//creation of voting system
	vil := make(map[string]*impl.VotingInstance)
	votingSystem := impl.NewVotingSystem(db, vil)

	//creation of controller (for the web interactions)
	ctrl := controller.NewController(content, contenthomepage, views, votingSystem)

	//creation of users
	userNoemien, err := votingSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), voting.Choice{})
	if err != nil {
		logger.Fatal("user noemien creation is incorrect.")
	}
	userGuillaume, err := votingSystem.NewUser("Guillaume", make(map[string]voting.Liquid), make(map[string]voting.Liquid), voting.Choice{})
	if err != nil {
		logger.Fatal("user guillaume creation is incorrect.")
	}
	userEtienne, err := votingSystem.NewUser("Etienne", make(map[string]voting.Liquid), make(map[string]voting.Liquid), voting.Choice{})
	if err != nil {
		logger.Fatal("user etienne creation is incorrect.")
	}

	//list of voters
	voters := []*voting.User{&userNoemien, &userGuillaume, &userEtienne}

	description := "Do you want fries every day at the restaurant?"
	description2 := "Do you want free vacations (365 days a year) ?"

	candidats := make([]string, 3)
	votingConfig, err := impl.NewVotingConfig(voters, "VoteRoom1", description, candidats)
	if err != nil {
		logger.Fatal("NewVotingConfig is incorrect")
	}
	voters2 := []*voting.User{}
	votingConfig2, err := impl.NewVotingConfig(voters2, "VoteRoom2", description2, candidats)
	if err != nil {
		logger.Fatal("NewVotingConfig is incorrect")
	}
	votes := make(map[string]voting.Choice)
	votes2 := make(map[string]voting.Choice)
	votingSystem.CreateAndAdd("001", votingConfig, "open", votes)
	votingSystem.CreateAndAdd("002", votingConfig2, "close", votes2)

	fmt.Println("VOTING INSTANCE LIST : ", votingSystem.VotingInstancesList)

	var vi = *votingSystem.VotingInstancesList["001"]

	fmt.Println("VI:", vi)
	yesChoice := make(map[string]voting.Liquid)
	noChoice := make(map[string]voting.Liquid)
	midChoice := make(map[string]voting.Liquid)

	liq100, err100 := impl.NewLiquid(100)
	liq50, err50 := impl.NewLiquid(50)
	liqid0, err0 := impl.NewLiquid(0)
	if (err100 != nil) || (err50 != nil) || (err0 != nil) {
		logger.Fatalf("Creation of liquid is incorrect.")
	}

	yesChoice["yes"] = liq100
	yesChoice["no"] = liqid0
	noChoice["no"] = liq100
	noChoice["yes"] = liqid0
	midChoice["no"] = liq50
	midChoice["yes"] = liq50
	choiceGuillaume, errG := impl.NewChoice(noChoice)
	choiceEtienne, errE := impl.NewChoice(midChoice)
	fmt.Println("CHOICE Guigui: ", choiceGuillaume)
	fmt.Println("CHOICE etien: ", choiceEtienne)
	if (errG != nil) || (errE != nil) {
		logger.Fatalf("Choices creation incorrect.")
	}

	vi.SetChoice(&userGuillaume, choiceGuillaume)
	fmt.Println(":::::: Result of the setchoice of guillaume", userGuillaume.MyChoice)
	vi.SetChoice(&userEtienne, choiceEtienne)
	fmt.Println(":::::: Result of the setchoice of etienne", userEtienne.MyChoice)

	vi.CastVote(&userGuillaume)
	vi.CastVote(&userEtienne)

	fmt.Println("RESULTS OF THE VOTE ====> ", vi.GetResults())

	mux.HandleFunc("/", ctrl.HandleHome)
	mux.HandleFunc("/homepage", ctrl.HandleHomePage)
	mux.HandleFunc("/election", ctrl.HandleShowElection)
	mux.HandleFunc("/results", ctrl.HandleShowResults)
	mux.HandleFunc("/manage", ctrl.HandleManageVoting)

	// serve assets
	mux.Handle("/web/static/", http.FileServer(http.FS(static)))
	mux.Handle("/web/images/", http.FileServer(http.FS(image)))

	// create connection
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("failed to create conn '%s': %v", listenAddr, err)
		return
	}

	// launch server
	err = server.Serve(ln)
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

// Utility function for logging

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
