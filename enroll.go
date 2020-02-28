package main

import (
	"os"
	"fmt"
	"flag"
	"path"
	"strings"
	"os/exec"
	"net/url"
	"net/http"
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
	u, err := url.Parse(config.Configuration.URL)
	if err != nil{
		slog.PrintError("Failed to parse the url", err)
		fmt.Println("Status: Unknown")
		os.Exit(1)
	}
	u.Path = path.Join(u.Path, "status")
	u.Path = path.Join(u.Path, serverid)

	r, err := http.Get(u.String())

	if err != nil{
		slog.PrintError(err)
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
	slog.PrintDebug("NOT READY")
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
