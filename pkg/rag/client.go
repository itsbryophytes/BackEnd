package rag

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	UploadDocument(ctx context.Context, req UploadRequest) (map[string]any, error)
	GetStaging(ctx context.Context, userID string, documentID string) (map[string]any, error)
	ListStaging(ctx context.Context, userID string) (map[string]any, error)
	ConfirmDocument(ctx context.Context, userID string, documentID string, body any) (map[string]any, error)
	UpdateDocument(ctx context.Context, userID string, documentID string, body any) (map[string]any, error)
	ManualDocument(ctx context.Context, userID string, body any) (map[string]any, error)
	DiscardDocument(ctx context.Context, req DiscardRequest) (map[string]any, error)
	GetPendingDocuments(ctx context.Context, userID string) (map[string]any, error)
	GetResults(ctx context.Context, userID string) (map[string]any, error)
	ListDocuments(ctx context.Context, userID string) (map[string]any, error)
	DeleteDocument(ctx context.Context, userID string, documentID string) (map[string]any, error)
	Health(ctx context.Context) error
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	UserID    string        `json:"user_id"`
	Message   string        `json:"message"`
	History   []ChatMessage `json:"history"`
	TopK      int           `json:"top_k,omitempty"`
	Threshold float64       `json:"threshold,omitempty"`
}

type ChatResponse struct {
	Reply string         `json:"reply"`
	Meta  map[string]any `json:"meta,omitempty"`
}

type UploadRequest struct {
	UserID       string
	DocumentType string
	FileName     string
	ContentType  string
	File         io.Reader
}

type DiscardRequest struct {
	UserID        string
	DocumentID    string
	RemoveFromRAG bool
}

func NewClient() Client {
	baseURL := os.Getenv("RAG_SERVICE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("PYTHON_SERVICE_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	return &HTTPClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 90 * time.Second},
	}
}

func (c *HTTPClient) Chat(ctx context.Context, reqBody ChatRequest) (*ChatResponse, error) {
	if reqBody.TopK == 0 {
		reqBody.TopK = 6
	}
	if reqBody.Threshold == 0 {
		reqBody.Threshold = 0.65
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/chat/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rag service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, decodeError(resp)
	}

	return parseChatStream(resp.Body)
}

func (c *HTTPClient) UploadDocument(ctx context.Context, upload UploadRequest) (map[string]any, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("user_id", upload.UserID); err != nil {
		return nil, err
	}
	documentType := upload.DocumentType
	if documentType == "" {
		documentType = "lab_result"
	}
	if err := writer.WriteField("document_type", documentType); err != nil {
		return nil, err
	}

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(upload.FileName)))
	if upload.ContentType != "" {
		partHeader.Set("Content-Type", upload.ContentType)
	}
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, upload.File); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/pipeline/upload", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doJSON(req)
}

func (c *HTTPClient) GetStaging(ctx context.Context, userID string, documentID string) (map[string]any, error) {
	return c.get(ctx, "/api/pipeline/staging/"+url.PathEscape(documentID), queryUser(userID))
}

func (c *HTTPClient) ListStaging(ctx context.Context, userID string) (map[string]any, error) {
	return c.get(ctx, "/api/pipeline/staging", queryUser(userID))
}

func (c *HTTPClient) ConfirmDocument(ctx context.Context, userID string, documentID string, body any) (map[string]any, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonData)
	}
	return c.post(ctx, "/api/pipeline/confirm/"+url.PathEscape(documentID), queryUser(userID), bodyReader)
}

func (c *HTTPClient) UpdateDocument(ctx context.Context, userID string, documentID string, body any) (map[string]any, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.url("/api/pipeline/results/"+url.PathEscape(documentID), queryUser(userID)), bodyReader)
	if err != nil {
		return nil, err
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.doJSON(req)
}


func (c *HTTPClient) ManualDocument(ctx context.Context, userID string, body any) (map[string]any, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Combine manual body with user_id for RAG schema
	var combined map[string]any
	if err := json.Unmarshal(bodyBytes, &combined); err != nil {
		return nil, err
	}
	combined["user_id"] = userID

	finalBytes, _ := json.Marshal(combined)
	return c.post(ctx, "/api/pipeline/manual", nil, bytes.NewReader(finalBytes))
}

func (c *HTTPClient) GetPendingDocuments(ctx context.Context, userID string) (map[string]any, error) {
	return c.get(ctx, "/api/pipeline/staging", queryUser(userID))
}

func (c *HTTPClient) DiscardDocument(ctx context.Context, discard DiscardRequest) (map[string]any, error) {
	q := queryUser(discard.UserID)
	q.Set("remove_from_rag", fmt.Sprintf("%t", discard.RemoveFromRAG))
	return c.post(ctx, "/api/pipeline/discard/"+url.PathEscape(discard.DocumentID), q, nil)
}

func (c *HTTPClient) GetResults(ctx context.Context, userID string) (map[string]any, error) {
	return c.get(ctx, "/api/pipeline/results", queryUser(userID))
}

func (c *HTTPClient) ListDocuments(ctx context.Context, userID string) (map[string]any, error) {
	return c.get(ctx, "/api/pipeline/documents", queryUser(userID))
}

func (c *HTTPClient) DeleteDocument(ctx context.Context, userID string, documentID string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.url("/api/pipeline/documents/"+url.PathEscape(documentID), queryUser(userID)), nil)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *HTTPClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/healthz", nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return decodeError(resp)
	}
	return nil
}

func (c *HTTPClient) get(ctx context.Context, path string, q url.Values) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(path, q), nil)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *HTTPClient) post(ctx context.Context, path string, q url.Values, body io.Reader) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url(path, q), body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.doJSON(req)
}

func (c *HTTPClient) doJSON(req *http.Request) (map[string]any, error) {
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rag service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, decodeError(resp)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *HTTPClient) url(path string, q url.Values) string {
	if len(q) == 0 {
		return c.baseURL + path
	}
	return c.baseURL + path + "?" + q.Encode()
}

func queryUser(userID string) url.Values {
	q := url.Values{}
	q.Set("user_id", userID)
	return q
}

func parseChatStream(r io.Reader) (*ChatResponse, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var reply strings.Builder
	var meta map[string]any

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" {
			continue
		}
		if payload == "[DONE]" {
			break
		}

		var event map[string]any
		if err := json.Unmarshal([]byte(payload), &event); err != nil {
			reply.WriteString(payload)
			continue
		}

		switch event["type"] {
		case "meta":
			if value, ok := event["meta"].(map[string]any); ok {
				meta = value
			}
		case "chunk":
			if text, ok := event["text"].(string); ok {
				reply.WriteString(text)
			}
		case "error":
			if text, ok := event["text"].(string); ok {
				return nil, errors.New(text)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ChatResponse{Reply: reply.String(), Meta: meta}, nil
}

func decodeError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		return fmt.Errorf("rag service returned status %d", resp.StatusCode)
	}
	return fmt.Errorf("rag service returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

func escapeQuotes(value string) string {
	return strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(value)
}
