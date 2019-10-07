package horizon

import (
	"encoding/json"
	"testing"

	"github.com/kinecosystem/go/protocols/horizon"
)

func TestAccountActions_Show(t *testing.T) {
	ht := StartHTTPTest(t, "allow_trust")
	defer ht.Finish()

	// existing account
	w := ht.Get(
		"/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
	)
	if ht.Assert.Equal(200, w.Code) {

		var result horizon.Account
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Require.NoError(err)
		ht.Assert.Equal("8589934593", result.Sequence)

		ht.Assert.NotEqual(0, result.LastModifiedLedger)
		for _, balance := range result.Balances {
			if balance.Type == "native" {
				ht.Assert.Equal(uint32(0), balance.LastModifiedLedger)
			} else {
				ht.Assert.NotEqual(uint32(0), balance.LastModifiedLedger)
			}
		}

	}

	// missing account
	w = ht.Get("/accounts/GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ")
	ht.Assert.Equal(404, w.Code)
}

func TestAccountActions_ShowRegressions(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	w := ht.Get(
		"/accounts/GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
	)
	var result horizon.Account
	err := json.Unmarshal(w.Body.Bytes(), &result)
	ht.Require.NoError(err)

	// Regression: no trades link
	ht.Assert.Contains(result.Links.Trades.Href, "/trades")

	// Regression: no data link
	ht.Assert.Contains(result.Links.Data.Href, "/data/{key}")
	ht.Assert.True(result.Links.Data.Templated)

	// Regression:  return 200 ok even when the history record cannot be found.

	// overwrite history with blank
	ht.T.ScenarioWithoutHorizon("base")
	w = ht.Get(
		"/accounts/GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
	)
	ht.Assert.Equal(200, w.Code)

}

func TestAccountActions_InvalidID(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// existing account
	w := ht.Get(
		"/accounts/=cr%FF%98%CB%F3%AF%E72%D85%FE%28%15y%8Fz%C4Ng%CE%98h%02%2A:%B6%FF%B9%CF%92%88O%91%10d&S%7C%9Bi%D4%CFI%28%CFo",
	)
	ht.Assert.Equal(400, w.Code)
}
