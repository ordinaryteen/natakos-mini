package ports

// InboundPort represents the "Door".
// It defines the contract for any adapter (HTTP, CLI, Queue) that wants to feed data into our business logic.
type InboundPort interface {
	ProcessIncomingMessage(senderJID, rawText, correlationID string) (string, error)
}
