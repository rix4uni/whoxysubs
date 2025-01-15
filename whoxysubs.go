package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
	"github.com/spf13/pflag"
)

// prints the version message
const version = "0.0.1"

func printVersion() {
	fmt.Printf("Current whoxysubs version %s\n", version)
}

// Prints the Colorful banner
func printBanner() {
	banner := `
            __                                      __         
 _      __ / /_   ____   _  __ __  __ _____ __  __ / /_   _____
| | /| / // __ \ / __ \ | |/_// / / // ___// / / // __ \ / ___/
| |/ |/ // / / // /_/ /_>  < / /_/ /(__  )/ /_/ // /_/ /(__  ) 
|__/|__//_/ /_/ \____//_/|_| \__, //____/ \__,_//_.___//____/  
                            /____/`
fmt.Printf("%s\n%70s\n\n", banner, "Current whoxysubs version "+version)
}

type DomainInfo struct {
	Num         int    `json:"num"`
	DomainName  string `json:"domain_name"`
	Registrar   string `json:"registrar"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	Expiry      string `json:"expiry"`
}

func main() {
	// Define the search type flag
	searchType := pflag.StringP("search", "s", "", "Search type: company, email, keyword, or name")
	silent := pflag.Bool("silent", false, "silent mode.")
	version := pflag.Bool("version", false, "Print the version of the tool and exit.")
	pflag.Parse()

	// Print version and exit if -version flag is provided
	if *version {
		printBanner()
		printVersion()
		return
	}

	// Don't Print banner if -silnet flag is provided
	if !*silent {
		printBanner()
	}

	// Validate the search type
	validSearchTypes := map[string]bool{
		"company": true,
		"email":   true,
		"keyword": true,
		"name":    true,
	}
	if !validSearchTypes[*searchType] {
		log.Fatalf("Invalid search type: %s. Valid options are: company, email, keyword, or name.", *searchType)
	}

	// Read input from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Get the query from stdin
		query := strings.TrimSpace(scanner.Text())
		query = url.QueryEscape(query)

		// Construct the appropriate search URL based on the search type
		var searchURL string
		switch *searchType {
		case "company":
			searchURL = fmt.Sprintf("https://www.whoxy.com/search.php?company=%s", query)
		case "email":
			searchURL = fmt.Sprintf("https://www.whoxy.com/search.php?email=%s", query)
		case "keyword":
			searchURL = fmt.Sprintf("https://www.whoxy.com/search.php?keyword=%s", query)
		case "name":
			searchURL = fmt.Sprintf("https://www.whoxy.com/search.php?name=%s", query)
		}

		// Fetch the redirected URL
		resp, err := http.Get(searchURL)
		if err != nil {
			log.Fatalf("Failed to fetch data: %v", err)
		}
		defer resp.Body.Close()

		// Parse the HTML response
		doc, err := html.Parse(resp.Body)
		if err != nil {
			log.Fatalf("Failed to parse HTML: %v", err)
		}

		var domains []DomainInfo
		var parseTable func(*html.Node)
		parseTable = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "table" {
				for _, attr := range n.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "grid first_col_center") {
						domains = extractTableData(n)
						return
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				parseTable(c)
			}
		}

		parseTable(doc)

		// Convert the data to JSON
		output, err := json.MarshalIndent(domains, "", "  ")
		if err != nil {
			log.Fatalf("Failed to convert data to JSON: %v", err)
		}

		// Print the JSON output
		fmt.Println(string(output))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

func extractTableData(table *html.Node) []DomainInfo {
	var rows []DomainInfo
	var parseRow func(*html.Node)
	parseRow = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var row DomainInfo
			var cellIdx int
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					cellIdx++
					text := extractText(c)
					switch cellIdx {
					case 1:
						fmt.Sscanf(text, "%d", &row.Num)
					case 2:
						row.DomainName = text
					case 3:
						row.Registrar = text
					case 4:
						row.Created = text
					case 5:
						row.Updated = text
					case 6:
						row.Expiry = text
					}
				}
			}
			if row.DomainName != "" {
				rows = append(rows, row)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseRow(c)
		}
	}

	for c := table.FirstChild; c != nil; c = c.NextSibling {
		parseRow(c)
	}

	return rows
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractText(c)
	}
	return strings.TrimSpace(result)
}
