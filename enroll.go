package main

import (
	"os"
	"log"
	"time"
	"github.com/Olling/Enroll/config"
	l "github.com/Olling/Enroll/logging"
)

func getHostname {
	fqdn, err := exec.Command("/bin/hostname", "--fqdn").Output()

	if err == nil {
		return fqdn
	}

	hostname, err := os.Hostname()

	if err != nil {
		l.ErrorLog.Println("Failed to get FQDN and hostname")
		os.Exit(1)
	}

	return hostname
}

func GetEnrolldStatus(fqdn string) {
	r, err := http.Get(config.Configuration.URL + "/server/" + fqdn)

	if err != nil{
		l.ErrorLog.Println("Failed to get status", err)
		fmt.Println("Status: Unknown")
		os.Exit(1)
	}

	if r.StatusCode == 200 {
		fmt.Println("Status: Enrolled")
		os.Exit(0)
	}

	if r.StatusCode == 404 {
		fmt.Println("Status: Not enrolled")
		os.Exit(0)
	}

	fmt.Println("Status: Unknown")
	os.Exit(1)
}

func Post() {
	fmt.Println("NOT READY")
}

func main() {
	l.InitializeLogging(os.Stdout, os.Stderr)
	config.Initialize()
	hostname := getHostname()

	if Configuration.URL == "" {
		l.ErrorLog.Println("No URL was provided")
		exit(1)
	}

	if Configuration.Status {
		GetEnrolldStatus(hostname)
	}

	var p Payload
	StructFromFile(config.Configuration.ConfigPath


	//	config.GetMainConfiguration("/etc/enroll/enroll.conf")
	//config.GetAdditionalConfiguration("/etc/enroll/enroll.d")

}
