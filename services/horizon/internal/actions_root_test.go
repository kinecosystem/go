package horizon

import (
	"encoding/json"
	"testing"

<<<<<<< HEAD
	"github.com/kinecosystem/go/services/horizon/internal/test"
	"github.com/kinecosystem/go/protocols/horizon"
=======
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/test"
>>>>>>> horizon-v0.15.4
)

func TestRootAction(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	server := test.NewStaticMockServer(`{
			"info": {
				"network": "test",
				"build": "test-core",
				"protocol_version": 4
			}
		}`)
	defer server.Close()

	ht.App.horizonVersion = "test-horizon"
	ht.App.config.StellarCoreURL = server.URL
	ht.App.config.NetworkPassphrase = "test"
	ht.App.UpdateStellarCoreInfo()

	w := ht.Get("/")
	if ht.Assert.Equal(200, w.Code) {
		var actual horizon.Root
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)
		ht.Assert.Equal("test-horizon", actual.HorizonVersion)
		ht.Assert.Equal("test-core", actual.StellarCoreVersion)
	}
}
