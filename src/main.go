package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os/exec"
)

var (
	laddr = flag.String("l", ":8080", "HTTP service address (e.g.address, ':8080')")
)

type Record interface {
	GetGitURL() string
	GetName() string
}

type BitbucketRecord struct {
	Repository struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
		URL  string `json:"absolute_url"`
	} `json:"repository"`
	BaseURL string `json:"canon_url"`
	User    string `json:"user"`
}

func (r BitbucketRecord) GetGitURL() string {
	return fmt.Sprintf("%s%s", r.BaseURL, r.Repository.URL)
}

func (r BitbucketRecord) GetName() string {
	return r.Repository.Name
}

type GithubRecord struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"git_url"`
	} `json:"repository"`
}

func (r GithubRecord) GetGitURL() string {
	return r.Repository.URL
}

func (r GithubRecord) GetName() string {
	return r.Repository.Name
}

func RecordFactory(hookname string) (Record, error) {
	switch hookname {
	case "bitbucket":
		return new(BitbucketRecord), nil
	case "github":
		return new(GithubRecord), nil
	default:
		return nil, errors.New("Unknown hookname.")
	}
}

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

func RunScript(w http.ResponseWriter, hook string, action string, params ...string) {
	fw := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}
	scriptname := fmt.Sprintf("./scripts/%s/%s.sh", hook, action)
	log.Println("Exec script: ", scriptname)
	cmd := exec.Command(scriptname, params...)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	cmd.Run()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hookname := params["hookname"]
	action := params["action"]

	log.Println("Hook name: ", hookname)

	var record, err = RecordFactory(hookname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	RunScript(w, hookname, action, record.GetGitURL(), record.GetName())
}

func main() {
	flag.Parse()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{hookname:[a-z]+}/{action:[a-z]+}", Handler).Methods("POST")

	http.Handle("/", rtr)

	log.Println("webhookd server listening...")
	log.Fatal(http.ListenAndServe(*laddr, nil))
}
