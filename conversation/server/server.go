package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/theothertomelliott/gptsql/conversation"
)

type Server interface {
	SampleQuestions() ([]string, error)
	Ask(question string) (*conversation.Response, error)

	// TODO: Allow editing the SQL for a given question
}

type conversationServer struct {
	conv *conversation.Conversation
}

func New(conv *conversation.Conversation) Server {
	return &conversationServer{conv: conv}
}

func (s *conversationServer) SampleQuestions() ([]string, error) {
	return s.conv.SampleQuestions()
}

func (s *conversationServer) Ask(question string) (*conversation.Response, error) {
	return s.conv.Ask(
		conversation.Request{
			Question: question,
		},
	)
}

type SampleQuestionsRequest struct {
}

type SampleQuestionsResponse struct {
	Questions []string `json:"questions"`
	Err       string   `json:"err,omitempty"`
}

func makeSampleQuestionsEndpoint(svc Server) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		v, err := svc.SampleQuestions()
		if err != nil {
			return SampleQuestionsResponse{v, err.Error()}, nil
		}
		return SampleQuestionsResponse{v, ""}, nil
	}
}

func GetSampleQuestionsHandler(svc Server) *httptransport.Server {
	return httptransport.NewServer(
		makeSampleQuestionsEndpoint(svc),
		func(_ context.Context, r *http.Request) (interface{}, error) {
			var request SampleQuestionsRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			return request, nil
		},
		func(_ context.Context, w http.ResponseWriter, response interface{}) error {
			return json.NewEncoder(w).Encode(response)
		},
	)
}

type AskRequest struct {
	Question string `json:"question"`
}

type AskResponse struct {
	Query   string `json:"query"`
	DataCsv string `json:"data_csv"`
	Err     string `json:"err,omitempty"`
}

func makeAskEndpoint(svc Server) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(AskRequest)
		v, err := svc.Ask(req.Question)
		if err != nil {
			return AskResponse{
				Err: err.Error(),
			}, nil
		}
		var errStr string
		if v.Error != nil {
			errStr = v.Error.Error()
		}
		return AskResponse{
			Query:   v.Query,
			DataCsv: v.DataCsv,
			Err:     errStr,
		}, nil
	}
}

func GetAskHandler(svc Server) *httptransport.Server {
	return httptransport.NewServer(
		makeAskEndpoint(svc),
		func(_ context.Context, r *http.Request) (interface{}, error) {
			var request AskRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			return request, nil
		},
		func(_ context.Context, w http.ResponseWriter, response interface{}) error {
			return json.NewEncoder(w).Encode(response)
		},
	)
}
