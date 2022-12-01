package flag

import (
	"flag"
	"log"
	"os"

	"github.com/yashkundu/falcon/pkg/parsing"
)

var (
	Rand string
)

func run() int {
	log.Println("No config file path provided")
	return 1
}

func init() {
	Rand = ""
	path := flag.String("config", "", "The absolute path to the config file")
	flag.Parse()
	parsing.ConfigFile.FilePath = *path

	if *path == "" {
		os.Exit(run())
	}

}
