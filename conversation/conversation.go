package conversation

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/joho/sqltocsv"
	"github.com/sashabaranov/go-openai"
	"github.com/theothertomelliott/gptsql/schema"
)

func New(
	client *openai.Client,
	db *sql.DB,
	schema schema.Schema,
) *Conversation {
	return &Conversation{
		client: client,
		db:     db,
		schema: schema,
	}
}

type Conversation struct {
	client *openai.Client
	schema schema.Schema
	db     *sql.DB

	history []Exchange
}

type Request struct {
	Question string
}

type Response struct {
	Query   string
	DataCsv string
}

func (c *Conversation) SampleQuestions() ([]string, error) {
	prompt := fmt.Sprintf(`
	Given this schema

	%v

	Provide three example questions that may be answered using SQL queries against this database.
	Ensure that these questions could be turned into SQL queries using only the schema provided.
	Lean towards questions that aggregate data rather than expecting the user to specify values.
	Do not provide the SQL queries themselves.
	Output questions one per line.
	`, c.schema)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatCompletion error: %w", err)
	}

	response := resp.Choices[0].Message.Content
	return strings.Split(response, "\n"), nil
}

func (c *Conversation) Ask(req Request) (*Response, error) {
	res := &Response{}

	prompt := fmt.Sprintf(`
	Given this schema
	
	%v
	
	Write a SQL query to answer the following question. Use only the content of the schema provided.
	Avoid queries with placeholders. Only output the query, please do not explain it:
	%v`, c.schema, req.Question)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatCompletion error: %w", err)
	}

	query := resp.Choices[0].Message.Content
	res.Query = query
	res.DataCsv, err = c.execQuery(query)
	if err != nil {
		return nil, err
	}

	c.history = append(c.history, Exchange{
		Request:  &req,
		Response: res,
	})
	return res, nil
}

type Exchange struct {
	Request  *Request
	Response *Response
}

// execQuery runs a db query and prints the results in csv format
func (c *Conversation) execQuery(query string) (string, error) {
	rows, err := c.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("query:\n%v\n%w", query, err)
	}
	defer rows.Close()

	out, err := sqltocsv.WriteString(rows)
	if err != nil {
		return "", fmt.Errorf("Rendering to csv: %w", err)
	}
	return out, nil
}
