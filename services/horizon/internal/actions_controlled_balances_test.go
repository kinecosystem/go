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
	resp := ht.Get(
		"/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/controlled_balances",
	)
	if ht.Assert.Equal(200, resp.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(resp.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		ht.Assert.Equal(0, len(controlledBalances.Embeded.Records))
	}

	// get controlled balance for kp1 - should return both kp1 and kp2
	resp2 := ht.Get(
		"/accounts/SBOTI576E6NH4DWMRSWFW2SWMY7CCCGUX33DDJCG4RZLDSJ6RKZHI7NB/controlled_balances",
	)
	if ht.Assert.Equal(200, resp2.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(resp2.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		ht.Assert.Equal(2, len(controlledBalances.Embeded.Records))
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Id, "GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7")
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Balance, 1050000000)
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Id, "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y")
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Balance, 949999800)
	}

	// missing account
	resp3 := ht.Get("/accounts/GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ/controlled_balances")
	ht.Assert.Equal(404, resp3.Code)
}
