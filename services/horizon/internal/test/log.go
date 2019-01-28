package test

import (
	"github.com/sirupsen/logrus"
<<<<<<< HEAD
	"github.com/kinecosystem/go/services/horizon/internal/log"
=======
	"github.com/stellar/go/support/log"
>>>>>>> horizon-v0.15.3
)

var testLogger *log.Entry

func init() {
	testLogger = log.New()
	testLogger.Entry.Logger.Formatter.(*logrus.TextFormatter).DisableColors = true
	testLogger.Entry.Logger.Level = logrus.DebugLevel
}
