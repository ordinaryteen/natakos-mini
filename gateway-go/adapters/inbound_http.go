package adapters

import (
	"encoding/json"
	"fmt"
	"gateway-go/ports"
	"io"
	"net/http"
	"strings"
	"time"
)

type WebhookRequestPayload struct {
	MessageID string `json:"message_id"`
	SenderJID string `json:"sender_jid"`
	TextBody  string `json:"text_body"`
}

type ContractResponsePayload struct {
	Status         string `json:"status"`
	IntentDetected string `json:"intent_detected"`
	ReplyMessage   string `json:"reply_message"`
}

// Define standard JSON log structure
type JsonLog struct {
	Timestamp     string `json:"timestamp"`
	Level         string `json:"level"`
	Service       string `json:"service"`
	CorrelationID string `json:"correlation_id"`
	Message       string `json:"message"`
}

func logStructured(level, correlationID, msg string) {
	logObj := JsonLog{
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Level:         level,
		Service:       "gateway-go",
		CorrelationID: correlationID,
		Message:       msg,
	}
	jsonBytes, _ := json.Marshal(logObj)
	fmt.Println(string(jsonBytes))
}

type HTTPWebhookAdapter struct {
	coreService ports.InboundPort // Abstract plug to our business brain
}

func NewHTTPWebhookAdapter(service ports.InboundPort) *HTTPWebhookAdapter {
	return &HTTPWebhookAdapter{coreService: service}
}

func (h *HTTPWebhookAdapter) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("X-Natakos-Token") != "surabaya_secret_2026" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 1. Generate Global Tracking Identifier
	correlationID := fmt.Sprintf("req_%d", time.Now().UnixNano())
	logStructured("INFO", correlationID, "Received inbound webhook packet from edge channel")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logStructured("ERROR", correlationID, fmt.Sprintf("Failed to read body: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req WebhookRequestPayload
	if err := json.Unmarshal(body, &req); err != nil {
		logStructured("ERROR", correlationID, fmt.Sprintf("Failed to unmarshal JSON payload: %v", err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if strings.TrimSpace(req.TextBody) == "" {
		logStructured("WARN", correlationID, "Empty text body received")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "empty text body received"})
		return
	}

	// Push the parsed transit items directly into the abstract inbound core gate
	replyMsg, err := h.coreService.ProcessIncomingMessage(req.SenderJID, req.TextBody, correlationID)
	if err != nil {
		logStructured("ERROR", correlationID, fmt.Sprintf("ProcessIncomingMessage error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logStructured("INFO", correlationID, "Inbound webhook processed successfully")

	// Format response payload to match our signed API Execution Contract exactly
	responseContract := ContractResponsePayload{
		Status:         "SUCCESS",
		IntentDetected: "RESOLVE_ROUTE",
		ReplyMessage:   replyMsg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseContract)
}
