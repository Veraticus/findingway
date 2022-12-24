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
				fmt.Printf("unknown option '%s'", arg)
			}
		}
	}

	token, tokenExists := os.LookupEnv("DISCORD_TOKEN")

	if !tokenExists {
		fmt.Printf("Please provide DISCORD_TOKEN")
		os.Exit(1)
	}

	sleepStr, sleepEnvExists := os.LookupEnv("SLEEP")

	if sleepEnvExists {
		sleep64, err := strconv.ParseInt(sleepStr, 10, 64)

		if err != nil {
			fmt.Printf("bad input for --sleep: %s\n", err)
			os.Exit(1)
		}

		sleep = sleep64
	}

	murult.Logger.Println("starting server...")
	server := murult.NewServer(token)

	if server == nil {
		return
	}

	defer server.CloseServer()
	server.Run(sleep)
}
