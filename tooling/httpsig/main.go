package main

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-fed/httpsig"
	configflag "github.com/ncarlier/webhookd/pkg/config/flag"
)

type config struct {
	KeyID   string `flag:"key-id" desc:"Signature key ID"`
	KeyFile string `flag:"key-file" desc:"Private key file (PEM format)" default:"./key.pem"`
	JSON    string `flag:"json" desc:"JSON payload"`
}

func main() {
	conf := &config{}
	configflag.Bind(conf, "HTTP_SIG")

	flag.Parse()

	if conf.KeyID == "" {
		log.Fatal("missing key ID")
	}

	args := flag.Args()
	if len(args) <= 0 {
		log.Fatal("missing target URL")
	}
	targetURL := args[0]
	if _, err := url.Parse(targetURL); err != nil {
		log.Fatal("invalid target URL")
	}

	keyBytes, err := os.ReadFile(conf.KeyFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	pemBlock, _ := pem.Decode(keyBytes)
	if pemBlock == nil {
		log.Fatal("invalid PEM format")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	var payload io.Reader
	var jsonBytes []byte
	if conf.JSON != "" {
		var err error
		jsonBytes, err = os.ReadFile(conf.JSON)
		if err != nil {
			log.Fatal(err.Error())
		}
		payload = bytes.NewReader(jsonBytes)
	}

	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256
	headers := []string{httpsig.RequestTarget, "date"}
	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headers, httpsig.Signature, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	req, err := http.NewRequest("POST", targetURL, payload)
	if err != nil {
		log.Fatal(err.Error())
	}
	if payload != nil {
		req.Header.Add("content-type", "application/json")
	}
	req.Header.Add("date", time.Now().UTC().Format(http.TimeFormat))

	if err = signer.SignRequest(privateKey, conf.KeyID, req, jsonBytes); err != nil {
		log.Fatal(err.Error())
	}

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	scanner := bufio.NewScanner(strings.NewReader(string(dump)))
	for scanner.Scan() {
		fmt.Println(">", scanner.Text())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	dump, err = httputil.DumpResponse(res, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	scanner = bufio.NewScanner(strings.NewReader(string(dump)))
	for scanner.Scan() {
		fmt.Println("<", scanner.Text())
	}
}
