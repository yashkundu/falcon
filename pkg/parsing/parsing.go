package parsing

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

// The configuration struct of config.toml
type Config struct {
	port   int     // The port on which reverse-proxy will run
	routes []Route // The routes which have to be proxied
}

type Route struct {
	endpoint string    //URL Path of incoming request
	match    string    // How to match the path prefix, regex or exact
	backends []Backend // List of the servers to transport requests
}

type Backend struct {
	url string
}

var (
	config Config
	once   sync.Once
)

const fileName = "config.toml"

func GetConfig() *Config {
	once.Do(func() {
		curDir, err1 := os.Getwd()
		if err1 != nil {
			log.Fatal(err1)
		}
		filePath := filepath.Join(filepath.Dir(filepath.Dir(curDir)), "configs", "configs.toml")
		log.Printf("filePath : %s", filePath)

		_, err := toml.DecodeFile(filePath, &config)
		if err != nil {
			log.Fatal(err)
		}

	})
	return &config
}
