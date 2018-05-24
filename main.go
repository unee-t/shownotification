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

var queueurl = aws.String("https://sqs.ap-southeast-1.amazonaws.com/812644853088/atest")

// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-long-polling.html
var longpollDuration = int64(20)

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

	go func() {
		for {
			if es.ConsumersCount() > 0 { // Should only run if there is an active client
				h.receiveSQS()
			}
			log.Infof("Long polling %ds", longpollDuration)
		}
	}()

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
		WaitTimeSeconds: aws.Int64(longpollDuration), // Long poll for 10
		QueueUrl:        queueurl,
	})

	log.Info("Requesting from SQS")
	msgs, err := msgReq.Send() // I kind of expected this to block until a message appeared
	if err != nil {
		panic(err)
	}
	var sqsm sqsmessage
	for _, msg := range msgs.Messages {
		log.Infof("Payload %v", msg)
		json.Unmarshal([]byte(*msg.Body), &sqsm)

		// in the Unee-T platform, only the consumer that processes the relevant message should delete it
		delReq := sqssvc.DeleteMessageRequest(&sqs.DeleteMessageInput{
			QueueUrl:      queueurl,
			ReceiptHandle: msg.ReceiptHandle,
		})

		_, err := delReq.Send()
		if err != nil {
			panic(err)
		}
		log.Infof("%v deleted", *msg.MessageId)

	}

	h.es.SendEventMessage(sqsm.Message, "", "")
}

func (h handler) hook(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	h.es.SendEventMessage(string(dump), "", "")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Robots-Tag", "none")
	t := template.Must(template.New("").ParseGlob("public/*.html"))
	t.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"COMMIT": os.Getenv("COMMIT"),
	})
}
