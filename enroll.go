package main

import (
	"os"
	"fmt"
	"flag"
	"os/exec"
	"net/http"
	"github.com/Olling/Enroll/config"
	"github.com/Olling/slog"
)

func getHostname() string{
	fqdn, err := exec.Command("/bin/hostname", "--fqdn").Output()

	if err == nil {
		return string(fqdn)
	}

	hostname, err := os.Hostname()

	if err != nil {
		slog.PrintError("Failed to get FQDN and hostname")
		os.Exit(1)
	}

	return hostname
}

func GetEnrolldStatus(fqdn string) {
	r, err := http.Get(config.Configuration.URL + "/status/" + fqdn)

	if err != nil{
		slog.PrintError("Failed to get status", err)
		fmt.Println("Status: Unknown")
		os.Exit(1)
	}

	if r.StatusCode == 200 {
		fmt.Println("Status: Enrolled")
		os.Exit(0)
	}

	if r.StatusCode == 202 {
		fmt.Println("Status: Enrolling")
		os.Exit(0)
	}

	if r.StatusCode == 404 {
		fmt.Println("Status: Not enrolled")
		os.Exit(0)
	}

	fmt.Println("Status: Unknown")
	os.Exit(0)
}

func Enroll(p config.Payload) {
	fmt.Println("NOT READY")
}

func main() {
	config.Initialize()
	hostname := getHostname()

	if config.Configuration.URL == "" {
		slog.PrintError("No URL was provided")
		os.Exit(1)
	}

	if config.Status {
		GetEnrolldStatus(hostname)
		os.Exit(0)
	}

	if config.Enroll {
		payload := config.GetPayload()
		Enroll(payload)
		GetEnrolldStatus(hostname)
		os.Exit(0)
	}

	flag.PrintDefaults()
}
