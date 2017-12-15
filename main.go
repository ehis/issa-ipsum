package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Message envelope for http responses
type Message struct {
	Data interface{} `json:"data"`
	Time time.Time   `json:"time"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/issa-ipsum", IssaIpsumHandler).Methods("GET")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		fmt.Println("====== Routes ======")
		t, err := route.GetPathTemplate()

		if err != nil {
			return err
		}

		fmt.Println(t)

		return nil
	})
	http.Handle("/", r)

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// IssaIpsumHandler generates a ipsum from 21 Savage's Issa Album
func IssaIpsumHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("sentences")
	sentences, err := strconv.Atoi(query)

	if err != nil {
		// default to 5 sentences
		sentences = 5
	}

	text := IssaMarkovChain(sentences)

	msg := Message{
		Data: string(text),
		Time: time.Now(),
	}

	resp, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// IssaMarkovChain takes ideas from markov chain algorithm
func IssaMarkovChain(sentences int) []byte {
	file, err := os.Open("./issa-album/bank-account.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var buffer bytes.Buffer
	for i := 0; i < sentences; i++ {
		success := scanner.Scan()

		if success == false {
			err = scanner.Err()
			if err == nil {
				log.Println("Scan completed and reached EOF")
			} else {
				log.Fatal(err)
			}
		}

		buffer.WriteString(fmt.Sprintf("Issa Ipsum %s ", scanner.Text()))
	}

	return buffer.Bytes()
}
