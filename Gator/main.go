package main

import (
	"fmt"

	config "github.com/Maxy-Maxwell/go-tutorials/Gator/internal"
)

func main() {
	conf := config.Read()
	conf.SetUser("Maxy")
	conf = config.Read()
	prettyPrintConfig(conf)
}

func prettyPrintConfig(c config.Config) {
	fmt.Println("Current_user_name:", c.Current_user_name)
	fmt.Println("Db_url:", c.Db_url)
}
