package conversation

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Exchange struct {
	Request  *Request
	Response *Response
}

type Request struct {
	Question string
}

type Response struct {
	Query   string
	DataCsv string
	Error   error
}

func (e *Exchange) toMessages() []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage
	messages = append(messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(
			"Please answer this question in the form of an SQL query, do not explain your response:\n%v",
			e.Request.Question,
		),
	})
	if e.Response != nil {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: e.Response.Query,
		})
		if e.Response.DataCsv != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("Sample data from the above query:\n%v", upTo5Lines(e.Response.DataCsv)),
			})
		}
		if e.Response.Error != nil {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("The above query returned the error: %v", e.Response.Error),
			})
		}
	}
	return messages
}
