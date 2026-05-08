package dto

const (
	MESSAGE_SEND_CHAT_SUCCESS      = "Successfully sent message"
	MESSAGE_GET_HISTORY_SUCCESS    = "Successfully retrieved chat history"
	MESSAGE_GET_SESSIONS_SUCCESS   = "Successfully retrieved chat sessions"
	MESSAGE_DELETE_SESSION_SUCCESS = "Successfully deleted session"
)

type SendMessageRequest struct {
	Message   string        `json:"message" binding:"required"`
	SessionID string        `json:"session_id,omitempty"`
	History   []ChatMessage `json:"history,omitempty"`
	TopK      int           `json:"top_k,omitempty"`
	Threshold float64       `json:"threshold,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SendMessageResponse struct {
	Reply     string        `json:"reply"`
	SessionID string        `json:"session_id"`
	History   []ChatMessage `json:"history,omitempty"`
	Meta      any           `json:"meta,omitempty"`
}

type ChatSessionResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
}

type ChatHistoryResponse struct {
	SessionID string        `json:"session_id"`
	Messages  []ChatMessage `json:"messages"`
}

type PythonChatRequest struct {
	UserID    string        `json:"user_id"`
	Message   string        `json:"message"`
	History   []ChatMessage `json:"history"`
	TopK      int           `json:"top_k,omitempty"`
	Threshold float64       `json:"threshold,omitempty"`
}

type PythonChatResponse struct {
	Reply string `json:"reply"`
	Meta  any    `json:"meta,omitempty"`
}
