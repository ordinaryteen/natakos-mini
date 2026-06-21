package adapters

import (
	"database/sql"
	"fmt"
	"gateway-go/domain"
	"os"
	"strconv"

	_ "github.com/lib/pq" // Driver side-effect initialization
)

type PostgresDBAdapter struct {
	db *sql.DB
}

// NewPostgresDBAdapter establishes and validates our live SQL connection pool
func NewPostgresDBAdapter() (*PostgresDBAdapter, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct connection string matching our network environment specs
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database pool connection: %w", err)
	}

	// Verify the network line to the container vault is physically open
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database container: %w", err)
	}

	fmt.Println("🔌 [GATEWAY-GO] Secure connection pool established with natakos-postgres-db!")
	return &PostgresDBAdapter{db: db}, nil
}

// FindTenantByJID executes an indexed relational query to extract tenant truth
func (a *PostgresDBAdapter) FindTenantByJID(phoneJID string) (*domain.Tenant, error) {
	var t domain.Tenant

	// Structured SQL query performing a cross-table lookup to capture room details
	query := `
		SELECT t.name, t.phone_jid, r.room_number, t.role 
		FROM tenants t
		LEFT JOIN rooms r ON t.room_id = r.id
		WHERE t.phone_jid = $1;
	`

	row := a.db.QueryRow(query, phoneJID)
	err := row.Scan(&t.Name, &t.PhoneJID, &t.RoomNumber, &t.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found for JID: %s", phoneJID)
		}
		return nil, fmt.Errorf("database query failure: %w", err)
	}

	return &t, nil
}

// FindRoomByID queries database to extract the room truth by its ID
func (a *PostgresDBAdapter) FindRoomByID(roomID string) (*domain.Room, error) {
	var r domain.Room
	var id int
	var propertyID int

	query := `
		SELECT id, property_id, room_number, price_per_month, status 
		FROM rooms 
		WHERE id = $1;
	`
	row := a.db.QueryRow(query, roomID)
	err := row.Scan(&id, &propertyID, &r.RoomNumber, &r.MonthlyPrice, &r.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room not found for ID: %s", roomID)
		}
		return nil, fmt.Errorf("database query failure: %w", err)
	}

	r.ID = strconv.Itoa(id)
	r.PropertyID = strconv.Itoa(propertyID)
	return &r, nil
}

func (a *PostgresDBAdapter) Close() error {
	return a.db.Close()
}
