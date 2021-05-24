package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render
var messageChannels = make(map[chan []byte]bool)

type abc struct{}

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "./templates/*.html",
	}

	rnd = renderer.New(opts)
}

var (
	broker = "localhost:29092"
	group  = "G1"
	topics = []string{"map_stream"}
	c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
		"enable.auto.commit":              true,
		"auto.offset.reset":               "earliest"})
	sigchan = make(chan os.Signal, 1)
)

func main() {

	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/topic/TOPICNAME", topic)
	port := ":9011"
	log.Println("Listening on port ", port)
	http.ListenAndServe(port, mux)

}

func home(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, abc{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func about(w http.ResponseWriter, r *http.Request) {
	jsonStructure, _ := json.Marshal(map[string]interface{}{
		"name":    "Susheel",
		"message": "hi",
		"number":  100,
	})

	go func() {
		for messageChannel := range messageChannels {
			messageChannel <- []byte(jsonStructure)
		}
	}()

	w.Write([]byte("ok."))
}

func formatSSE(event string, data string) []byte {
	// fmt.Println("Event::, ", event)
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}

func topic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	run := true

	n := 0
	t1 := time.Now()

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				w.Write(formatSSE("message", string(e.Value)))
				w.(http.Flusher).Flush()
				n++
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
		if time.Now().Sub(t1) > 1*time.Second {
			fmt.Printf("n:%d\n", n)
			n = 0
			t1 = time.Now()
		}

	}

	fmt.Printf("Closing consumer\n")
	c.Close()

}
