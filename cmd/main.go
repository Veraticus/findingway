package main

import (
	"fmt"
	"os"
	"strconv"

	murult "github.com/yuyuriyuri/trappingway-go/internal"
)

func main() {
	var sleep int64 = 3

	for i := 1; i < len(os.Args); i++ {
		arg := string(os.Args[i])
		switch arg {
		default:
			{
				fmt.Printf("Unknown option '%s'\n", arg)
			}
		}
	}

	token, tokenExists := os.LookupEnv("DISCORD_TOKEN")

	if !tokenExists {
		fmt.Println("Please provide DISCORD_TOKEN")
		os.Exit(1)
	}

	path, pathExists := os.LookupEnv("SQLITE_DB")

	if !pathExists {
		fmt.Println("Please provide SQLITE_DB")
		os.Exit(1)
	}

	sleepStr, sleepEnvExists := os.LookupEnv("SLEEP")

	if sleepEnvExists {
		sleep64, err := strconv.ParseInt(sleepStr, 10, 64)

		if err != nil {
			fmt.Printf("Bad input for SLEEP: %s\n", err)
			os.Exit(1)
		}

		sleep = sleep64
	}

	murult.Logger.Println("Starting server...")
	server := murult.NewServer(token, path)

	if server == nil {
		return
	}

	defer server.CloseServer()
	server.Run(sleep)
}
