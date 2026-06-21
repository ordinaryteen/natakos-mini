package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway-go/domain"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisherAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQPublisherAdapter sets up a persistent connection to the message highway
func NewRabbitMQPublisherAdapter() (*RabbitMQPublisherAdapter, error) {
	amqpURL := os.Getenv("AMQP_SERVER_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	// Connect to the RabbitMQ container server broker
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ broker: %w", err)
	}

	// Open an isolated virtual channel for data transit operations
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open AMQP communication channel: %w", err)
	}

	// Declare our reliable queue buffer room
	q, err := ch.QueueDeclare(
		"interaction.queue", // Name of the queue highway
		true,                // Durable (survives container crashes/reboots)
		false,               // Delete when unused
		false,               // Exclusive
		false,               // No-wait
		nil,                 // Arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue boundary line: %w", err)
	}

	fmt.Println("📟 [GATEWAY-GO] Asynchronous AMQP Event Highway Connected & Queue Declared!")
	return &RabbitMQPublisherAdapter{conn: conn, channel: ch, queue: q}, nil
}

// PublishInteractionToQueue drops the packed interaction data right onto the event buffer line
func (p *RabbitMQPublisherAdapter) PublishInteractionToQueue(tenant domain.Tenant, rawText, msgID, correlationID string) error {
	// Re-map into our signed JSON transport payload contract block
	payload := map[string]string{
		"message_id":  msgID,
		"sender_jid":  tenant.PhoneJID,
		"raw_text":    rawText,
		"tenant_name": tenant.Name,
		"room_number": tenant.RoomNumber,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal outbound event schema: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Slam the packet onto the queue channel
	err = p.channel.PublishWithContext(ctx,
		"",           // Exchange name (default direct exchange)
		p.queue.Name, // Routing Key matches our Queue name
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Keeps message safe on disk inside the broker box!
			Headers: amqp.Table{
				"X-User-Role":      tenant.Role, // Injected structural identity header for RBAC enforcement
				"X-Correlation-ID": correlationID,
			},
			Body: jsonBytes,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to transmit packet payload to queue: %w", err)
	}

	logStructured("INFO", correlationID, fmt.Sprintf("Successfully pushed event for '%s' to RabbitMQ queue buffer", tenant.Name))
	return nil
}

func (p *RabbitMQPublisherAdapter) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
