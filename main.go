package main

import (
	"net/http"
	"os"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"gopkg.in/antage/eventsource.v1"
)

type handler struct {
	es eventsource.EventSource
}

func init() {
	if os.Getenv("UP_STAGE") == "" {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(jsonhandler.Default)
	}
}

func main() {
	es := eventsource.New(nil, nil)
	defer es.Close()
	h := handler{es: es}

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/events", es)
	http.Handle("/hook", http.HandlerFunc(h.hook))

	addr := ":" + os.Getenv("PORT")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func (h handler) hook(w http.ResponseWriter, r *http.Request) {
	h.es.SendEventMessage("hello", "", "")

}
