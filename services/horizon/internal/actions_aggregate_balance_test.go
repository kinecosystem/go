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
		ht.Assert.Equal(999999996999999700, aggregateBalances.Embeded.Records[0].AggregateBalance)
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
		ht.Assert.Equal(999999996999999700, aggregateBalances.Embeded.Records[0].AggregateBalance)
		ht.Assert.Equal("GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", aggregateBalances.Embeded.Records[0].Id)
		ht.Assert.Equal(999999996999999700, aggregateBalances.Embeded.Records[1].AggregateBalance)
	}


	resp = ht.Get(
		"/accounts/SBOTI576E6NH4DWMRSWFW2SWMY7CCCGUX33DDJCG4RZLDSJ6RKZHI7NB/aggregate_balance",
	)
	if ht.Assert.Equal(200, resp.Code) {
		var aggregateBalances horizon.AggregatedBalances
		err := json.Unmarshal(resp.Body.Bytes(), &aggregateBalances)
		ht.Require.NoError(err)
		fmt.Printf("%+v\n", aggregateBalances)
		fmt.Printf("%+s\n", resp.Body)
		ht.Assert.Equal("GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", aggregateBalances.Embeded.Records[0].Id)
		ht.Assert.Equal(999999996999999700, aggregateBalances.Embeded.Records[0].AggregateBalance)
		ht.Assert.Equal("GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", aggregateBalances.Embeded.Records[0].Id)
		ht.Assert.Equal(999999996999999700, aggregateBalances.Embeded.Records[1].AggregateBalance)
	}
}
