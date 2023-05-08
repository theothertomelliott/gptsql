package conversation

import (
	"context"
	"database/sql"
	"fmt"

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

func (c *Conversation) Ask(req Request) (*Response, error) {
	res := &Response{}

	prompt := fmt.Sprintf(`
	Given this schema
	
	%v
	
	Write a SQL query to answer the following question, but only output the query, please do not explain it:
	%v`, c.schema, req.Question)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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
		return "", err
	}
	defer rows.Close()

	out, err := sqltocsv.WriteString(rows)
	if err != nil {
		return "", fmt.Errorf("Rendering to csv: %w", err)
	}
	return out, nil
}
