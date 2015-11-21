/* pwnpaste - An OSINT script which queries haveibeenpwned.com for compromised emails.
Pastes respective to the compromised emails are outputted to stdout or exported
to stdout in HTML format. Pwnpaste is meant to act as a quick way to identify
compromises and gather information useful to assessments.

Authors: Jonathan Broche and Tom Steele
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"os"
	"sync"

	"github.com/gojhonny/pwnpaste/haveibeenpwnd"
	"github.com/olekukonko/tablewriter"
	"github.com/toashd/gopher"
)

const version = "1.0.0"

type pasteData struct {
	ID             string
	Data           string
	PasteBinLink   string
	CachedViewLink string
}

type result struct {
	Email  string
	Pastes []pasteData
}

func main() {
	inputfilename := flag.String("i", "", "File containing new line delimited email addresses")
	ccount := flag.Int("c", 10, "Maximum number of concurrent requests")
	htmlOut := flag.Bool("html", false, "Output HTML (has webcache links)")
	showVersion := flag.Bool("v", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	emailAddrs := make(map[string]bool)

	if i := len(flag.Args()); i > 0 {
		for j := 0; j < i; j++ {
			emailAddrs[flag.Arg(j)] = true
		}
	}

	if *inputfilename != "" {
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

	g := gopher.New()

	if !*htmlOut {
		g.Start()
		g.SetSuffix(fmt.Sprintf("pwnpaste %v", version))
		g.SetColor(gopher.Green)
	}

	// Start up *ccount goroutine pool to make requests.
	for i := 0; i < *ccount; i++ {
		go func(j int) {
			for email := range work {
				pasteAccounts, _ := hibp.GetPasteAccount(email)
				if len(pasteAccounts) > 0 {
					pasteDatas := []pasteData{}

					for _, p := range pasteAccounts {
						pasteDatas = append(pasteDatas, pasteData{
							ID:             p.ID,
							PasteBinLink:   fmt.Sprintf("https://pastebin.com/raw.php?i=%s", p.ID),
							CachedViewLink: fmt.Sprintf("https://webcache.googleusercontent.com/search?q=cache:https://pastebin.com/%s", p.ID),
						})
					}
					resultsMutex.Lock()
					results = append(results, result{
						Email:  email,
						Pastes: pasteDatas,
					})
					if !*htmlOut {
						if j%2 == 0 {
							g.SetColor(gopher.Magenta)
							g.SetActivity(gopher.Loving)
						} else {
							g.SetColor(gopher.Green)
							g.SetActivity(gopher.Boring)
						}
					}
					resultsMutex.Unlock()
				}
				wg.Done()
			}
		}(i)
	}

	// Send each email to the goroutine pool
	for email := range emailAddrs {
		wg.Add(1)
		work <- email
	}
	close(work)

	wg.Wait()
	g.Stop()

	if *htmlOut {
		outputHTML(results)
	} else {
		outputTable(results)
	}

}

func outputTable(results []result) {
	fmt.Printf("\npwnpaste %s\n\n", version)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Email", "Pastebin URL", "CachedView URL"})
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	for _, r := range results {
		pastes := ""
		cviews := ""
		for _, p := range r.Pastes {
			pastes += p.PasteBinLink + "\n"
			cviews += p.CachedViewLink + "\n"
		}
		table.Append([]string{r.Email, pastes, cviews})

	}
	table.Render()
}

func outputHTML(results []result) {
	tmpl := `
<html>
<head>
    <style>
        body {
            font-family: arial, sans-serif;
            font-size: 14px;
            margin: 10px 0 0 20px;
        }
        table {
            border-collapse: collapse;
            margin-top: 25px;
        }
        th {
            background-color: #0057b8;
            color: #fff;
        }
        table th,
        td {
            border: 1px solid #000;
            padding: 10px;
        }
        tr:nth-child(odd) {
            background: #e1e1e1;
        }
        ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }
        li:not(:last-child) {
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <h2>pwnpaste</h2>

    <table>
        <tr>
            <th>Email</th>
            <th>PasteBin URL</th>
            <th>CachedView URL</th>
        </tr>
        {{range $result := .}}
        <tr>
            <td style="vertical-align: top;">
                {{$result.Email}}
            </td>
            <td>
                <ul>
                    {{range $paste := $result.Pastes}}
                    <li><a href="{{$paste.PasteBinLink}}" target="_blank">{{$paste.PasteBinLink}}</a>
                    </li>
                    {{end}}
            </td>
            <td>
                <ul>
                    {{range $paste := $result.Pastes}}
                    <li><a href="{{$paste.CachedViewLink}}" target="_blank">{{$paste.CachedViewLink}}</a>
                    </li>
                    {{end}}
            </td>
        </tr>
        {{end}}
    </table>
</body>
</html>
`
	t, _ := template.New("PWNPASTE").Parse(tmpl)
	t.Execute(os.Stdout, results)
}
