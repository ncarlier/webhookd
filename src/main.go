package main

import (
	"./hooks"
	"encoding/json"
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

	var record, err = hooks.RecordFactory(hookname)
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
