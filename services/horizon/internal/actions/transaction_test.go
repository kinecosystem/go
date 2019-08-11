package actions

import (
	"context"
	"testing"

	"github.com/stellar/go/services/horizon/internal/db2"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/test"
)

var defaultPage db2.PageQuery = db2.PageQuery{
	Order:  db2.OrderAscending,
	Limit:  db2.DefaultPageSize,
	Cursor: "",
}

func TestTransactionPage(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()

	ctx := context.Background()

	// filter by account
	page, err := TransactionPage(ctx, &history.Q{tt.HorizonSession()}, "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", 0, true, defaultPage, true)
	tt.Assert.NoError(err)
	tt.Assert.Equal(3, len(page.Embedded.Records))

	page, err = TransactionPage(ctx, &history.Q{tt.HorizonSession()}, "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", 0, true, defaultPage, true)
	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(page.Embedded.Records))

	page, err = TransactionPage(ctx, &history.Q{tt.HorizonSession()}, "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", 0, true, defaultPage, true)
	tt.Assert.NoError(err)
	tt.Assert.Equal(2, len(page.Embedded.Records))

	// filter by ledger
	page, err = TransactionPage(ctx, &history.Q{tt.HorizonSession()}, "", 1, true, defaultPage, true)
	tt.Assert.NoError(err)
	tt.Assert.Equal(0, len(page.Embedded.Records))

	page, err = TransactionPage(ctx, &history.Q{tt.HorizonSession()}, "", 2, true, defaultPage, true)
	tt.Assert.NoError(err)
	tt.Assert.Equal(3, len(page.Embedded.Records))

	page, err = TransactionPage(ctx, &history.Q{tt.HorizonSession()}, "", 3, true, defaultPage, true)
	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(page.Embedded.Records))

	// conflict fields
	_, err = TransactionPage(ctx, &history.Q{tt.HorizonSession()}, "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", 1, true, defaultPage, true)
	tt.Assert.Error(err)
}

func TestLoadTransactionRecords(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()

	// filter by account
	records, err := loadTransactionRecords(&history.Q{tt.HorizonSession()}, "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", 0, true, defaultPage)
	tt.Assert.NoError(err)
	tt.Assert.Equal(3, len(records))

	records, err = loadTransactionRecords(&history.Q{tt.HorizonSession()}, "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", 0, true, defaultPage)
	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(records))

	records, err = loadTransactionRecords(&history.Q{tt.HorizonSession()}, "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", 0, true, defaultPage)
	tt.Assert.NoError(err)
	tt.Assert.Equal(2, len(records))

	// filter by ledger
	records, err = loadTransactionRecords(&history.Q{tt.HorizonSession()}, "", 1, true, defaultPage)
	tt.Assert.NoError(err)
	tt.Assert.Equal(0, len(records))

	records, err = loadTransactionRecords(&history.Q{tt.HorizonSession()}, "", 2, true, defaultPage)
	tt.Assert.NoError(err)
	tt.Assert.Equal(3, len(records))

	records, err = loadTransactionRecords(&history.Q{tt.HorizonSession()}, "", 3, true, defaultPage)
	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(records))

	// conflict fields
	_, err = loadTransactionRecords(&history.Q{tt.HorizonSession()}, "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", 1, true, defaultPage)
	tt.Assert.Error(err)
}
