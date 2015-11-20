// Author:  Jonathan Broche and Tom Steele

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gojhonny/pwnpaste/haveibeenpwnd"
	"github.com/olekukonko/tablewriter"
	"os"
)

type result struct {
	Email        string
	PasteAccount hibp.PasteAccount
}

func main() {

	//arguments: default one email, input file [-i], download pastebin data (if 404, offer cacheview) or links [default]

	inputfilename := flag.String("i", "", "file containing new line delimited email addresses")
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

	for email := range emailAddrs {
		p, _ := hibp.GetPasteAccount(email)
		if len(p) > 0 {
			results = append(results, result{
				Email:        email,
				PasteAccount: p,
			})
		}
	}

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
