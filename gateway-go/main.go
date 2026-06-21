package main

import (
	"fmt"
	"gateway-go/adapters"
	"gateway-go/service"
	"net/http"
	"os"
)

func main() {
	fmt.Println("🚀 Booting up Decoupled Gateway-Go Engine on Go 1.25.4...")

	// 1. Initialize Relational DB Engine Connection
	dbAdapter, err := adapters.NewPostgresDBAdapter()
	if err != nil {
		fmt.Printf("❌ Database connection critical failure: %v\n", err)
		os.Exit(1)
	}
	defer dbAdapter.Close()

	// 2. Initialize Asynchronous RabbitMQ Broker Highway Connection
	mqAdapter, err := adapters.NewRabbitMQPublisherAdapter()
	if err != nil {
		fmt.Printf("❌ Message broker pipeline connection failure: %v\n", err)
		os.Exit(1)
	}
	defer mqAdapter.Close()

	// 3. Inject both decoupled drivers directly into Core Service Layer
	coreService := service.NewRoomService(dbAdapter, mqAdapter)

	// 4. Attach Service into Inbound Border Control HTTP Adapter
	httpWebhookAdapter := adapters.NewHTTPWebhookAdapter(coreService)

	http.HandleFunc("/whatsapp/webhook", httpWebhookAdapter.HandleWebhook)

	fmt.Println("⚡ Decoupled Ingress boundary listening live on port :8085!")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		fmt.Printf("❌ Network server startup crashed: %v\n", err)
	}
}
