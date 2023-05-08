package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sashabaranov/go-openai"
	"github.com/theothertomelliott/gptsql/conversation"
	"github.com/theothertomelliott/gptsql/schema"
)

func main() {
	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	schema, err := schema.Load(db)
	if err != nil {
		log.Fatal(err)
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_TOKEN"))
	conv := conversation.New(client, db, schema)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Ask a question: ")
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		text := scanner.Text()
		if len(text) != 0 {
			res, err := conv.Ask(conversation.Request{
				Question: text,
			})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(res.Query)
			fmt.Println()
			fmt.Println("Data sample:")
			print5Lines(res.DataCsv)
		} else {
			break
		}

	}
}

// print5Lines prints the first 5 lines of the given string
func print5Lines(input string) {
	inputReader := bufio.NewReader(strings.NewReader(input))
	scanner := bufio.NewScanner(inputReader)
	for i := 0; i < 5; i++ {
		scanner.Scan()
		fmt.Println(scanner.Text())
	}
}
