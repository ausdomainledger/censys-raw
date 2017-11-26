package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	out := map[string][]int64{}

	defer f.Close()

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		l = strings.ToLower(strings.TrimSpace(l))

		split := strings.Split(l, ",")
		if len(split) != 3 {
			fmt.Fprintf(os.Stderr, "Invalid line: %v", l)
			continue
		}

		d := split[0]
		start, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		end, err := strconv.ParseInt(split[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		existing, ok := out[d]
		if !ok {
			out[d] = []int64{start, end}
		} else {
			if existing[0] > start {
				existing[0] = start
			}
			if existing[1] < end {
				existing[1] = end
			}
		}
	}

	for d, t := range out {
		fmt.Println(d, t)
	}
}
