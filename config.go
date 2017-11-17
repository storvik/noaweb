package noaweb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Helper type to decode config
type sections map[string]json.RawMessage

// LoadConfig loads configuration file.
func parseConfig(configFile string, instance *Instance) error {

	// Read file to string, f
	var dat []byte
	var err error
	if strings.HasPrefix(configFile, "/") {
		dat, err = ioutil.ReadFile(configFile)
	} else {
		dat, err = ioutil.ReadFile(instance.AssetsDir + "/" + configFile)
	}
	if err != nil {
		fmt.Println(err)
	}
	f := string(dat)

	var s sections
	decoder := json.NewDecoder(strings.NewReader(f))
	err = decoder.Decode(&s)
	if err != nil {
		return err
	}

	if instance.commandFlags.devMode {
		err = json.Unmarshal(s["dev"], &instance)
	} else {
		err = json.Unmarshal(s["prod"], &instance)
	}
	if err != nil {
		return err
	}

	return nil
}
