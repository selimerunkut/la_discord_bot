package main

import (
	"fmt"
	"la_discord_bot/internal/webserver"
)

func main() {

	err := webserver.Init()
	if err != nil {
		fmt.Println("Webserver Init Error: ", err)
	}

	return

}
