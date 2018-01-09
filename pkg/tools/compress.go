package tools

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"

	"github.com/ncarlier/webhookd/pkg/logger"
)

// CompressFile is a simple file gzipper.
func CompressFile(filename string) (zipfile string, err error) {
	zipfile = fmt.Sprintf("%s.gz", filename)
	in, err := os.Open(filename)
	if err != nil {
		return
	}
	out, err := os.Create(zipfile)
	if err != nil {
		logger.Error.Println("Unable to create gzip file", err)
		return
	}

	// buffer readers from file, writes to pipe
	bufin := bufio.NewReader(in)

	// gzip wraps buffer writer and wr
	gw := gzip.NewWriter(out)
	defer gw.Close()

	_, err = bufin.WriteTo(gw)
	if err != nil {
		logger.Error.Println("Unable to write into the gzip file", err)
		return
	}
	logger.Debug.Println("Gzip file created: ", zipfile)
	return
}
