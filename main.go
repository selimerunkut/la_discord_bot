package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"la_discord_bot/internal/config"
	"la_discord_bot/internal/webserver"
	"log"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Start Server")
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	c := config.Config{}
	exists := false

	if c.PathToStorage, exists = os.LookupEnv("PATH_TO_STORAGE"); !exists {
		log.Fatalf("PATH_TO_STORAGE does not exits in .env")
	}

	f, err := os.OpenFile(c.PathToStorage+"logs/log.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)

	if c.PathToWWW, exists = os.LookupEnv("PATH_TO_WWW"); !exists {
		log.Fatalf("PATH_TO_WWW does not exits in .env")
	}
	if c.GinMode, exists = os.LookupEnv("GIN_MODE"); !exists {
		log.Fatalf("GIN_MODE does not exits in .env")
	}
	if c.ServerPort, exists = os.LookupEnv("PORT"); !exists {
		log.Fatalf("SERVER_PORT does not exits in .env")
	}
	if c.Login, exists = os.LookupEnv("LOGIN"); !exists {
		log.Fatalf("LOGIN does not exits in .env")
	}
	if c.Password, exists = os.LookupEnv("PASSWORD"); !exists {
		log.Fatalf("PASSWORD does not exits in .env")
	}
	if c.AutoStartTasks, exists = os.LookupEnv("AUTO_START_TASKS_AFTER_APP_RELOAD"); !exists {
		log.Fatalf("AUTO_START_TASKS_AFTER_APP_RELOAD does not exits in .env")
	}

	delay, exists := os.LookupEnv("DELAY_MIN")
	if !exists {
		log.Fatalf("DELAY_MIN does not exits in .env")
	}
	if c.DelayMin, err = strconv.ParseFloat(delay, 32); err != nil {
		log.Fatalf("DELAY_MIN not float in .env")
	}

	delay, exists = os.LookupEnv("DELAY_MAX")
	if !exists {
		log.Fatalf("DELAY_MAX does not exits in .env")
	}
	if c.DelayMax, err = strconv.ParseFloat(delay, 32); err != nil {
		log.Fatalf("DELAY_MAX not float in .env")
	}

	//c.DelayMin = f
	//t.SendDelayMax = f

	err = webserver.Init(c)

	if err != nil {
		log.Fatalf("Webserver Init Error: ", err)
	}

	return

}
