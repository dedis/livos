package controller

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"

	//"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting/impl"
)

// NewController ...
func NewController(homeHTML embed.FS, vs impl.VotingSystem) Controller {
	return Controller{
		homeHTML: homeHTML,
		vs:       vs,
	}
}

// Controller ...
type Controller struct {
	homeHTML embed.FS
	vs       impl.VotingSystem
}

// HandleHome ...
func (c Controller) HandleHome(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(c.homeHTML, "index.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//req.Form.Get(election//)

	name := req.FormValue("username")
	description := req.FormValue("description")
	roomID := req.FormValue("roomID")

	//Only print. Have to be stored on database
	fmt.Fprintln(w, "Username = \n", name)
	fmt.Fprintln(w, "Description = \n", description)
	fmt.Fprintln(w, "RoomID = \n", roomID)
}

func (c Controller) HandleHomePage(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(c.homeHTML, "homepage.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c Controller) HandleQuit(w http.ResponseWriter, req *http.Request) {
	//action
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	<-stop

	fmt.Fprintf(w, "shutting down ...\n")

	//fmt.Fprintf(w, "Server is shut down")
	//ctx, cancel := context.WithCancel(context.Background())
	//if cn, ok := w.(http.CloseNotifier); ok {
	//go func(done <-chan struct{}, closed <-chan bool) {
	//	cancel()
	//}(ctx.Done(), cn.CloseNotify())
	//}

}
