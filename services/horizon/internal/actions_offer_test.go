package horizon

import (
	"testing"
	"time"
)

func TestOfferActions_Index(t *testing.T) {
	ht := StartHTTPTest(t, "trades")
	defer ht.Finish()

	w := ht.Get(
		"/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/offers",
	)

	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)

		//test last modified timestamp
		var records []map[string]interface{}
		ht.UnmarshalPage(w.Body, &records)
		t2018, err := time.Parse("2006-01-02", "2018-01-01")
		ht.Assert.NoError(err)
		recordTime, err := time.Parse("2006-01-02T15:04:05Z", records[2]["last_modified_time"].(string))
		ht.Assert.True(recordTime.After(t2018))
		ht.Assert.EqualValues(5, records[2]["last_modified_ledger"])
	}
}

func TestOfferActions_IndexNoLedgerData(t *testing.T) {
	ht := StartHTTPTest(t, "trades")
	defer ht.Finish()

	// Remove ledger data
	_, err := ht.App.HistoryQ().ExecRaw("DELETE FROM history_ledgers WHERE sequence=?", 5)
	ht.Assert.NoError(err)

	w := ht.Get(
		"/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/offers",
	)

	// Since 0.15.0 Horizon returns empty data instead of 500 error
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)

		//test last modified timestamp
		var records []map[string]interface{}
		ht.UnmarshalPage(w.Body, &records)
		ht.Assert.NotEmpty(records[2]["last_modified_ledger"])
		ht.Assert.Nil(records[2]["last_modified_time"])
	}
}
