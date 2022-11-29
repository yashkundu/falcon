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
	Core     Core     `toml:"core"`
	LimitReq LimitReq `toml:"limitReq"`
	Log      Log      `toml:"log"`
	Proxy    Proxy    `toml:"proxy"`
}

type Core struct {
	Listen       int `toml:"listen"`
	ApiPort      int `toml:"apiport"`
	LimitMaxConn int `toml:"limitMaxConn"`
	ReadTimeout  int `toml:"readTimeout"`
	WriteTimeout int `toml:"writeTimeout"`
	IdleTimeout  int `toml:"idleTimeout"`
}

type LimitReq struct {
	Enable    bool `toml:"enable"`
	Interval  int  `toml:"interval"`
	Frequency int  `toml:"frequency"`
	Mode      int  `toml:"mode"`
}

type Log struct {
	Level string `toml:"level"`
}

type Proxy struct {
	Routes []Route `toml:"routes"`
}

type Route struct {
	Endpoint string    `toml:"endpoint"`
	Match    int       `toml:"match"`
	Balancer int       `toml:"balancer"`
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
		curDir, _ := os.Getwd()
		filePath2 := filepath.Join(filepath.Dir(filepath.Dir(curDir)), "bin", "config", "config.toml")
		filePath1 := "C:\\Users\\Hp\\OneDrive\\Desktop\\falcon\\bin\\config\\config.toml"
		log.Printf("filePath : %s\n", filePath2)

		_, err := toml.DecodeFile(filePath1, &config)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Config -> ")
		log.Println(config)
	})
	return &config
}
