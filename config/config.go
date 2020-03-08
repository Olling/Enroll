package config

import (
	"os"
	"flag"
	"io/ioutil"
	"encoding/json"
	"github.com/Olling/slog"
)

var (
	Configuration	configuration
	Status		*bool
	Enroll		*bool
	ConfigPath	*string
)

type fragment struct {
	Inventories       []string
	AnsibleProperties map[string]string
}

type Payload struct {
	ServerID          string
	NewServer         bool
	Inventories       []string
	AnsibleProperties map[string]string
}

type configuration struct {
	ConfigFragments	string
	URL		string
	Payload		Payload
}


func StructToJson(s interface{}) (string, error) {
	bytes, marshalErr := json.MarshalIndent(s, "", "\t")
	return string(bytes), marshalErr
}


func StructFromJson(input string, output interface{}) error {
	return json.Unmarshal([]byte(input), &output)
}


func FileToStruct(path string, s interface{}) error {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	file.Close()

	return nil
}

func GetPayload() (p Payload) {
	p = Configuration.Payload
	LoadFragments(Configuration.ConfigFragments, &p)
	return p
}

func LoadFragments(path string, output *Payload) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		slog.PrintDebug("Failed to load files from", path)
		return
	}

	for _, file := range files {
		filepath := path + "/" + file.Name()

		var f fragment
		FileToStruct(filepath, f)

		output.Inventories = append(output.Inventories, f.Inventories...)

		for key, value := range output.AnsibleProperties {
			f.AnsibleProperties[key] = value
		}
		output.AnsibleProperties = f.AnsibleProperties
	}
}

func Initialize() {
	flagURL := flag.String("url","","Enrolld host url")
	ConfigPath := flag.String("config","/etc/enroll/enroll.conf","Main configuration file (enroll.conf)")
	flagConfigFragments := flag.String("fragments","","Path to additional configuration files (enroll.d)")
	Status = flag.Bool("status", false, "Get current enrollment status")
	Enroll = flag.Bool("enroll", true, "Call the Enrolld server")

	flag.Parse()

	file, err := os.Open(*ConfigPath)
	if err != nil {
		slog.PrintError("Failed to load config (" + *ConfigPath + "):", err)
		os.Exit(1)
	}

	err = json.NewDecoder(file).Decode(&Configuration)
	if err != nil {
		slog.PrintError("Failed to decode config (" + *ConfigPath + "):", err)
	}

	file.Close()

	if *flagConfigFragments != "" && Configuration.ConfigFragments != "" {
		Configuration.ConfigFragments = *flagConfigFragments
	}

	if *flagURL != "" {
		Configuration.URL = *flagURL
	}
}
