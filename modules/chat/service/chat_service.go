package service

import (
	"context"
	"fmt"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/modules/chat/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/chat/repository"
	"github.com/Caknoooo/go-gin-clean-starter/pkg/rag"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService interface {
	SendMessage(ctx context.Context, userID string, req dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	GetSessions(ctx context.Context, userID string) ([]dto.ChatSessionResponse, error)
	GetHistory(ctx context.Context, userID string, sessionID string) (*dto.ChatHistoryResponse, error)
	DeleteSession(ctx context.Context, userID string, sessionID string) error
	ClearAllHistory(ctx context.Context, userID string) error
}

type chatService struct {
	ragClient rag.Client
	repo      repository.ChatRepository
	db        *gorm.DB
}

func NewChatService(ragClient rag.Client, repo repository.ChatRepository, db *gorm.DB) ChatService {
	return &chatService{
		ragClient: ragClient,
		repo:      repo,
		db:        db,
	}
}

func (s *chatService) SendMessage(ctx context.Context, userID string, req dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	sessionID, err := s.resolveSession(ctx, userID, req.SessionID, req.Message)
	if err != nil {
		return nil, err
	}

	dbMessages, _ := s.repo.GetRecentMessages(ctx, s.db, sessionID, 20)
	ragHistory := make([]rag.ChatMessage, 0, len(dbMessages))
	for _, m := range dbMessages {
		ragHistory = append(ragHistory, rag.ChatMessage{Role: m.Role, Content: m.Content})
	}
	if len(req.History) > 0 {
		ragHistory = make([]rag.ChatMessage, 0, len(req.History))
		for _, m := range req.History {
			ragHistory = append(ragHistory, rag.ChatMessage{Role: m.Role, Content: m.Content})
		}
	}

	resp, err := s.ragClient.Chat(ctx, rag.ChatRequest{
		UserID:    userID,
		Message:   req.Message,
		History:   ragHistory,
		TopK:      req.TopK,
		Threshold: req.Threshold,
	})
	if err != nil {
		return nil, err
	}

	sessionUUID, _ := uuid.Parse(sessionID)
	s.repo.AddMessage(ctx, s.db, entities.ChatMessage{
		SessionID: sessionUUID,
		Role:      "user",
		Content:   req.Message,
	})
	s.repo.AddMessage(ctx, s.db, entities.ChatMessage{
		SessionID: sessionUUID,
		Role:      "assistant",
		Content:   resp.Reply,
	})

	updatedMessages, _ := s.repo.GetRecentMessages(ctx, s.db, sessionID, 20)
	respHistory := make([]dto.ChatMessage, 0, len(updatedMessages))
	for _, m := range updatedMessages {
		respHistory = append(respHistory, dto.ChatMessage{Role: m.Role, Content: m.Content})
	}

	return &dto.SendMessageResponse{
		Reply:     resp.Reply,
		SessionID: sessionID,
		History:   respHistory,
		Meta:      resp.Meta,
	}, nil
}

func (s *chatService) GetSessions(ctx context.Context, userID string) ([]dto.ChatSessionResponse, error) {
	sessions, err := s.repo.GetSessionsByUser(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ChatSessionResponse, 0, len(sessions))
	for _, sess := range sessions {
		result = append(result, dto.ChatSessionResponse{
			ID:        sess.ID.String(),
			Title:     sess.Title,
			UpdatedAt: sess.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

func (s *chatService) GetHistory(ctx context.Context, userID string, sessionID string) (*dto.ChatHistoryResponse, error) {
	session, err := s.repo.GetSessionWithMessages(ctx, s.db, sessionID, userID)
	if err != nil {
		return nil, err
	}
	messages := make([]dto.ChatMessage, 0, len(session.Messages))
	for _, m := range session.Messages {
		messages = append(messages, dto.ChatMessage{Role: m.Role, Content: m.Content})
	}
	return &dto.ChatHistoryResponse{
		SessionID: session.ID.String(),
		Messages:  messages,
	}, nil
}

func (s *chatService) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	return s.repo.DeleteSession(ctx, s.db, sessionID, userID)
}

func (s *chatService) ClearAllHistory(ctx context.Context, userID string) error {
	return s.repo.ClearAllSessions(ctx, s.db, userID)
}

func (s *chatService) resolveSession(ctx context.Context, userID string, sessionID string, firstMessage string) (string, error) {
	if sessionID != "" {
		return sessionID, nil
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}
	title := firstMessage
	if len(title) > 60 {
		title = fmt.Sprintf("%s...", title[:57])
	}
	session, err := s.repo.CreateSession(ctx, s.db, entities.ChatSession{
		UserID: parsedUserID,
		Title:  title,
	})
	if err != nil {
		return "", err
	}
	return session.ID.String(), nil
}
