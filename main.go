package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type JobBoard struct {
	Url string `json:"url"`
}

const jobBoardsFile = "jobboards.json"

var jobBoards []JobBoard

func loadJobBoardsFromDisk() {
	_, err := os.Stat(jobBoardsFile)
	if os.IsNotExist(err) {
		fmt.Printf("You must specify a list of jobboards in %s\n", jobBoardsFile)
		panic(err)
	}

	data, err := os.ReadFile(jobBoardsFile)
	if err != nil {
		fmt.Printf("Failed to read %s\n", jobBoardsFile)
		panic(err)
	}

	err = json.Unmarshal(data, &jobBoards)
	if err != nil {
		fmt.Printf("Failed to json-decode %s\n", jobBoardsFile)
		panic(err)
	}

	fmt.Printf("Loaded %s\n", jobBoardsFile)
}

func main() {
	fmt.Println("Jobmaxxing...")

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file")
		panic(err)
	}

	loadJobBoardsFromDisk()
	LoadSavedHashes()
	BeginJobWatch(jobBoards)
}
