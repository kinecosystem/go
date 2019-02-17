package main

import (
	"go/types"
	stdLog "log"
	"os"

	"github.com/kinecosystem/go/network"
	"github.com/kinecosystem/go/services/horizon/internal"
	"github.com/kinecosystem/go/support/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
<<<<<<< HEAD
=======
	horizon "github.com/stellar/go/services/horizon/internal"
	"github.com/stellar/go/services/horizon/internal/db2/schema"
	apkg "github.com/stellar/go/support/app"
	support "github.com/stellar/go/support/config"
	"github.com/stellar/go/support/log"
>>>>>>> horizon-v0.16.0
	"github.com/throttled/throttled"
)

var app *horizon.App
var config horizon.Config

var rootCmd *cobra.Command

// validateBothOrNeither ensures that both options are provided, if either is provided
func validateBothOrNeither(option1, option2 string) {
	arg1, arg2 := viper.GetString(option1), viper.GetString(option2)
	switch {
	case arg1 != "" && arg2 == "":
		stdLog.Fatalf("Invalid config: %s = %s, but corresponding option %s is not configured", option1, arg1, option2)
	case arg1 == "" && arg2 != "":
		stdLog.Fatalf("Invalid config: %s = %s, but corresponding option %s is not configured", option2, arg2, option1)
	}
}

// checkMigrations looks for necessary database migrations and fails with a descriptive error if migrations are needed
func checkMigrations() {
	migrationsToApplyUp := schema.GetMigrationsUp(viper.GetString("db-url"))
	if len(migrationsToApplyUp) > 0 {
		stdLog.Printf(`There are %v migrations to apply in the "up" direction.`, len(migrationsToApplyUp))
		stdLog.Printf("The necessary migrations are: %v", migrationsToApplyUp)
		stdLog.Printf("A database migration is required to run this version (%v) of Horizon. Run \"horizon db migrate up\" to update your DB. Consult the Changelog (https://github.com/stellar/horizon/blob/master/CHANGELOG.md) for more information.", apkg.Version())
		os.Exit(1)
	}

	nMigrationsDown := schema.GetNumMigrationsDown(viper.GetString("db-url"))
	if nMigrationsDown > 0 {
		stdLog.Printf("A database migration DOWN to an earlier version of the schema is required to run this version (%v) of Horizon. Consult the Changelog (https://github.com/stellar/horizon/blob/master/CHANGELOG.md) for more information.", apkg.Version())
		stdLog.Printf("In order to migrate the database DOWN, using the HIGHEST version number of Horizon you have installed (not this binary), run \"horizon db migrate down %v\".", nMigrationsDown)
		os.Exit(1)
	}
}

