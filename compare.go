package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

const hashesSaveFile string = "jobwatch.json"

var urlToHashMu sync.RWMutex
var urlToHash = make(map[string][]byte)

func hashSha256(body io.ReadCloser) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, body); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func LoadSavedHashes() {
	_, err := os.Stat(hashesSaveFile)
	if os.IsNotExist(err) {
		file, err := os.Create(hashesSaveFile)
		if err != nil {
			fmt.Printf("Failed to create %s\n", hashesSaveFile)
			panic(err)
		}
		defer file.Close()
		fmt.Printf("Created %s\n", hashesSaveFile)
		return
	}

	data, err := os.ReadFile(hashesSaveFile)
	if err != nil {
		fmt.Printf("Failed to read %s\n", hashesSaveFile)
		panic(err)
	}

	urlToHashMu.Lock()
	err = json.Unmarshal(data, &urlToHash)
	urlToHashMu.Unlock()
	if err != nil {
		fmt.Printf("Failed to json-decode %s\n", hashesSaveFile)
		os.Remove(hashesSaveFile)
		panic(err)
	}

	fmt.Printf("Successfully read %s\n", hashesSaveFile)
}

func SaveHashes() {
	urlToHashMu.RLock()
	json, err := json.Marshal(urlToHash)
	urlToHashMu.RUnlock()
	if err != nil {
		fmt.Printf("Failed to json-encode hashes\n")
		return
	}
	err = os.WriteFile(hashesSaveFile, json, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write to %s\n", hashesSaveFile)
		return
	}
}

func VerifyJobBoard(jobBoard JobBoard, body io.ReadCloser) {
	hash, err := hashSha256(body)
	if err != nil {
		fmt.Printf("Failed to process '%s'\n", jobBoard.Url)
		return
	}

	urlToHashMu.RLock()
	oldHash, keyExists := urlToHash[jobBoard.Url]
	urlToHashMu.RUnlock()
	if !keyExists {
		fmt.Printf("First time load for '%s'\n", jobBoard.Url)
		urlToHashMu.Lock()
		urlToHash[jobBoard.Url] = hash
		urlToHashMu.Unlock()
		return
	}

	if !bytes.Equal(hash, oldHash) {
		fmt.Printf("Detected change in '%s'\n", jobBoard.Url)
		go IssueDiscordNotification(jobBoard)
		urlToHashMu.Lock()
		urlToHash[jobBoard.Url] = hash
		urlToHashMu.Unlock()
		return
	}

	fmt.Printf("No change in '%s'\n", jobBoard.Url)
}
