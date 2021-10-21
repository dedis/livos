package controller

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	//"github.com/dedis/livos/storage/bbolt"
)

// NewController ...
func NewController(homeHTML embed.FS) Controller {
	return Controller{
		homeHTML: homeHTML,
	}
}

// Controller ...
type Controller struct {
	homeHTML embed.FS
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
