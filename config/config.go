package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

//Configuration struct holds the config and constants file throughout the projects

type Config struct {
	AppName       string `json:"app_name"`
	HTTPAddress   string `json:"http_address"`
	Port          int    `json:"port"`
	ReadTimeout   int    `json:"read_timeout"`
	WriteTimeout  int    `json:"write_timeout"`
	DatabaseName  string `json:"dbName"`
	RetryDBInsert int    `json:"retry_db_insert"`
	DelimsL       string `json:"delims_l"`
	DelimsR       string `json:"delims_r"`
	Endpoint      string `json:"endpoint"`
}

var (
	cfg  Config
	once sync.Once
)

// Parse parses the json configuration file
// And converting it into native type
func Parse(file string) *Config {
	once.Do(func() {
		// Reading the flags
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalln("config: ioutil.ReadFile failed: ", err)
		}

		if err := json.Unmarshal(data, &cfg); err != nil {
			log.Fatalln("config: json.unmarshal failed: ", err)
		}
	})

	return &cfg
}
