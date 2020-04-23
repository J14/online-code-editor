package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"html/template"
	"fmt"
	"bytes"
	"./run"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/run", runCode)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("index").Parse(string(file))

	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func runCode(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(body)
	code := buf.String()

	result := run.Run(code)

	fmt.Fprintf(w, result)
}