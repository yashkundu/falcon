package main

import (
	"encoding/json"
	"fmt"

	"github.com/yashkundu/falcon/pkg/parsing"
)

func main() {
	config := *(parsing.GetConfig())
	js, _ := json.MarshalIndent(config, "", " ")
	fmt.Println(string(js))
}
