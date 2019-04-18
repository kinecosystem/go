package horizon

import (
	"encoding/json"
	//"fmt"
	"github.com/kinecosystem/go/protocols/horizon"
	"testing"
)

func TestActionsControlledBalances_Show(t *testing.T) {
	ht := StartHTTPTest(t, "two_signatures")
	defer ht.Finish()

	// get controlled balance for scott - should not return any records
	w := ht.Get(
		"/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/controlled_balances",
	)
	if ht.Assert.Equal(200, w.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(w.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		ht.Assert.Equal(0, len(controlledBalances.Embeded.Records))
	}

	// get controlled balance for andrew - should return 2 accounts (scott and andrew)
	w1 := ht.Get(
		"/accounts/GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON/controlled_balances",
	)
	if ht.Assert.Equal(200, w1.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(w1.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		ht.Assert.Equal(2, len(controlledBalances.Embeded.Records))
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Id, "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON")
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Balance, 1050000000)
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Id, "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU")
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Balance, 949999800)
	}

	// missing account
	w2 := ht.Get("/accounts/GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ/controlled_balances")
	ht.Assert.Equal(404, w2.Code)
}