// configOpts defines the complete flag configuration for horizon
// Add a new entry here to connect a new field in the horizon.Config struct
var configOpts = []*support.ConfigOption{
	&support.ConfigOption{
		Name:      "db-url",
		EnvVar:    "DATABASE_URL",
		ConfigKey: &config.DatabaseURL,
		OptType:   types.String,
		Required:  true,
		Usage:     "horizon postgres database to connect with",
	},
	&support.ConfigOption{
		Name:      "stellar-core-db-url",
		EnvVar:    "STELLAR_CORE_DATABASE_URL",
		ConfigKey: &config.StellarCoreDatabaseURL,
		OptType:   types.String,
		Required:  true,
		Usage:     "stellar-core postgres database to connect with",
	},
	&support.ConfigOption{
		Name:      "stellar-core-url",
		ConfigKey: &config.StellarCoreURL,
		OptType:   types.String,
		Required:  true,
		Usage:     "stellar-core to connect with (for http commands)",
	},
	&support.ConfigOption{
		Name:        "port",
		ConfigKey:   &config.Port,
		OptType:     types.Uint,
		FlagDefault: uint(8000),
		Usage:       "tcp port to listen on for http requests",
	},
	&support.ConfigOption{
		Name:        "max-db-connections",
		ConfigKey:   &config.MaxDBConnections,
		OptType:     types.Int,
		FlagDefault: 20,
		Usage:       "max db connections (per DB), may need to be increased when responses are slow but DB CPU is normal",
	},
	&support.ConfigOption{
		Name:           "sse-update-frequency",
		ConfigKey:      &config.SSEUpdateFrequency,
		OptType:        types.Int,
		FlagDefault:    5,
		CustomSetValue: support.SetDuration,
		Usage:          "defines how often streams should check if there's a new ledger (in seconds), may need to increase in case of big number of streams",
	},
	&support.ConfigOption{
		Name:           "connection-timeout",
		ConfigKey:      &config.ConnectionTimeout,
		OptType:        types.Int,
		FlagDefault:    55,
		CustomSetValue: support.SetDuration,
		Usage:          "defines the timeout of connection after which 504 response will be sent or stream will be closed, if Horizon is behind a load balancer with idle connection timeout, this should be set to a few seconds less that idle timeout",
	},
	&support.ConfigOption{
		Name:        "per-hour-rate-limit",
		ConfigKey:   &config.RateLimit,
		OptType:     types.Int,
		FlagDefault: 3600,
		CustomSetValue: func(co *support.ConfigOption) {
			var rateLimit *throttled.RateQuota = nil
			perHourRateLimit := viper.GetInt(co.Name)
			if perHourRateLimit != 0 {
				rateLimit = &throttled.RateQuota{
					MaxRate:  throttled.PerHour(perHourRateLimit),
					MaxBurst: 100,
				}
				*(co.ConfigKey.(**throttled.RateQuota)) = rateLimit
			}
		},
		Usage: "max count of requests allowed in a one hour period, by remote ip address",
	},
	&support.ConfigOption{
		Name:      "rate-limit-redis-key",
		ConfigKey: &config.RateLimitRedisKey,
		OptType:   types.String,
		Usage:     "redis key for storing rate limit data, useful when deploying a cluster of Horizons, ignored when redis-url is empty",
	},
	&support.ConfigOption{
		Name:      "redis-url",
		ConfigKey: &config.RedisURL,
		OptType:   types.String,
		Usage:     "redis to connect with, for rate limiting",
	},
	&support.ConfigOption{
		Name:           "friendbot-url",
		ConfigKey:      &config.FriendbotURL,
		OptType:        types.String,
		CustomSetValue: support.SetURL,
		Usage:          "friendbot service to redirect to",
	},
	&support.ConfigOption{
		Name:        "log-level",
		ConfigKey:   &config.LogLevel,
		OptType:     types.String,
		FlagDefault: "info",
		CustomSetValue: func(co *support.ConfigOption) {
			ll, err := logrus.ParseLevel(viper.GetString(co.Name))
			if err != nil {
				stdLog.Fatalf("Could not parse log-level: %v", viper.GetString(co.Name))
			}
			*(co.ConfigKey.(*logrus.Level)) = ll
		},
		Usage: "minimum log severity (debug, info, warn, error) to log",
	},
	&support.ConfigOption{
		Name:      "log-file",
		ConfigKey: &config.LogFile,
		OptType:   types.String,
		Usage:     "name of the file where logs will be saved (leave empty to send logs to stdout)",
	},
	&support.ConfigOption{
		Name:        "max-path-length",
		ConfigKey:   &config.MaxPathLength,
		OptType:     types.Uint,
		FlagDefault: uint(4),
		Usage:       "the maximum number of assets on the path in `/paths` endpoint",
	},
	&support.ConfigOption{
		Name:      "network-passphrase",
		ConfigKey: &config.NetworkPassphrase,
		OptType:   types.String,
		Required:  true,
		Usage:     "Override the network passphrase",
	},
	&support.ConfigOption{
		Name:      "sentry-dsn",
		ConfigKey: &config.SentryDSN,
		OptType:   types.String,
		Usage:     "Sentry URL to which panics and errors should be reported",
	},
	&support.ConfigOption{
		Name:      "loggly-token",
		ConfigKey: &config.LogglyToken,
		OptType:   types.String,
		Usage:     "Loggly token, used to configure log forwarding to loggly",
	},
	&support.ConfigOption{
		Name:        "loggly-tag",
		ConfigKey:   &config.LogglyTag,
		OptType:     types.String,
		FlagDefault: "horizon",
		Usage:       "Tag to be added to every loggly log event",
	},
	&support.ConfigOption{
		Name:      "tls-cert",
		ConfigKey: &config.TLSCert,
		OptType:   types.String,
		Usage:     "TLS certificate file to use for securing connections to horizon",
	},
	&support.ConfigOption{
		Name:      "tls-key",
		ConfigKey: &config.TLSKey,
		OptType:   types.String,
		Usage:     "TLS private key file to use for securing connections to horizon",
	},
	&support.ConfigOption{
		Name:        "ingest",
		ConfigKey:   &config.Ingest,
		OptType:     types.Bool,
		FlagDefault: false,
		Usage:       "causes this horizon process to ingest data from stellar-core into horizon's db",
	},
	&support.ConfigOption{
		Name:        "history-retention-count",
		ConfigKey:   &config.HistoryRetentionCount,
		OptType:     types.Uint,
		FlagDefault: uint(0),
		Usage:       "the minimum number of ledgers to maintain within horizon's history tables.  0 signifies an unlimited number of ledgers will be retained",
	},
	&support.ConfigOption{
		Name:        "history-stale-threshold",
		ConfigKey:   &config.StaleThreshold,
		OptType:     types.Uint,
		FlagDefault: uint(0),
		Usage:       "the maximum number of ledgers the history db is allowed to be out of date from the connected stellar-core db before horizon considers history stale",
	},
	&support.ConfigOption{
		Name:        "skip-cursor-update",
		ConfigKey:   &config.SkipCursorUpdate,
		OptType:     types.Bool,
		FlagDefault: false,
		Usage:       "causes the ingester to skip reporting the last imported ledger state to stellar-core",
	},
	&support.ConfigOption{
		Name:        "enable-asset-stats",
		ConfigKey:   &config.EnableAssetStats,
		OptType:     types.Bool,
		FlagDefault: false,
		Usage:       "enables asset stats during the ingestion and expose `/assets` endpoint, Enabling it has a negative impact on CPU",
	},
}

