package horizon

import (
	"encoding/json"
	"fmt"
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

	// get controlled balance for kp1 - should not return anything as kp1 isnt a signer anywhere else
	resp1 := ht.Get(
		"/accounts/GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7/controlled_balances",
	)
	if ht.Assert.Equal(200, resp1.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(resp1.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		fmt.Print(resp1.Body)
		ht.Assert.Equal(0, len(controlledBalances.Embeded.Records))
	}

	// get controlled balance for kp2 - should return both kp1 and kp2
	resp2 := ht.Get(
		"/accounts/GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y/controlled_balances",
	)
	if ht.Assert.Equal(200, resp2.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(resp2.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		fmt.Print(resp2.Body)
		ht.Assert.Equal(2, len(controlledBalances.Embeded.Records))
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Id, "GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7") // kp1
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Balance, "9999.99900")
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Id, "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y") // kp2
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Balance, "19999.99900")
	}

	// get controlled balance for kp3 - should return both kp2 and kp3
	resp3 := ht.Get(
		"/accounts/GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ/controlled_balances",
	)
	if ht.Assert.Equal(200, resp3.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(resp3.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		fmt.Print(resp3.Body)
		ht.Assert.Equal(2, len(controlledBalances.Embeded.Records))
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Id, "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y") // kp2
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Balance, "19999.99900")
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Id, "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ") // kp3
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Balance, "29999.99900")
	}

	// get controlled balance for kp4 - should return both kp3 and kp4
	resp4 := ht.Get(
		"/accounts/GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6/controlled_balances",
	)
	if ht.Assert.Equal(200, resp4.Code) {
		var controlledBalances horizon.ControlledBalances
		err := json.Unmarshal(resp4.Body.Bytes(), &controlledBalances)
		ht.Require.NoError(err)
		fmt.Print(resp4.Body)
		ht.Assert.Equal(2, len(controlledBalances.Embeded.Records))
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Id, "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ") // kp3
		ht.Assert.Equal(controlledBalances.Embeded.Records[1].Balance, "29999.99900")
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Id, "GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6") // kp4
		ht.Assert.Equal(controlledBalances.Embeded.Records[0].Balance, "40000.00000") //didnt pay any fees
	}


	// try a non exsiting account
	resp5 := ht.Get("/accounts/GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ/controlled_balances")
	ht.Assert.Equal(404, resp5.Code)
}


// kp1
// GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7
// SAYYMGTGIWMR6VO243LU5GG2J5YYMAPFK5ZSTITHPUV5DVUS3K6Q4ZNF
// should have ~1000 kins

// kp2
// GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y
// SBMTHT5NRNVJFBYDON5COPZRCRCMDOQNRPAOMVJJV6W7DFW32EKRN6RV
// should have ~2000 kins

// kp3
// GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ
// SBOTI576E6NH4DWMRSWFW2SWMY7CCCGUX33DDJCG4RZLDSJ6RKZHI7NB
// should have ~3000 kins

// kp4
// GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6
// SC3CILTVMLKPF7YJXXIESDXNIA2QX6XX47WMLAY3UI5YQL75Y2SRXOZN
// should have ~4000 kins