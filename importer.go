// LICENCE: No licence is provided for this project

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/net/publicsuffix"
)

var (
	db *sqlx.DB
)

func main() {
	var err error
	db, err = sqlx.Open("postgres", os.Getenv("SCANNER_DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := doImport(); err != nil {
		log.Fatalf("Failed to import: %v", err)
	}

}

func doImport() error {
	names := map[string]int64{}
	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}

	sc := bufio.NewReader(f)
	for {
		if line, err := sc.ReadString('\n'); err != nil {
			log.Printf("Error in scanner: %v\n", err)
			break
		} else {
			split := strings.Split(strings.TrimSpace(strings.ToLower(line)), ",")
			if len(split) != 3 {
				log.Fatalf("Line is bad: %v", split)
			}
			start, err := strconv.ParseInt(split[1], 10, 64)
			if err != nil {
				panic(err)
			}
			end, err := strconv.ParseInt(split[2], 10, 64)
			if err != nil {
				panic(err)
			}
			if start < 0 || end < 0 {
				continue
			}
			names[split[0]] = start
		}
	}

	submitNames(names)

	return nil
}

func submitNames(domains map[string]int64) {
	for name, ts := range domains {
		etld, err := publicsuffix.EffectiveTLDPlusOne(name)
		if err != nil {
			log.Printf("Couldn't determine etld for %s: %v", name, err)
		}

		if _, err := db.Exec(`INSERT INTO domains (domain, first_seen, last_seen, etld) `+
			`VALUES ($1, $2, $2, $3) ON CONFLICT (domain) DO UPDATE `+
			`SET last_seen = GREATEST($2,domains.last_seen), first_seen = LEAST(domains.first_seen, $2);`, name, ts, etld); err != nil {
			log.Printf("Failed to insert/update %s: %v", name, err)
		}
	}
}
