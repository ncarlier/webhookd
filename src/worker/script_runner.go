package worker

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

var (
	workingdir = os.Getenv("APP_WORKING_DIR")
	scriptsdir = os.Getenv("APP_SCRIPTS_DIR")
)

func RunScript(work *WorkRequest) (string, error) {
	if workingdir == "" {
		workingdir = os.TempDir()
	}
	if scriptsdir == "" {
		scriptsdir = "scripts"
	}

	scriptname := path.Join(scriptsdir, work.Name, fmt.Sprintf("%s.sh", work.Action))
	fmt.Println("Exec script: ", scriptname, "...")

	// Exec script...
	cmd := exec.Command(scriptname, work.Args...)

	// Open the out file for writing
	outfilename := path.Join(workingdir, fmt.Sprintf("%s-%s.txt", work.Name, work.Action))
	outfile, err := os.Create(outfilename)
	if err != nil {
		return "", err
	}

	defer outfile.Close()
	cmd.Stdout = outfile
	cmd.Stderr = outfile

	err = cmd.Start()
	if err != nil {
		return outfilename, err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Exec script: ", scriptname, "KO!")
		return outfilename, err
	}

	fmt.Println("Exec script: ", scriptname, "OK")
	return outfilename, nil
}
