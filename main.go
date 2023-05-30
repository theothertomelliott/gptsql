package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	_ "github.com/lib/pq"
	"github.com/sashabaranov/go-openai"
	sf "github.com/snowflakedb/gosnowflake"
	"github.com/theothertomelliott/gptsql/conversation/server"
	"github.com/theothertomelliott/gptsql/schema"
)

func main() {
	var dsn, dbType string
	var useDevFrontEnd bool
	var err error

	if os.Getenv("USE_DEV_FRONTEND") != "" {
		useDevFrontEnd = true
	}

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

	svr := server.New(client, db, dbType, schema)

	mux := http.NewServeMux()

	newConversationHandler := server.GetNewConversationHandler(svr)
	mux.Handle("/new", newConversationHandler)

	askHandler := server.GetAskHandler(svr)
	mux.Handle("/ask", askHandler)

	sampleQuestionsHandler := server.GetSampleQuestionsHandler(svr)
	mux.Handle("/sample-questions", sampleQuestionsHandler)

	if useDevFrontEnd {
		remote, err := url.Parse("http://localhost:3000")
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		mux.Handle("/", proxy)
	} else {
		mux.HandleFunc("/", handleStatic)
	}

	fmt.Println("Listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Println("server failed:", err)
	}
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
