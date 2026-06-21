package service

import (
	"errors"
	"fmt"
	"gateway-go/adapters"
	"gateway-go/ports"
	"strings"
	"time"
)

type RoomService struct {
	dbRepo    ports.OutboundPort
	mqAdapter *adapters.RabbitMQPublisherAdapter
}

// Update constructor to take our newly added MQ adapter resource block
func NewRoomService(repo ports.OutboundPort, mq *adapters.RabbitMQPublisherAdapter) *RoomService {
	return &RoomService{dbRepo: repo, mqAdapter: mq}
}

func (s *RoomService) ProcessIncomingMessage(senderJID, rawText, correlationID string) (string, error) {
	if strings.TrimSpace(rawText) == "" {
		return "", errors.New("empty text body received")
	}

	// 1. Identity lookup against live Postgres container record
	tenant, err := s.dbRepo.FindTenantByJID(senderJID)
	if err != nil {
		return "Maaf Kak, nomor kamu belum terdaftar di sistem Kos Natakos.", nil
	}

	generatedMsgID := fmt.Sprintf("MSG-%d", time.Now().UnixNano())

	// 2. ASYNCHRONOUS SHIFT: Drop the packet payload into the queue instead of waiting on HTTP calls
	err = s.mqAdapter.PublishInteractionToQueue(*tenant, rawText, generatedMsgID, correlationID)
	if err != nil {
		fmt.Printf("🚨 [GATEWAY-GO] Queue submission failure: %v\n", err)
		return "Waduh, antrean pemrosesan kami sedang padat. Sila coba beberapa saat lagi.", nil
	}

	// 3. Instant acknowledgment strategy!
	// We release the user immediately while the worker runs behind the scenes.
	return "Pesan kamu sudah diterima dan sedang diproses oleh sistem AI Natakos, mohon tunggu sebentar ya Kak!", nil
}
