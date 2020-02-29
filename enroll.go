package main

import (
	"os"
	"fmt"
	"flag"
	"path"
	"bytes"
	"strings"
	"os/exec"
	"net/url"
	"net/http"
	"io/ioutil"
	"github.com/Olling/slog"
	"github.com/Olling/Enroll/config"
)

func getHostname() string{
	bytefqdn, err := exec.Command("/bin/hostname", "--fqdn").Output()

	if err == nil {
		fqdn := strings.TrimSpace(string(bytefqdn))
		return fqdn
	}

	hostname, err := os.Hostname()

	if err != nil {
		slog.PrintError("Failed to get FQDN and hostname")
		os.Exit(1)
	}

	return hostname
}

func GetEnrolldStatus(serverid string) {
	url, err := CreateURL(config.Configuration.URL, "status", serverid)
	if err != nil{
		slog.PrintError("Failed to parse the url", err)
		fmt.Println("Status: Unknown")
		os.Exit(1)
	}

	resp, err := http.Get(url)

	if err != nil{
		slog.PrintError(err)
		fmt.Println("Status: Unknown")
		os.Exit(1)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Status: Enrolled")
		os.Exit(0)
	}

	if resp.StatusCode == 202 {
		fmt.Println("Status: Enrolling")
		os.Exit(0)
	}

	if resp.StatusCode == 404 {
		fmt.Println("Status: Not enrolled")
		os.Exit(0)
	}

	fmt.Println("Status: Unknown")
	os.Exit(0)
}

func CreateURL(baseURL string, function string, serverid string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil{
		return "",err
	}

	u.Path = path.Join(u.Path, function)
	u.Path = path.Join(u.Path, serverid)

	return u.String(), nil
}


func Enroll(p config.Payload) {
	slog.PrintDebug("Preparing Enroll")
	url, err := CreateURL(config.Configuration.URL, "server", p.ServerID)
	if err != nil {
		slog.PrintError("Failed to create URL", err)
		return
	}

	slog.PrintDebug("Post url:", url)

	json, err := config.StructToJson(p)

	if err != nil {
		slog.PrintError("Failed to create json payload", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")

	httpclient := &http.Client{}

	resp, err := httpclient.Do(req)
	if err != nil {
		slog.PrintError("Failed to post data", err)
		return
	}

	defer resp.Body.Close()

	slog.PrintDebug("ReponseStatus:",resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	slog.PrintDebug("Server response:", body)

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		slog.PrintError("Got bad response from Enrolld server")
	}
}

func main() {
	config.Initialize()

	if config.Configuration.URL == "" {
		slog.PrintError("No URL was provided")
		os.Exit(1)
	}

	if config.Configuration.Payload.ServerID == "" {
		config.Configuration.Payload.ServerID = getHostname()

	}

	if config.Status {
		GetEnrolldStatus(config.Configuration.Payload.ServerID)
		os.Exit(0)
	}

	if config.Enroll {
		payload := config.GetPayload()

		Enroll(payload)
		GetEnrolldStatus(config.Configuration.Payload.ServerID)
		os.Exit(0)
	}

	flag.PrintDefaults()
}
