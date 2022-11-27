package parsing

import (
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

// The configuration struct of config.toml
type Config struct {
	Core     Core     `toml:"core"`
	LimitReq LimitReq `toml:"limitReq"`
	Log      Log      `toml:"log"`
	Proxy    Proxy    `toml:"proxy"`
}

type Core struct {
	Listen       int `toml:"listen"`
	LimitMaxConn int `toml:"limitMaxConn"`
	ReadTimeout  int `toml:"readTimeout"`
	WriteTimeout int `toml:"writeTimeout"`
	IdleTimeout  int `toml:"idleTimeout"`
}

type LimitReq struct {
	Enable    bool `toml:"enable"`
	Interval  int  `toml:"interval"`
	Frequency int  `toml:"frequency"`
}

type Log struct {
	Level string `toml:"level"`
}

type Proxy struct {
	Routes []Route `toml:"routes"`
}

type Route struct {
	Endpoint string    `toml:"endpoint"`
	Match    string    `toml:"match"`
	Backends []Backend `toml:"backends"`
}

type Backend struct {
	Url string `toml:"url"`
}

var (
	config Config
	once   sync.Once
)

const fileName = "config.toml"

func GetConfig() *Config {
	once.Do(func() {
		filePath := "C:\\Users\\Hp\\OneDrive\\Desktop\\falcon\\configs\\config.toml"
		log.Printf("filePath : %s", filePath)

		_, err := toml.DecodeFile(filePath, &config)
		if err != nil {
			log.Fatal(err)
		}

	})
	return &config
}
