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

const model = openai.GPT3Dot5Turbo

func New(
	client *openai.Client,
	db *sql.DB,
	dbType string,
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
	dbType string

	history []Exchange
}

func (c *Conversation) schemaPromptMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: fmt.Sprintf("Use the following schema to answer questions\nThe database type is %v\n\n%v\n\n", c.dbType, c.schema),
	}
}

func (c *Conversation) SampleQuestions() ([]string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				c.schemaPromptMessage(),
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `
					Provide three example questions that may be answered using SQL queries against this database.
					Ensure that these questions could be turned into SQL queries using only the schema provided.
					Lean towards questions that aggregate data rather than expecting the user to specify values.
					Do not provide the SQL queries themselves.
					Output questions one per line.
					`,
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

	var messages []openai.ChatCompletionMessage
	messages = append(messages, c.schemaPromptMessage())
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: `You are a chatbot that answers questions about a database in the form of SQL queries.
		You will only use the content from the schema provided to answer questions.
		Avoid queries with placeholders.`,
	})

	for _, exchange := range c.history {
		messages = append(messages, exchange.toMessages()...)
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(
			"Please answer this question in the form of an SQL query, do not explain your response:\n%v",
			req.Question,
		),
	})

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
			N:        3,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatCompletion error: %w", err)
	}

	c.history = append(c.history, Exchange{
		Request:  &req,
		Response: res,
	})

	for _, choice := range resp.Choices {
		res.Query = choice.Message.Content
		res.DataCsv, res.Error = c.execQuery(choice.Message.Content)
		if res.Error == nil {
			break
		}
	}

	if res.Error != nil {
		return nil, res.Error
	}
	return res, nil
}

// upTo5Lines returns up to the first 5 lines of the given string
func upTo5Lines(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) <= 5 {
		return input
	}

	return strings.Join(lines[0:5], "\n")
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
		return "", fmt.Errorf("rendering query: %w", err)
	}
	return out, nil
}