func main() {
	rootCmd.Execute()
}

func init() {
<<<<<<< HEAD
	viper.SetDefault("port", 8000)
	viper.SetDefault("history-retention-count", 0)
	viper.SetDefault("horizon-db-max-open-connections", 12)
	viper.SetDefault("horizon-db-max-idle-connections", 4)
	viper.SetDefault("core-db-max-open-connections", 12)
	viper.SetDefault("core-db-max-idle-connections", 4)
	viper.SetDefault("cursor-name", "HORIZON")

	viper.BindEnv("port", "PORT")
	viper.BindEnv("db-url", "DATABASE_URL")
	viper.BindEnv("stellar-core-db-url", "STELLAR_CORE_DATABASE_URL")
	viper.BindEnv("stellar-core-url", "STELLAR_CORE_URL")
	viper.BindEnv("connection-timeout", "CONNECTION_TIMEOUT")
	viper.BindEnv("per-hour-rate-limit", "PER_HOUR_RATE_LIMIT")
	viper.BindEnv("rate-limit-redis-key", "RATE_LIMIT_REDIS_KEY")
	viper.BindEnv("redis-url", "REDIS_URL")
	viper.BindEnv("ruby-horizon-url", "RUBY_HORIZON_URL")
	viper.BindEnv("friendbot-url", "FRIENDBOT_URL")
	viper.BindEnv("log-level", "LOG_LEVEL")
	viper.BindEnv("log-file", "LOG_FILE")
	viper.BindEnv("sentry-dsn", "SENTRY_DSN")
	viper.BindEnv("loggly-token", "LOGGLY_TOKEN")
	viper.BindEnv("loggly-tag", "LOGGLY_TAG")
	viper.BindEnv("tls-cert", "TLS_CERT")
	viper.BindEnv("tls-key", "TLS_KEY")
	viper.BindEnv("ingest", "INGEST")
	viper.BindEnv("cursor-name", "CURSOR_NAME")
	viper.BindEnv("network-passphrase", "NETWORK_PASSPHRASE")
	viper.BindEnv("history-retention-count", "HISTORY_RETENTION_COUNT")
	viper.BindEnv("history-stale-threshold", "HISTORY_STALE_THRESHOLD")
	viper.BindEnv("skip-cursor-update", "SKIP_CURSOR_UPDATE")
	viper.BindEnv("horizon-db-max-open-connections", "HORIZON_DB_MAX_OPEN_CONNECTIONS")
	viper.BindEnv("horizon-db-max-idle-connections", "HORIZON_DB_MAX_IDLE_CONNECTIONS")
	viper.BindEnv("core-db-max-open-connections", "CORE_DB_MAX_OPEN_CONNECTIONS")
	viper.BindEnv("core-db-max-idle-connections", "CORE_DB_MAX_IDLE_CONNECTIONS")
	viper.BindEnv("enable-asset-stats", "ENABLE_ASSET_STATS")
	viper.BindEnv("max-path-length", "MAX_PATH_LENGTH")

=======
>>>>>>> horizon-v0.16.0
	rootCmd = &cobra.Command{
		Use:   "horizon",
		Short: "client-facing api server for the stellar network",
		Long:  "client-facing api server for the stellar network",
		Run: func(cmd *cobra.Command, args []string) {
			initApp(cmd, args)
			app.Serve()
		},
	}

<<<<<<< HEAD
	rootCmd.PersistentFlags().String(
		"db-url",
		"",
		"horizon postgres database to connect with",
	)

	rootCmd.PersistentFlags().String(
		"stellar-core-db-url",
		"",
		"stellar-core postgres database to connect with",
	)

	rootCmd.PersistentFlags().String(
		"stellar-core-url",
		"",
		"stellar-core to connect with (for http commands)",
	)

	rootCmd.PersistentFlags().Int(
		"port",
		8000,
		"tcp port to listen on for http requests",
	)

	rootCmd.PersistentFlags().Int(
		"per-hour-rate-limit",
		3600,
		"max count of requests allowed in a one hour period, by remote ip address",
	)

	rootCmd.PersistentFlags().Int(
		"connection-timeout",
		55,
		"defines the timeout of connection after which 504 response will be sent or stream will be closed, if Horizon is behind a load balancer with idle connection timeout, this should be set to a few seconds less that idle timeout",
	)

	rootCmd.PersistentFlags().String(
		"rate-limit-redis-key",
		"",
		"redis key for storing rate limit data, useful when deploying a cluster of Horizons, ignored when redis-url is empty",
	)

	rootCmd.PersistentFlags().String(
		"redis-url",
		"",
		"redis to connect with, for rate limiting",
	)

	rootCmd.PersistentFlags().String(
		"friendbot-url",
		"",
		"friendbot service to redirect to",
	)

	rootCmd.PersistentFlags().String(
		"log-level",
		"info",
		"Minimum log severity (debug, info, warn, error) to log",
	)

	rootCmd.PersistentFlags().String(
		"log-file",
		"",
		"Name of the file where logs will be saved (leave empty to send logs to stdout)",
	)

	rootCmd.PersistentFlags().String(
		"sentry-dsn",
		"",
		"Sentry URL to which panics and errors should be reported",
	)

	rootCmd.PersistentFlags().String(
		"loggly-token",
		"",
		"Loggly token, used to configure log forwarding to loggly",
	)

	rootCmd.PersistentFlags().String(
		"loggly-tag",
		"horizon",
		"Tag to be added to every loggly log event",
	)

	rootCmd.PersistentFlags().String(
		"tls-cert",
		"",
		"The TLS certificate file to use for securing connections to horizon",
	)

	rootCmd.PersistentFlags().String(
		"tls-key",
		"",
		"The TLS private key file to use for securing connections to horizon",
	)

	rootCmd.PersistentFlags().Bool(
		"ingest",
		false,
		"causes this horizon process to ingest data from stellar-core into horizon's db",
	)

	rootCmd.PersistentFlags().String(
		"cursor-name",
		"HORIZON",
		"Set the cursor name used in stellar-core",
	)

	rootCmd.PersistentFlags().String(
		"network-passphrase",
		network.TestNetworkPassphrase,
		"Override the network passphrase",
	)

	rootCmd.PersistentFlags().Uint(
		"history-retention-count",
		0,
		"the minimum number of ledgers to maintain within horizon's history tables.  0 signifies an unlimited number of ledgers will be retained",
	)

	rootCmd.PersistentFlags().Uint(
		"history-stale-threshold",
		0,
		"the maximum number of ledgers the history db is allowed to be out of date from the connected stellar-core db before horizon considers history stale",
	)

	rootCmd.PersistentFlags().Bool(
		"enable-asset-stats",
		false,
		"enables asset stats during the ingestion and expose `/assets` endpoint,  Enabling it has a negative impact on CPU",
	)

	rootCmd.PersistentFlags().Uint(
		"max-path-length",
		4,
		"the maximum number of assets on the path in `/paths` endpoint",
	)
=======
	for _, co := range configOpts {
		err := co.Init(rootCmd)
		if err != nil {
			stdLog.Fatal(err.Error())
		}
	}
>>>>>>> horizon-v0.16.0

	rootCmd.AddCommand(dbCmd)
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initApp(cmd *cobra.Command, args []string) *horizon.App {
	initConfig()

	var err error
	app, err = horizon.NewApp(config)

	if err != nil {
		stdLog.Fatal(err.Error())
	}

	return app
}

func initConfig() {
	// Verify required options and load the config struct
	for _, co := range configOpts {
		co.Require()
		co.SetValue()
	}

	// Migrations should be checked as early as possible
	checkMigrations()

	// Validate options that should be provided together
	validateBothOrNeither("tls-cert", "tls-key")
	validateBothOrNeither("rate-limit-redis-key", "redis-url")

	// Configure log file
	if config.LogFile != "" {
		logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			log.DefaultLogger.Logger.Out = logFile
		} else {
			stdLog.Fatalf("Failed to open file to log: %s", err)
		}
	}

<<<<<<< HEAD
	cert, key := viper.GetString("tls-cert"), viper.GetString("tls-key")

	switch {
	case cert != "" && key == "":
		stdLog.Fatal("Invalid TLS config: key not configured")
	case cert == "" && key != "":
		stdLog.Fatal("Invalid TLS config: cert not configured")
	}

	var friendbotURL *url.URL
	friendbotURLString := viper.GetString("friendbot-url")
	if friendbotURLString != "" {
		friendbotURL, err = url.Parse(friendbotURLString)
		if err != nil {
			stdLog.Fatalf("Unable to parse URL: %s/%v", friendbotURLString, err)
		}
	}

	var rateLimit *throttled.RateQuota = nil
	perHourRateLimit := viper.GetInt("per-hour-rate-limit")
	if perHourRateLimit != 0 {
		rateLimit = &throttled.RateQuota{
			MaxRate:  throttled.PerHour(perHourRateLimit),
			MaxBurst: 100,
		}
	}

	config = horizon.Config{
		DatabaseURL:            viper.GetString("db-url"),
		StellarCoreDatabaseURL: viper.GetString("stellar-core-db-url"),
		StellarCoreURL:         viper.GetString("stellar-core-url"),
		Port:                   viper.GetInt("port"),
		HorizonDBMaxOpenConnections: viper.GetInt("horizon-db-max-open-connections"),
		HorizonDBMaxIdleConnections: viper.GetInt("horizon-db-max-idle-connections"),
		CoreDBMaxOpenConnections:    viper.GetInt("core-db-max-open-connections"),
		CoreDBMaxIdleConnections:    viper.GetInt("core-db-max-idle-connections"),
		ConnectionTimeout:           time.Duration(viper.GetInt("connection-timeout")) * time.Second,
		RateLimit:                   rateLimit,
		RateLimitRedisKey:           viper.GetString("rate-limit-redis-key"),
		RedisURL:                    viper.GetString("redis-url"),
		FriendbotURL:                friendbotURL,
		LogLevel:                    ll,
		LogFile:                     lf,
		MaxPathLength:               uint(viper.GetInt("max-path-length")),
		NetworkPassphrase:           viper.GetString("network-passphrase"),
		SentryDSN:                   viper.GetString("sentry-dsn"),
		LogglyToken:                 viper.GetString("loggly-token"),
		LogglyTag:                   viper.GetString("loggly-tag"),
		TLSCert:                     cert,
		TLSKey:                      key,
		Ingest:                      viper.GetBool("ingest"),
		HistoryRetentionCount:       uint(viper.GetInt("history-retention-count")),
		StaleThreshold:              uint(viper.GetInt("history-stale-threshold")),
		CursorName:                  viper.GetString("cursor-name"),
		SkipCursorUpdate:            viper.GetBool("skip-cursor-update"),
		EnableAssetStats:            viper.GetBool("enable-asset-stats"),
	}
=======
	// Configure log level
	log.DefaultLogger.Level = config.LogLevel
>>>>>>> horizon-v0.16.0
}
