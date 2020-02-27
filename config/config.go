package config

import (
	"os"
	"encoding/json"
	l "github.com/Olling/Enrolld/logging"
)

var (
	Configuration configuration
	Status		bool
	ConfigPath	string
)

type fragment struct {
	Inventories       []string
	AnsibleProperties map[string]string
}

type payload struct {
	FQDN              string
	NewServer         string
	Inventories       []string
	AnsibleProperties map[string]string
}

type configuration struct {
	ConfigFragments	string
	URL		string
	Payload		payload
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


func LoadFragments(path string, output *payload) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		l.Errorlog.Println("Failed to load files from", path)
		return
	}

	for _, file := range files {
		filepath := path + "/" + file.Name()

		var f fragment
		FileToStruct(filepath, f)

		output.Inventories = append(output.Inventories, f.Inventories)

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
	Configuration.Status := flag.Bool("status", false, "Get current enrollment status")
	flag.Parse()

	file, _ := os.Open(flagConfigPath)
	err := json.NewDecoder(file).Decode(&Configuration)

	if err != nil {
		fmt.Println("Failed to load config (" + ConfigPath + "):", err
	}
	file.Close()

	if flagConfigFragments == "" && ConfigurationConfigFragments == "" {
		ConfigurationConfigFragments = "/etc/enroll/conf.d"
	} else {
		ConfigurationConfigFragments = flagConfigFragments
	}

	if flagURL != "" {
		ConfigurationURL = flagURL
	}
}
