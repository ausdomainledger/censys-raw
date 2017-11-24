package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	workers = 16
)

var (
	workCh chan []byte
	wg     sync.WaitGroup
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	workCh = make(chan []byte, 1024)

	for i := 0; i < workers; i++ {
		go worker()
		wg.Add(1)
	}

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}
		workCh <- l
	}

	close(workCh)
	wg.Wait()
}

func worker() {
	defer wg.Done()
	var out struct {
		Parsed struct {
			Names    []string `json:"names"`
			Validity struct {
				Start time.Time `json:"start"`
				End   time.Time `json:"end"`
			} `json:"validity"`
		} `json:"parsed"`
	}
	var names []string
	for b := range workCh {

		if err := json.Unmarshal(b, &out); err != nil {
			log.Fatal(err)
		}

		names = nil
		for _, v := range out.Parsed.Names {
			v = strings.ToLower(strings.TrimSpace(v))
			if !strings.HasSuffix(v, ".au") {
				continue
			}
			names = append(names, v)
		}

		for _, n := range names {
			fmt.Println(n + "," + strconv.FormatInt(out.Parsed.Validity.Start.Unix(), 10) + "," + strconv.FormatInt(out.Parsed.Validity.End.Unix(), 10))
		}

	}
}
