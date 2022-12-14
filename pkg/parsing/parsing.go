package parsing

import (
	"log"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
)

// The configuration struct of config.toml
type Config struct {
	Core     Core     `toml:"core"`
	LimitReq LimitReq `toml:"limitReq"`
	Proxy    Proxy    `toml:"proxy"`
}

type Core struct {
	Listen            int  `toml:"listen"`
	ApiPort           int  `toml:"apiport"`
	LimitMaxConn      int  `toml:"limitMaxConn"`
	ReadTimeout       int  `toml:"readTimeout"`
	WriteTimeout      int  `toml:"writeTimeout"`
	IdleTimeout       int  `toml:"idleTimeout"`
	EnableServerStats bool `toml:"enableServerStats"`
}

type LimitReq struct {
	Enable    bool `toml:"enable"`
	Interval  int  `toml:"interval"`
	Frequency int  `toml:"frequency"`
	Mode      int  `toml:"mode"`
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
	Url     string `toml:"url"`
	VarName string `toml:"varName"`
}

type File struct {
	FilePath string
}

var (
	config     Config
	once       sync.Once
	ConfigFile *File
)

func init() {
	ConfigFile = new(File)
}

const fileName = "config.toml"

func GetConfig() *Config {
	once.Do(func() {
		// curDir, _ := os.Getwd()
		// filePath := filepath.Join(curDir, "bin", "config", "config.toml")

		_, err := toml.DecodeFile(ConfigFile.FilePath, &config)
		if err != nil {
			log.Fatal(err)
		}
		spew.Dump(config, "\n")
	})
	return &config
}
