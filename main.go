package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"gopkg.in/antage/eventsource.v1"
)

type handler struct {
	es eventsource.EventSource
}

type sqsmessage struct {
	Message string
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

func (h handler) receiveSQS() {

	cfg, err := external.LoadDefaultAWSConfig(external.WithSharedConfigProfile("uneet-dev"))
	if err != nil {
		log.WithError(err).Error("failed to load config")
		return
	}

	// https://godoc.org/github.com/aws/aws-sdk-go-v2/service/sqs#ReceiveMessageRequest
	sqssvc := sqs.New(cfg)

	msgReq := sqssvc.ReceiveMessageRequest(&sqs.ReceiveMessageInput{
		QueueUrl: aws.String("https://sqs.ap-southeast-1.amazonaws.com/812644853088/atest"),
	})

	log.Info("Requesting from SQS")
	msgs, err := msgReq.Send()
	if err != nil {
		panic(err)
	}
	var sqsm sqsmessage
	for _, msg := range msgs.Messages {
		log.Infof("SQS Receive Message %s", &msg.Body)
		json.Unmarshal([]byte(*msg.Body), &sqsm)
		log.Infof("here %s", sqsm.Message)
	}

	h.es.SendEventMessage(sqsm.Message, "", "")
}

func (h handler) hook(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	h.es.SendEventMessage(string(dump), "", "")
	h.receiveSQS()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Robots-Tag", "none")
	t := template.Must(template.New("").ParseGlob("public/*.html"))
	t.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"COMMIT": os.Getenv("COMMIT"),
	})
}
