package horizon

import (
	"encoding/json"
	"testing"
)

func TestOperationFeeTestsActions_Show(t *testing.T) {

	testCases := []struct {
		scenario    string
		min         string
		mode        string
		lastbasefee string
	}{
		// happy path
		{
			"operation_fee_stats_1",
			"100",
			"100",
			"100",
		},
		// no transactions in last 5 ledgers
		{
			"operation_fee_stats_2",
			"100",
			"100",
			"100",
		},
		// transactions with varying fees
		{
			"operation_fee_stats_3",
			"200",
			"400",
			"100",
		},
	}

	for _, kase := range testCases {
		t.Run("/operation_fee_stats", func(t *testing.T) {
			ht := StartHTTPTest(t, kase.scenario)
			defer ht.Finish()

			w := ht.Get("/operation_fee_stats")

			if ht.Assert.Equal(200, w.Code) {
				var result map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &result)
				ht.Require.NoError(err)
				ht.Assert.Equal(kase.min, result["min_accepted_fee"])
				ht.Assert.Equal(kase.mode, result["mode_accepted_fee"])
				ht.Assert.Equal(kase.lastbasefee, result["last_ledger_base_fee"])
			}
		})
	}
}
