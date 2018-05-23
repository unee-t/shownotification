package main

import (
	"fmt"
	"strings"
	"html/template"
	"net/http"
	"net/http/httputil"
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
	http.Handle("/hook-neighbor", http.HandlerFunc(h.hook_neighbor))

	addr := ":" + os.Getenv("PORT")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func (h handler) hook_neighbor(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	h.es.SendEventMessage(string(dump), "", "")
}

func (h handler) hook(w http.ResponseWriter, r *http.Request) {
	_, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
        // TODO cahce strings for performance
        // TODO use etcd service discovery

        http.Get("http://0.0.0.0:"+os.Getenv("PORT")+"/hook-neighbor")
        for _, neighbor := range strings.Split( os.Getenv("NEIGHBORS"), "," ) {
            http.Get("http://"+neighbor+"/hook-neighbor")
        }
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Robots-Tag", "none")
	t := template.Must(template.New("").ParseGlob("public/*.html"))
	t.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"COMMIT": os.Getenv("COMMIT"),
	})
}
