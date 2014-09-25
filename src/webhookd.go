package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ncarlier/webhookd/hook"
	"github.com/ncarlier/webhookd/notification"
	"github.com/ncarlier/webhookd/tools"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

var (
	laddr      = flag.String("l", ":8080", "HTTP service address (e.g.address, ':8080')")
	workingdir = os.Getenv("APP_WORKING_DIR")
	scriptsdir = os.Getenv("APP_SCRIPTS_DIR")
)

type HookContext struct {
	Hook   string
	Action string
	args   []string
}

func Notify(subject string, text string, outfilename string) {
	var notifier, err = notification.NotifierFactory()
	if err != nil {
		log.Println(err)
		return
	}
	if notifier == nil {
		log.Println("Notification provider not found.")
		return
	}

	var zipfile string
	if outfilename != "" {
		zipfile, err = tools.CompressFile(outfilename)
		if err != nil {
			log.Println(err)
			zipfile = outfilename
		}
	}

	notifier.Notify(subject, text, zipfile)
}

func RunScript(w http.ResponseWriter, context *HookContext) {
	scriptname := path.Join(scriptsdir, context.Hook, fmt.Sprintf("%s.sh", context.Action))
	log.Println("Exec script: ", scriptname)

	cmd := exec.Command(scriptname, context.args...)
	var ErrorHandler func(err error, out string)
	ErrorHandler = func(err error, out string) {
		subject := fmt.Sprintf("Webhook %s/%s FAILED.", context.Hook, context.Action)
		Notify(subject, err.Error(), out)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// open the out file for writing
	outfilename := path.Join(workingdir, fmt.Sprintf("%s-%s.txt", context.Hook, context.Action))
	outfile, err := os.Create(outfilename)
	if err != nil {
		ErrorHandler(err, "")
		return
	}

	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		ErrorHandler(err, "")
		return
	}

	err = cmd.Wait()
	if err != nil {
		ErrorHandler(err, outfilename)
		return
	}

	subject := fmt.Sprintf("Webhook %s/%s SUCCEEDED.", context.Hook, context.Action)
	Notify(subject, "See attached file for logs.", outfilename)
	fmt.Fprintf(w, subject)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	context := new(HookContext)
	context.Hook = params["hookname"]
	context.Action = params["action"]

	var record, err = hook.RecordFactory(context.Hook)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Println("Using hook: ", context.Hook)

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&record)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	context.args = []string{record.GetURL(), record.GetName()}

	RunScript(w, context)
}

func main() {
	if workingdir == "" {
		workingdir = os.TempDir()
	}
	if scriptsdir == "" {
		scriptsdir = "scripts"
	}

	flag.Parse()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{hookname:[a-z]+}/{action:[a-z]+}", Handler).Methods("POST")

	http.Handle("/", rtr)

	log.Println("webhookd server listening...")
	log.Fatal(http.ListenAndServe(*laddr, nil))
}
