package horizon

import (
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/throttled/throttled"
)

// Config is the configuration for horizon.  It get's populated by the
// app's main function and is provided to NewApp.
type Config struct {
	DatabaseURL            string
	StellarCoreDatabaseURL string
	StellarCoreURL         string
	Port                   int
	ConnectionTimeout      time.Duration
	RateLimit              *throttled.RateQuota
	RateLimitRedisKey      string
	RedisURL               string
	FriendbotURL           *url.URL
	LogLevel               logrus.Level
	LogFile                string
	SentryDSN              string
	LogglyTag              string
	LogglyToken            string
	// Maximum length of the path returned by `/paths` endpoint.
	MaxPathLength uint
	// TLSCert is a path to a certificate file to use for horizon's TLS config
	TLSCert string
	// TLSKey is the path to a private key file to use for horizon's TLS config
	TLSKey string
	// Ingest is a boolean that indicates whether or not this horizon instance
	// should run the data ingestion subsystem.
	Ingest bool
	// CursorName indentifying horizon cursor, so multiple horizons will be able
	// to communicate with a single core without collisions in cursor info
	CursorName string
	// HistoryRetentionCount represents the minimum number of ledgers worth of
	// history data to retain in the horizon database. For the purposes of
	// determining a "retention duration", each ledger roughly corresponds to 10
	// seconds of real time.
	HistoryRetentionCount uint
	// StaleThreshold represents the number of ledgers a history database may be
	// out-of-date by before horizon begins to respond with an error to history
	// requests.
	StaleThreshold uint
	// SkipCursorUpdate causes the ingestor to skip reporting the "last imported
	// ledger" state to stellar-core.
	SkipCursorUpdate bool
	// Database connection pool configuration
	HorizonDBMaxOpenConnections int
	HorizonDBMaxIdleConnections int
	CoreDBMaxOpenConnections    int
	CoreDBMaxIdleConnections    int
	// EnableAssetStats is a feature flag that determines whether to calculate
	// asset stats during the ingestion and expose `/assets` endpoint.
	// Enabling it has a negative impact on CPU when ingesting ledgers full of
	// many different assets related operations.
	EnableAssetStats bool
}
