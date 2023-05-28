package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/theothertomelliott/gptsql/conversation"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

var _ Server = &client{}

type client struct {
	client *http.Client

	newConversationEndpoint endpoint.Endpoint
	sampleQuestionsEndpoint endpoint.Endpoint
	askEndpoint             endpoint.Endpoint
}

func NewClient(host string) *client {
	c := &client{
		client: http.DefaultClient,
	}

	newConversationURL, err := url.Parse(fmt.Sprintf("%v/new", host))
	if err != nil {
		log.Fatal(err)
	}

	c.newConversationEndpoint = httptransport.NewClient(
		"GET",
		newConversationURL,
		encodeRequest,
		func(_ context.Context, r *http.Response) (interface{}, error) {
			var response NewConversationResponse
			if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
				fmt.Println("Error decoding response body: ", err)
				return nil, err
			}
			return response, nil
		},
	).Endpoint()

	sampleQuestionsURL, err := url.Parse(fmt.Sprintf("%v/sample-questions", host))
	if err != nil {
		log.Fatal(err)
	}

	c.sampleQuestionsEndpoint = httptransport.NewClient(
		"GET",
		sampleQuestionsURL,
		encodeRequest,
		func(_ context.Context, r *http.Response) (interface{}, error) {
			var response SampleQuestionsResponse
			if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
				fmt.Println("Error decoding response body: ", err)
				return nil, err
			}
			return response, nil
		},
	).Endpoint()

	askURL, err := url.Parse(fmt.Sprintf("%v/ask", host))
	if err != nil {
		log.Fatal(err)
	}

	c.askEndpoint = httptransport.NewClient(
		"GET",
		askURL,
		encodeRequest,
		func(_ context.Context, r *http.Response) (interface{}, error) {
			var response AskResponse
			if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
				fmt.Println("Error decoding response body: ", err)
				return nil, err
			}
			return response, nil
		},
	).Endpoint()

	return c
}

func (c *client) NewConversation() (ConversationID, error) {
	response, err := c.newConversationEndpoint(
		context.Background(),
		NewConversationRequest{},
	)
	if err != nil {
		return "", err
	}
	resp := response.(NewConversationResponse)
	if resp.Err != "" {
		return "", fmt.Errorf(resp.Err)
	}
	return ConversationID(resp.ConversationID), nil
}

func (c *client) SampleQuestions(cid ConversationID) ([]string, error) {
	response, err := c.sampleQuestionsEndpoint(
		context.Background(),
		SampleQuestionsRequest{
			ConversationID: string(cid),
		},
	)
	if err != nil {
		return nil, err
	}
	resp := response.(SampleQuestionsResponse)
	if resp.Err != "" {
		return nil, fmt.Errorf(resp.Err)
	}
	return resp.Questions, nil
}

func (c *client) Ask(cid ConversationID, question string) (*conversation.Response, error) {
	response, err := c.askEndpoint(
		context.Background(),
		AskRequest{
			ConversationID: string(cid),
			Question:       question,
		},
	)
	if err != nil {
		return nil, err
	}
	resp := response.(AskResponse)
	if resp.Err != "" {
		return nil, fmt.Errorf(resp.Err)
	}
	out := &conversation.Response{}
	out.Query = resp.Query
	out.DataCsv = resp.DataCsv
	if resp.Err != "" {
		out.Error = fmt.Errorf(resp.Err)
	}

	return out, nil
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}
