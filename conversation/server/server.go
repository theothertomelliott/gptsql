package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"

	"github.com/theothertomelliott/gptsql/conversation"
	"github.com/theothertomelliott/gptsql/schema"
)

type ConversationID string

var ErrConversationNotFound = fmt.Errorf("conversation not found")

type Server interface {
	NewConversation() (ConversationID, error)
	SampleQuestions(ConversationID) ([]string, error)
	Ask(cid ConversationID, question string) (*conversation.Response, error)

	// TODO: Allow editing the SQL for a given question
}

type conversationServer struct {
	conversations map[ConversationID]*conversation.Conversation

	client *openai.Client
	db     *sql.DB
	dbType string
}

func New(client *openai.Client, db *sql.DB, dbType string) Server {
	return &conversationServer{
		conversations: make(map[ConversationID]*conversation.Conversation),

		client: client,
		db:     db,
		dbType: dbType,
	}
}

func (s *conversationServer) NewConversation() (ConversationID, error) {
	cid := ConversationID(uuid.New().String())

	schema, err := schema.Load(s.dbType, s.db)
	if err != nil {
		log.Fatal(err)
	}

	s.conversations[cid] = conversation.New(s.client, s.db, s.dbType, schema)
	return cid, nil
}

type NewConversationRequest struct {
}

type NewConversationResponse struct {
	ConversationID string `json:"conversation_id"`
	Err            string `json:"err,omitempty"`
}

func makeNewConversationEndpoint(svc Server) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		v, err := svc.NewConversation()
		if err != nil {
			return NewConversationResponse{
				Err: err.Error(),
			}, nil
		}
		return NewConversationResponse{
			ConversationID: string(v),
		}, nil
	}
}

func GetNewConversationHandler(svc Server) *httptransport.Server {
	return httptransport.NewServer(
		makeNewConversationEndpoint(svc),
		func(_ context.Context, r *http.Request) (interface{}, error) {
			return NewConversationRequest{}, nil
		},
		func(_ context.Context, w http.ResponseWriter, response interface{}) error {
			return json.NewEncoder(w).Encode(response)
		},
	)
}

func (s *conversationServer) SampleQuestions(cid ConversationID) ([]string, error) {
	conv, ok := s.conversations[cid]
	if !ok {
		return nil, ErrConversationNotFound
	}
	return conv.SampleQuestions()
}

func (s *conversationServer) Ask(cid ConversationID, question string) (*conversation.Response, error) {
	conv, ok := s.conversations[cid]
	if !ok {
		return nil, ErrConversationNotFound
	}
	return conv.Ask(
		conversation.Request{
			Question: question,
		},
	)
}

type SampleQuestionsRequest struct {
	ConversationID string `json:"conversation_id"`
}

type SampleQuestionsResponse struct {
	Questions []string `json:"questions"`
	Err       string   `json:"err,omitempty"`
}

func makeSampleQuestionsEndpoint(svc Server) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(SampleQuestionsRequest)
		v, err := svc.SampleQuestions(ConversationID(req.ConversationID))
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
	ConversationID string `json:"conversation_id"`
	Question       string `json:"question"`
}

type AskResponse struct {
	Query   string `json:"query"`
	DataCsv string `json:"data_csv"`
	Err     string `json:"err,omitempty"`
}

func makeAskEndpoint(svc Server) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(AskRequest)
		v, err := svc.Ask(ConversationID(req.ConversationID), req.Question)
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
