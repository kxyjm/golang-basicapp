package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	path := os.Getenv("PATH")
	cwd, _ := os.Getwd()
	if port == "" {
		port = "5000"
	}

	f, _ := os.Create("/var/log/golang-server.log")
	defer f.Close()
	log.SetOutput(f)
	const result = "public/index.html"
	var smokeData map[string]string
	smokeData = make(map[string]string)
	output, _ := exec.Command("sh", "-c", "echo This is version one").Output()
	smokeData["version"] = string(output)
	smokeData["golang_version"] = runtime.Version()
	output, _ = exec.Command("sh", "-c", "yum list installed openssl").Output()
	smokeData["openssl"] = string(output)
	output, _ = exec.Command("sh", "-c", "ls -ald /var/app/current").Output()
	smokeData["deploy_dir_stat"] = string(output)
	output, _ = exec.Command("sh", "-c", "printenv").Output()
	smokeData["env_varibale"] = string(output)
	output, _ = ioutil.ReadFile("/etc/yum/vars/guid")
	smokeData["GUID"] = strings.Trim(string(output), "\n")
	smokeData["PATH"] = path
	smokeData["PORT"] = port
	smokeData["CWD"] = cwd
	smokeData["ENVS"] = ""
	smokeData["banner"] = "Congratulations"
	smokeData["procfile"] = "This is the basic source test with build and procfile."
	smokeData["THE ANSWER in GO"] = "42"

	for _, s := range os.Environ() {
		smokeData["ENVS"] += s + ";"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if buf, err := ioutil.ReadAll(r.Body); err == nil {
				log.Printf("Received message: %s\n", string(buf))
			}
		} else {
			log.Printf("Serving %s to %s...\n", result, r.RemoteAddr)
			jsonString, _ := JSONMarshal(smokeData, true)
			fmt.Fprintf(w, string(jsonString))
			fmt.Printf("----------------\n")
			for _, s := range os.Environ() {
				fmt.Printf("%s\n", s)
			}
			fmt.Printf("----------------\n")
			fmt.Printf("%s", string(jsonString))
		}
	})

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}

func JSONMarshal(v interface{}, unescape bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if unescape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}
