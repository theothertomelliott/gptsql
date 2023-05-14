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
	sf "github.com/snowflakedb/gosnowflake"
	"github.com/theothertomelliott/gptsql/conversation"
	"github.com/theothertomelliott/gptsql/schema"
)

func main() {
	var dsn, dbType string
	var err error

	if os.Getenv("POSTGRES_CONN_STRING") != "" {
		dsn = os.Getenv("POSTGRES_CONN_STRING")
		dbType = "postgres"
	}

	if os.Getenv("SNOWFLAKE_ACCOUNT") != "" {
		dsn, _, err = getSnowflakeDSN()
		if err != nil {
			log.Fatal(err)
		}
		dbType = "snowflake"
	}

	if dsn == "" {
		log.Fatal("No database connection config was provided.")
	}

	db, err := sql.Open(dbType, dsn)
	if err != nil {
		log.Fatal(err)
	}

	schema, err := schema.Load(dbType, db)
	if err != nil {
		log.Fatal(err)
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_TOKEN"))
	conv := conversation.New(client, db, dbType, schema)

	samples, err := conv.SampleQuestions()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Welcome to GPTSQL! Your schema has been read and you may ask questions like the below:")
	fmt.Println()
	for _, sample := range samples {
		fmt.Println(sample)
	}
	fmt.Println()

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
				log.Println(err)
				continue
			}

			fmt.Println(res.Query)
			fmt.Println()
			fmt.Println("Data sample:")
			print5Lines(res.DataCsv)
			fmt.Println()
		} else {
			break
		}

	}
}

// print5Lines prints the first 5 lines of the given string
func print5Lines(input string) {
	lines := strings.Split(input, "\n")
	if len(lines) <= 5 {
		fmt.Println(input)
		return
	}

	fmt.Println(strings.Join(lines[0:5], "\n"))
}

// getSnowflakeDSN constructs a DSN based on the test connection parameters
func getSnowflakeDSN() (string, *sf.Config, error) {
	cfg := &sf.Config{
		Authenticator: sf.AuthTypeExternalBrowser,
		Account:       os.Getenv("SNOWFLAKE_ACCOUNT"),
		User:          os.Getenv("SNOWFLAKE_USER"),
	}

	if os.Getenv("SNOWFLAKE_DATABASE") != "" {
		cfg.Database = os.Getenv("SNOWFLAKE_DATABASE")
	}
	if os.Getenv("SNOWFLAKE_SCHEMA") != "" {
		cfg.Schema = os.Getenv("SNOWFLAKE_SCHEMA")
	}

	dsn, err := sf.DSN(cfg)
	return dsn, cfg, err
}
