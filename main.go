package main

import (
	"fmt"
	"os"

	"github.com/mikkelstb/feedfetcher/config"
)

var sq_config config.RepositoryConfig

func main() {
	fmt.Println("Hello World")

	var err error
	cfg, err := config.Read("/Users/mikkel/feedfetcher/config.json")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error: config file not read, aborting")
		os.Exit(1)
	}

	for r := range cfg.Repositories {
		if cfg.Repositories[r].Type == "sqlite3" {
			sq_config = cfg.Repositories[r]
		}
	}

}
