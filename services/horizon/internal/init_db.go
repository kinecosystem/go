package horizon

import (
<<<<<<< HEAD
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
	"github.com/kinecosystem/go/services/horizon/internal/db2/history"
	"github.com/kinecosystem/go/services/horizon/internal/log"
	"github.com/kinecosystem/go/support/db"
=======
	"github.com/stellar/go/services/horizon/internal/db2/core"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/support/log"
>>>>>>> horizon-v0.15.3
)

func initHorizonDb(app *App) {
	session, err := db.Open("postgres", app.config.DatabaseURL)

	if err != nil {
		log.Panic(err)
	}
<<<<<<< HEAD
	session.DB.SetMaxIdleConns(app.config.HorizonDBMaxIdleConnections)
	session.DB.SetMaxOpenConns(app.config.HorizonDBMaxOpenConnections)
=======

	// Make sure MaxIdleConns is equal MaxOpenConns. In case of high variance
	// in number of requests closing and opening connections may slow down Horizon.
	session.DB.SetMaxIdleConns(app.config.MaxDBConnections)
	session.DB.SetMaxOpenConns(app.config.MaxDBConnections)
>>>>>>> horizon-v0.15.3

	app.historyQ = &history.Q{session}
}

func initCoreDb(app *App) {
	session, err := db.Open("postgres", app.config.StellarCoreDatabaseURL)

	if err != nil {
		log.Panic(err)
	}

<<<<<<< HEAD
	session.DB.SetMaxIdleConns(app.config.CoreDBMaxIdleConnections)
	session.DB.SetMaxOpenConns(app.config.CoreDBMaxOpenConnections)
=======
	// Make sure MaxIdleConns is equal MaxOpenConns. In case of high variance
	// in number of requests closing and opening connections may slow down Horizon.
	session.DB.SetMaxIdleConns(app.config.MaxDBConnections)
	session.DB.SetMaxOpenConns(app.config.MaxDBConnections)
>>>>>>> horizon-v0.15.3
	app.coreQ = &core.Q{session}
}

func init() {
	appInit.Add("horizon-db", initHorizonDb, "app-context", "log")
	appInit.Add("core-db", initCoreDb, "app-context", "log")
}
