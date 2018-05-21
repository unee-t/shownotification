package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/apex/log"
	"gopkg.in/antage/eventsource.v1"
)

type handler struct {
	es eventsource.EventSource
}

func main() {
	es := eventsource.New(nil, nil)
	defer es.Close()
	h := handler{es: es}

	http.Handle("/", http.HandlerFunc(handleIndex))
	http.Handle("/events", es)
	http.Handle("/hook", http.HandlerFunc(h.hook))

	addr := ":" + os.Getenv("PORT")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func (h handler) hook(w http.ResponseWriter, r *http.Request) {
	h.es.SendEventMessage(r.UserAgent(), "", "")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Robots-Tag", "none")
	t := template.Must(template.New("").ParseGlob("public/*.html"))
	t.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"COMMIT": os.Getenv("COMMIT"),
	})
}
