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

//go:embed web/homepage/001.html
var content001 embed.FS

//go:embed web/static
var static embed.FS

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
	vil := make(map[string]impl.VotingInstance)
	votingSystem := impl.NewVotingSystem(db, vil)

	//creation of controller (for the web interactions)
	ctrl := controller.NewController(content, contenthomepage, content001, votingSystem)

	voters := make([]string, 3)
	voters = append(voters, "Noemien", "Guillaume", "Etienne")
	//fmt.Println(voters)

	title := "VoteRoom1"
	description := "Do you want fries every day at the restaurant?"

	candidats := make([]string, 3)
	votingConfig := impl.NewVotingConfig(voters, title, description, candidats)
	votes := make(map[string]voting.Choice)
	votingSystem.Create("001", votingConfig, "open", votes)
	votingSystem.Create("002", votingConfig, "close", votes)

	//fmt.Println("Test de listVoting", votingSystem.CastVote())
	//fmt.Println("VOTING INSTANCE LIST : ", votingSystem.VotingInstancesList)

	var vi = votingSystem.VotingInstancesList["001"]

	deleg := make(map[string]voting.Liquid)
	yesChoice := make(map[string]voting.Liquid)
	noChoice := make(map[string]voting.Liquid)
	midChoice := make(map[string]voting.Liquid)

	liq100 := impl.NewLiquid(100)
	liq50 := impl.NewLiquid(50)
	liqid0 := impl.NewLiquid(0)

	yesChoice["yes"] = liq100
	yesChoice["no"] = liqid0
	noChoice["no"] = liq100
	noChoice["yes"] = liqid0
	midChoice["no"] = liq50
	midChoice["yes"] = liq50
	choiceNoemien := impl.NewChoice(deleg, yesChoice, 0)
	choiceGuillaume := impl.NewChoice(deleg, noChoice, 0)
	choiceEtienne := impl.NewChoice(deleg, midChoice, 0)
	vi.CastVote("Noemien", choiceNoemien)
	vi.CastVote("Guillaume", choiceGuillaume)
	vi.CastVote("Etienne", choiceEtienne)

	fmt.Println("RESULTS OF THE VOTE ====> ", vi.GetResults())

	mux.HandleFunc("/", ctrl.HandleHome)
	mux.HandleFunc("/homepage", ctrl.HandleHomePage)
	mux.HandleFunc("/homepage/001", ctrl.Handle001)

	// serve assets
	mux.Handle("/static/", http.FileServer(http.FS(static)))

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
