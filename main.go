// Author:  Jonathan Broche and Tom Steele

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/gojhonny/pwnpaste/haveibeenpwnd"
	"github.com/olekukonko/tablewriter"
)

type result struct {
	Email        string
	PasteAccount hibp.PasteAccount
}

func main() {

	//arguments: default one email, input file [-i], download pastebin data (if 404, offer cacheview) or links [default]

	inputfilename := flag.String("i", "", "file containing new line delimited email addresses")
	ccount := flag.Int("c", 10, "maximum number of concurrent requests")
	//downloaddata := flag.Bool("d", false, "download associated pastes")
	flag.Parse()

	emailAddrs := make(map[string]bool) //[key]value

	if i := len(flag.Args()); i > 0 {
		for j := 0; j < i; j++ {
			emailAddrs[flag.Arg(j)] = true
		}
	}

	if *inputfilename != "" { //file
		f, err := os.Open(*inputfilename)

		if err != nil {
			fmt.Printf("[!] Error: %s", err.Error())
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			emailAddrs[scanner.Text()] = true
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("[!] Error: %s", err.Error())
		}

		f.Close()
	}

	results := []result{}

	// Mutex to protect writing in Goroutines.
	resultsMutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	work := make(chan string, *ccount)

	// Start up *ccount goroutine pool to make requests.
	for i := 0; i < *ccount; i++ {
		go func() {
			for email := range work {
				p, _ := hibp.GetPasteAccount(email)
				if len(p) > 0 {
					resultsMutex.Lock()
					results = append(results, result{
						Email:        email,
						PasteAccount: p,
					})
					resultsMutex.Unlock()
				}
				wg.Done()
			}
		}()
	}

	// Send each email to the goroutine pool
	for email := range emailAddrs {
		wg.Add(1)
		work <- email
	}
	close(work)

	wg.Wait()

	outputTable(results)
}

func outputTable(results []result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Email", "Pastebin URL"})
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	for _, r := range results {
		pastes := ""
		for _, p := range r.PasteAccount {
			pastes += fmt.Sprintf("https://pastebin.com/%s\n", p.ID)
		}
		table.Append([]string{r.Email, pastes})
	}
	table.Render() // Send output
}
