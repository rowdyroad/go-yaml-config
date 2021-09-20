package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/jinzhu/copier"
	log "github.com/rowdyroad/go-simple-logger"
)

//LoadConfigFromFile loading config from yaml file
func LoadConfigFromFile(config interface{}, configFile string, defaultValue interface{}) string {
	log.Debugf("Reading configuration from '%s'", configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Warn("Configuration not found")
		if defaultValue != nil {
			log.Warn("Default value is defined. Using it.")
			copier.Copy(config, defaultValue)
			return ""
		}
		panic(err)
	}

	if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(data))), config); err != nil {
		log.Warn("Configuration incorrect ")
		if defaultValue != nil {
			log.Warn("Default value is defined. Use it.")
			copier.Copy(config, defaultValue)
			return ""
		}
		panic(err)
	}

	lastConfigFile := configFile
	customDir := filepath.Dir(configFile)
	ext := filepath.Ext(configFile)
	baseNameNoExt := strings.TrimSuffix(filepath.Base(configFile), ext)
	customConfigFile := filepath.Join(
		customDir,
		fmt.Sprintf("%s%s%s", baseNameNoExt, ".custom", ext),
	)
	for i := 1; i < 100; i++ {
		log.Debugf("Try to read custom configuration from '%s'...", customConfigFile)
		data, err = ioutil.ReadFile(customConfigFile)
		if err == nil {
			log.Debugf("Reading custom configuration from '%s'", customConfigFile)
			if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(data))), config); err != nil {
				panic(err)
			}
			log.Debugf("Config loaded successfully with custom config file '%s'", customConfigFile)
			lastConfigFile = customConfigFile
			customConfigFile = filepath.Join(
				customDir,
				fmt.Sprintf("%s%s.%02d%s", baseNameNoExt, ".custom", i, ext),
			)
		} else {
			break
		}
	}

	log.Debug("Config loaded successfully")
	return lastConfigFile
}

//LoadConfig from command line argument
func LoadConfig(config interface{}, defaultFilename string, defaultValue interface{}) string {
	var (
		configFile                      string
		dumpJson, dumpYaml, shouldPanic bool
	)
	flag.StringVar(&configFile, "c", defaultFilename, "Config file")
	flag.StringVar(&configFile, "config", defaultFilename, "Config file")
	flag.BoolVar(&dumpJson, "dump-json", false, "Dump result config at json")
	flag.BoolVar(&dumpYaml, "dump-yaml", false, "Dump result config at yaml")
	flag.BoolVar(&shouldPanic, "panic", false, "panic after config")
	flag.Parse()
	r := LoadConfigFromFile(config, configFile, defaultValue)
	if dumpJson {
		s, e := DumpJson(config)
		fmt.Println(s)
		if e != nil {
			panic(e)
		}
	}
	if dumpYaml {
		s, e := DumpYaml(config)
		fmt.Println(s)
		if e != nil {
			panic(e)
		}
	}

	if shouldPanic {
		panic(config)
	}
	return r
}

func DumpYaml(config interface{}) (string, error) {
	r, e := yaml.Marshal(config)
	return string(r), e
}

func DumpJson(config interface{}) (string, error) {
	r, e := json.Marshal(config)
	return string(r), e
}
