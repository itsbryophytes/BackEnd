package repository

import (
	"context"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateSession(ctx context.Context, db *gorm.DB, session entities.ChatSession) (entities.ChatSession, error)
	GetSessionsByUser(ctx context.Context, db *gorm.DB, userID string) ([]entities.ChatSession, error)
	GetSessionWithMessages(ctx context.Context, db *gorm.DB, sessionID string, userID string) (entities.ChatSession, error)
	AddMessage(ctx context.Context, db *gorm.DB, message entities.ChatMessage) (entities.ChatMessage, error)
	GetRecentMessages(ctx context.Context, db *gorm.DB, sessionID string, limit int) ([]entities.ChatMessage, error)
	DeleteSession(ctx context.Context, db *gorm.DB, sessionID string, userID string) error
	ClearAllSessions(ctx context.Context, db *gorm.DB, userID string) error
}

type chatRepository struct{}

func NewChatRepository() ChatRepository {
	return &chatRepository{}
}

func (r *chatRepository) CreateSession(ctx context.Context, db *gorm.DB, session entities.ChatSession) (entities.ChatSession, error) {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	err := db.WithContext(ctx).Create(&session).Error
	return session, err
}

func (r *chatRepository) GetSessionsByUser(ctx context.Context, db *gorm.DB, userID string) ([]entities.ChatSession, error) {
	var sessions []entities.ChatSession
	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *chatRepository) GetSessionWithMessages(ctx context.Context, db *gorm.DB, sessionID string, userID string) (entities.ChatSession, error) {
	var session entities.ChatSession
	err := db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		First(&session).Error
	return session, err
}

func (r *chatRepository) AddMessage(ctx context.Context, db *gorm.DB, message entities.ChatMessage) (entities.ChatMessage, error) {
	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}
	err := db.WithContext(ctx).Create(&message).Error
	if err != nil {
		return entities.ChatMessage{}, err
	}
	// Touch session updated_at
	db.WithContext(ctx).Model(&entities.ChatSession{}).
		Where("id = ?", message.SessionID).
		Update("updated_at", gorm.Expr("NOW()"))
	return message, nil
}

func (r *chatRepository) GetRecentMessages(ctx context.Context, db *gorm.DB, sessionID string, limit int) ([]entities.ChatMessage, error) {
	var messages []entities.ChatMessage
	err := db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *chatRepository) DeleteSession(ctx context.Context, db *gorm.DB, sessionID string, userID string) error {
	return db.WithContext(ctx).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Delete(&entities.ChatSession{}).Error
}

func (r *chatRepository) ClearAllSessions(ctx context.Context, db *gorm.DB, userID string) error {
	return db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.ChatSession{}).Error
}
