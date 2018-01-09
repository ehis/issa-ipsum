package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var trackList = []string{
	"7-min-freestyle",
	"baby-girl",
	"bad-business",
	"bank-account",
	"close-my-eyes",
	"dead-people",
	"facetime",
	"famous",
	"money-convo",
	"nothing-new",
	"numb",
	"special",
	"thug-life",
	"whole-lot",
}

// Message envelope for http responses
type Message struct {
	Body interface{} `json:"body"`
	Time time.Time   `json:"time"`
}

func extractSentences(r *http.Request) int {
	query := r.FormValue("sentences")
	sentences, err := strconv.Atoi(query)

	if err != nil {
		// default to 5 sentences
		sentences = 5
	}

	return sentences
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
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
	sentences := extractSentences(r)
	text := IssaMarkovChain(sentences)

	msg := Message{
		Body: string(text),
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
	rand.Seed(time.Now().Unix())
	track := trackList[rand.Intn(len(trackList))]
	filePath := fmt.Sprintf("./issa-album/%s.txt", track)
	file, err := os.Open(filePath)

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
				return IssaMarkovChain(i)
			}

			log.Fatal(err)
		}

		buffer.WriteString(fmt.Sprintf("Issa Ipsum %s ", scanner.Text()))
	}

	return buffer.Bytes()
}

// IndexHandler loads the index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sentences := extractSentences(r)

	title := "Issa Ipsum"

	page := map[string]interface{}{
		"Title":     title,
		"Sentences": sentences,
		"Body":      string(IssaMarkovChain(sentences)),
	}

	t, err := template.ParseFiles("pages/index.html")

	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, page)
}
