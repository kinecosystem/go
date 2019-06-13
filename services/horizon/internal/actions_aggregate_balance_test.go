package horizon

import (
	"encoding/json"
	"fmt"
	"github.com/kinecosystem/go/protocols/horizon"
	"testing"
)

func TestActionsAggregateBalance_Show(t *testing.T) {
	//test one plain account with no slaves
	//and one account that doesn't exist

	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// existing account
	resp := ht.Get(
		"/accounts/GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H/aggregate_balance",
	)
	if ht.Assert.Equal(200, resp.Code) {
		var aggregateBalances horizon.AggregatedBalances
		err := json.Unmarshal(resp.Body.Bytes(), &aggregateBalances)
		ht.Require.NoError(err)
		fmt.Printf("%+v\n", aggregateBalances)
		fmt.Printf("%+s\n", resp.Body)
		ht.Assert.Equal("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", aggregateBalances.Embeded.Records[0].Id)
		ht.Assert.Equal("9999999969999.99700", aggregateBalances.Embeded.Records[0].AggregateBalance)
	}

	// missing account
	resp = ht.Get("/accounts/GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ/aggregate_balance")
	ht.Assert.Equal(404, resp.Code)
}

func TestActionsAggregateBalance_Show2(t *testing.T) {
	ht := StartHTTPTest(t, "two_signatures")
	defer ht.Finish()

	// andrew's account - should have a total of 2 acconts - andrew and scott
	resp := ht.Get(
		"/accounts/GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON/aggregate_balance",
	)
	if ht.Assert.Equal(200, resp.Code) {
		var aggregateBalances horizon.AggregatedBalances
		err := json.Unmarshal(resp.Body.Bytes(), &aggregateBalances)
		ht.Require.NoError(err)
		fmt.Printf("%+v\n", aggregateBalances)
		fmt.Printf("%+s\n", resp.Body)
		ht.Assert.Equal("GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", aggregateBalances.Embeded.Records[0].Id)
		ht.Assert.Equal("19999.99800", aggregateBalances.Embeded.Records[0].AggregateBalance) // andrew + scott
	}

	// get kp1's aggregate balance
	resp2 := ht.Get(
		"/accounts/GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7/aggregate_balance",
	)
	if ht.Assert.Equal(200, resp2.Code) {
		var aggregateBalances horizon.AggregatedBalances
		err := json.Unmarshal(resp2.Body.Bytes(), &aggregateBalances)
		ht.Require.NoError(err)
		fmt.Printf("%+v\n", aggregateBalances)
		fmt.Printf("%+s\n", resp2.Body)
		ht.Assert.Equal("GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7", aggregateBalances.Embeded.Records[0].Id)
		ht.Assert.Equal("9999.99900", aggregateBalances.Embeded.Records[0].AggregateBalance) // kp1 only has itself, and payed for kp2
	}

	// get kp2's aggregate balance
	resp3 := ht.Get(
		"/accounts/GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y/aggregate_balance",
	)
	if ht.Assert.Equal(200, resp3.Code) {
		var aggregateBalances horizon.AggregatedBalances
		err := json.Unmarshal(resp3.Body.Bytes(), &aggregateBalances)
		ht.Require.NoError(err)
		fmt.Printf("%+v\n", aggregateBalances)
		fmt.Printf("%+s\n", resp3.Body)
		ht.Assert.Equal("GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y", aggregateBalances.Embeded.Records[0].Id)
		ht.Assert.Equal("29999.99800", aggregateBalances.Embeded.Records[0].AggregateBalance) // kp2 (payed for kp3) + kp1 (payed for kp2) =~3000
	}
}
