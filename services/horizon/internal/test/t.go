package test

import (
	"io"

	"encoding/json"

	"github.com/guregu/null"
	"github.com/stellar/go/services/horizon/internal/ledger"
	"github.com/stellar/go/services/horizon/internal/operationfeestats"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/support/render/hal"
)

// CoreSession returns a db.Session instance pointing at the stellar core test database
func (t *T) CoreSession() *db.Session {
	return &db.Session{
		DB:  t.CoreDB,
		Ctx: t.Ctx,
	}
}

// Finish finishes the test, logging any accumulated horizon logs to the logs
// output
func (t *T) Finish() {
	RestoreLogger()
	// Reset cached ledger state
	ledger.SetState(ledger.State{})
	operationfeestats.ResetState()

	if t.LogBuffer.Len() > 0 {
		t.T.Log("\n" + t.LogBuffer.String())
	}
}

// HorizonSession returns a db.Session instance pointing at the horizon test
// database
func (t *T) HorizonSession() *db.Session {
	return &db.Session{
		DB:  t.HorizonDB,
		Ctx: t.Ctx,
	}
}

// Scenario loads the named sql scenario into the database
func (t *T) Scenario(name string) *T {
	LoadScenario(name)
	t.UpdateLedgerState()
	t.UpdateOperationFeeStatsState()
	return t
}

// ScenarioWithoutHorizon loads the named sql scenario into the database
func (t *T) ScenarioWithoutHorizon(name string) *T {
	LoadScenarioWithoutHorizon(name)
	t.UpdateLedgerState()
	t.UpdateOperationFeeStatsState()
	return t
}

// UnmarshalPage populates dest with the records contained in the json-encoded page in r
func (t *T) UnmarshalPage(r io.Reader, dest interface{}) hal.Links {
	var env struct {
		Embedded struct {
			Records json.RawMessage `json:"records"`
		} `json:"_embedded"`
		Links struct {
			Self hal.Link `json:"self"`
			Next hal.Link `json:"next"`
			Prev hal.Link `json:"prev"`
		} `json:"_links"`
	}

	err := json.NewDecoder(r).Decode(&env)
	t.Require.NoError(err, "failed to decode page")

	err = json.Unmarshal(env.Embedded.Records, dest)
	t.Require.NoError(err, "failed to decode records")

	return env.Links
}

// UnmarshalNext extracts and returns the next link
func (t *T) UnmarshalNext(r io.Reader) string {
	var env struct {
		Links struct {
			Next struct {
				Href string `json:"href"`
			} `json:"next"`
		} `json:"_links"`
	}

	err := json.NewDecoder(r).Decode(&env)
	t.Require.NoError(err, "failed to decode page")
	return env.Links.Next.Href
}

// UnmarshalExtras extracts and returns extras content
func (t *T) UnmarshalExtras(r io.Reader) map[string]string {
	var resp struct {
		Extras map[string]string `json:"extras"`
	}

	err := json.NewDecoder(r).Decode(&resp)
	t.Require.NoError(err, "failed to decode page")

	return resp.Extras
}

// UpdateLedgerState updates the cached ledger state (or panicing on failure).
func (t *T) UpdateLedgerState() {
	var next ledger.State

	err := t.CoreSession().GetRaw(&next, `
		SELECT
			COALESCE(MAX(ledgerseq), 0) as core_latest
		FROM ledgerheaders
	`)

	if err != nil {
		panic(err)
	}

	err = t.HorizonSession().GetRaw(&next, `
			SELECT
				COALESCE(MIN(sequence), 0) as history_elder,
				COALESCE(MAX(sequence), 0) as history_latest
			FROM history_ledgers
		`)

	if err != nil {
		panic(err)
	}

	ledger.SetState(next)
	return
}

// UpdateOperationFeeStatsState updates the cached operation fees state
// or panics on failure.
func (t *T) UpdateOperationFeeStatsState() {
	var err error
	var next operationfeestats.State

	var latest struct {
		BaseFee  int32 `db:"base_fee"`
		Sequence int32 `db:"sequence"`
	}
	var feeStats struct {
		Min  null.Int `db:"min"`
		Mode null.Int `db:"mode"`
	}

	cur := operationfeestats.CurrentState()

	err = t.HorizonSession().GetRaw(&latest, `
		SELECT base_fee, sequence
		FROM history_ledgers
		WHERE sequence = (SELECT COALESCE(MAX(sequence), 0) FROM history_ledgers)
	`)
	if err != nil {
		return
	}

	next.LastBaseFee = int64(latest.BaseFee)
	next.LastLedger = int64(latest.Sequence)

	// finish early if no new ledgers
	if cur.LastLedger == int64(latest.Sequence) {
		return
	}

	err = t.HorizonSession().GetRaw(&feeStats, `
		SELECT min(fee_paid/operation_count),  mode() within group (order by fee_paid/operation_count)
		FROM history_transactions
		WHERE ledger_sequence > $1 AND ledger_sequence <= $2
	`, latest.Sequence-5, latest.Sequence)
	if err != nil {
		return
	}

	// if no transactions in last X ledgers, return
	// latest ledger's base fee for all
	if !feeStats.Mode.Valid && !feeStats.Min.Valid {
		next.Min = next.LastBaseFee
		next.Mode = next.LastBaseFee
	} else {
		next.Min = feeStats.Min.Int64
		next.Mode = feeStats.Mode.Int64
	}

	operationfeestats.SetState(next)
	return
}
