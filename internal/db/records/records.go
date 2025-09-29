package records

import (
	"database/sql"
	"time"
)

type Account struct {
	Id       int
	Username string
	Hostname string
	Label    sql.NullString
	Tags     sql.NullString
	Serial   int
	IsActive bool
}

type PublicKey struct {
	Id        int
	Algorithm string
	KeyData   string
	Comment   string
	IsGlobal  bool
}

type SystemKey struct {
	Id         int
	Serial     int
	PublicKey  string
	PrivateKey string
	IsActive   bool
}

type AuditLogEntry struct {
	Id        int
	Timestamp string
	Username  string
	Action    string
	Details   string
}

type BootstrapSession struct {
	Id            string
	Username      string
	Hostname      string
	Label         string
	Tags          string
	TempPublicKey string
	CreatedAt     time.Time
	ExpiresAt     time.Time
	Status        string
}
