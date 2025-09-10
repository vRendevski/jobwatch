package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func watchUrl(jobBoard JobBoard) {
	resp, err := http.Get(jobBoard.Url)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer resp.Body.Close()
	VerifyJobBoard(jobBoard, resp.Body)
}

func BeginJobWatch(jobBoards []JobBoard) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		var wg sync.WaitGroup
		for _, jobBoard := range jobBoards {
			wg.Go(func() {
				watchUrl(jobBoard)
			})
		}

		wg.Wait()
		SaveHashes()
	}
}
