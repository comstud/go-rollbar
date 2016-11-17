package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/comstud/go-rollbar/rollbar"
)

var commands = map[string]func(*rollbar.Client) int{
	"get_item":             getItem,
	"get_item_by_counter":  getItemByCounter,
	"get_occurrence":       getOccurrence,
	"get_item_occurrences": getItemOccurrences,
	"get_occurrences":      getOccurrences,
}

func main() {
	apiToken := os.Getenv("ROLLBARCLI_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Please set ROLLBARCLI_API_TOKEN environment variable")
	}

	client, err := rollbar.NewClient(apiToken)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [<args>]\n", os.Args[0])
		os.Exit(1)
	}

	cmd := os.Args[1]
	fn, ok := commands[cmd]
	if !ok {
		keys := ""
		for k, _ := range commands {
			if keys == "" {
				keys = k
			} else {
				keys += "," + k
			}
		}

		fmt.Fprintf(os.Stderr, "Unknown command. Valid commands are: %s\n", keys)
		os.Exit(1)
	}

	os.Exit(fn(client))
}

func getItem(client *rollbar.Client) int {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s %s <identifier>\n", os.Args[0], os.Args[1])
		return 1
	}
	identifier, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	response, err := client.GetItem(identifier)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}
	if response.IsError() {
		fmt.Printf("Got error: %s\n", response.Message)
		return 1
	}

	fmt.Printf("Got item: %s\n", response.Item.AsPrettyJSON())

	return 0
}

func getItemByCounter(client *rollbar.Client) int {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s %s <identifier>\n", os.Args[0], os.Args[1])
		return 1
	}
	identifier, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	response, err := client.GetItemByCounter(identifier)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}
	if response.IsError() {
		fmt.Printf("Got error: %s\n", response.Message)
		return 1
	}
	fmt.Printf("Got item: %s\n", response.Item.AsPrettyJSON())
	return 0
}

func getOccurrence(client *rollbar.Client) int {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s %s <identifier>\n", os.Args[0], os.Args[1])
		return 1
	}
	identifier, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	response, err := client.GetOccurrence(identifier)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}
	if response.IsError() {
		fmt.Printf("Got error: %s\n", response.Message)
		return 1
	}
	fmt.Printf("Got occurrence: %s\n", response.Occurrence.AsPrettyJSON())
	return 0
}

func getOccurrences(client *rollbar.Client) int {
	var page uint64 = 1
	var err error

	if len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s %s [<page>]\n", os.Args[0], os.Args[1])
		return 1
	}

	if len(os.Args) == 3 {
		page, err = strconv.ParseUint(os.Args[2], 10, 0)
		if err != nil {
			fmt.Printf("%s\n", err)
			return 1
		}
	}

	response, err := client.GetOccurrencesWithPage(page)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}
	if response.IsError() {
		fmt.Printf("Got error: %s\n", response.Message)
		return 1
	}
	fmt.Printf("Got occurrences: %s\n", response.AsPrettyJSON())
	return 0
}

func getItemOccurrences(client *rollbar.Client) int {
	var page uint64 = 1
	var err error

	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s %s <item_id> [<page>]\n", os.Args[0], os.Args[1])
		return 1
	}

	item_id, err := strconv.ParseUint(os.Args[2], 10, 0)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	if len(os.Args) == 4 {
		page, err = strconv.ParseUint(os.Args[3], 10, 0)
		if err != nil {
			fmt.Printf("%s\n", err)
			return 1
		}
	}

	response, err := client.GetItemOccurrencesWithPage(item_id, page)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}
	if response.IsError() {
		fmt.Printf("Got error: %s\n", response.Message)
		return 1
	}
	fmt.Printf("Got occurrences: %s\n", response.AsPrettyJSON())
	return 0
}
