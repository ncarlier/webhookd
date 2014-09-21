package main

import (
	"./hooks"
	"./notifications"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
)

var (
	laddr = flag.String("l", ":8080", "HTTP service address (e.g.address, ':8080')")
)

type HookContext struct {
	Hook   string
	Action string
	args   []string
}

func Notify(text string, context *HookContext) {
	var subject = fmt.Sprintf("Action %s executed.", context.Action)
	var notifier, err = notifications.NotifierFactory()
	if err != nil {
		log.Println(err)
		return
	}
	if notifier == nil {
		log.Println("Notification provider not found.")
		return
	}
	notifier.Notify(text, subject)
}

func RunScript(w http.ResponseWriter, context *HookContext) {
	scriptname := fmt.Sprintf("./scripts/%s/%s.sh", context.Hook, context.Action)
	log.Println("Exec script: ", scriptname)

	out, err := exec.Command(scriptname, context.args...).Output()
	if err != nil {
		Notify(err.Error(), context)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Notify(fmt.Sprintf("%s", out), context)
	fmt.Fprintf(w, "Action '%s' executed!", context.Action)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	context := new(HookContext)
	context.Hook = params["hookname"]
	context.Action = params["action"]

	log.Println("Hook name: ", context.Hook)

	var record, err = hooks.RecordFactory(context.Hook)
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

	context.args = []string{record.GetURL(), record.GetName()}

	RunScript(w, context)
}

func main() {
	flag.Parse()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{hookname:[a-z]+}/{action:[a-z]+}", Handler).Methods("POST")

	http.Handle("/", rtr)

	log.Println("webhookd server listening...")
	log.Fatal(http.ListenAndServe(*laddr, nil))
}
