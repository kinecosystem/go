package io

import (
	"io"
	"sync"

	"github.com/stellar/go/exp/ingest/ledgerbackend"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

// LedgerTransaction represents the data for a single transaction within a ledger.
type LedgerTransaction struct {
	Index      uint32
	Envelope   xdr.TransactionEnvelope
	Result     xdr.TransactionResultPair
	Meta       xdr.TransactionMeta
	FeeChanges xdr.LedgerEntryChanges
}

// DBLedgerReadCloser is a database-backed implementation of the io.LedgerReadCloser interface.
type DBLedgerReadCloser struct {
	sequence     uint32
	backend      ledgerbackend.LedgerBackend
	header       xdr.LedgerHeaderHistoryEntry
	transactions []LedgerTransaction
	readIdx      int
	initOnce     sync.Once
	readMutex    sync.Mutex
}

// Ensure DBLedgerReadCloser implements LedgerReadCloser
var _ LedgerReadCloser = (*DBLedgerReadCloser)(nil)

// MakeLedgerReadCloser is a factory method for LedgerReadCloser.
func MakeLedgerReadCloser(sequence uint32, backend ledgerbackend.LedgerBackend) *DBLedgerReadCloser {
	return &DBLedgerReadCloser{
		sequence: sequence,
		backend:  backend,
	}
}

// GetSequence returns the sequence number of the ledger data stored by this object.
func (dblrc *DBLedgerReadCloser) GetSequence() uint32 {
	return dblrc.sequence
}

// GetHeader returns the XDR Header data associated with the stored ledger.
func (dblrc *DBLedgerReadCloser) GetHeader() (xdr.LedgerHeaderHistoryEntry, error) {
	var err error
	dblrc.initOnce.Do(func() { err = dblrc.init() })
	if err != nil {
		return xdr.LedgerHeaderHistoryEntry{}, err
	}
	return dblrc.header, nil
}

// Read returns the next transaction in the ledger, ordered by tx number, each time it is called. When there
// are no more transactions to return, an EOF error is returned.
func (dblrc *DBLedgerReadCloser) Read() (LedgerTransaction, error) {
	var err error
	dblrc.initOnce.Do(func() { err = dblrc.init() })
	if err != nil {
		return LedgerTransaction{}, err
	}

	// Protect all accesses to dblrc.readIdx
	dblrc.readMutex.Lock()
	defer dblrc.readMutex.Unlock()

	if dblrc.readIdx < len(dblrc.transactions) {
		dblrc.readIdx++
		return dblrc.transactions[dblrc.readIdx-1], nil
	}
	return LedgerTransaction{}, io.EOF
}

// Close moves the read pointer so that subsequent calls to Read() will return EOF.
func (dblrc *DBLedgerReadCloser) Close() error {
	dblrc.readMutex.Lock()
	dblrc.readIdx = len(dblrc.transactions)
	dblrc.readMutex.Unlock()

	return nil
}

// Init pulls data from the backend to set this object up for use.
func (dblrc *DBLedgerReadCloser) init() error {
	exists, ledgerCloseMeta, err := dblrc.backend.GetLedger(dblrc.sequence)

	if err != nil {
		return errors.Wrap(err, "error reading ledger from backend")
	}
	if !exists {
		return errors.Wrap(err, "ledger was not found")
	}

	dblrc.header = ledgerCloseMeta.LedgerHeader

	dblrc.storeTransactions(ledgerCloseMeta)

	return nil
}

// storeTransactions maps the close meta data into a slice of LedgerTransaction structs, to provide
// a per-transaction view of the data when Read() is called.
func (dblrc *DBLedgerReadCloser) storeTransactions(lcm ledgerbackend.LedgerCloseMeta) {
	for i := range lcm.TransactionEnvelope {
		dblrc.transactions = append(dblrc.transactions, LedgerTransaction{
			Index:      uint32(i + 1), // Transactions start at '1'
			Envelope:   lcm.TransactionEnvelope[i],
			Result:     lcm.TransactionResult[i],
			Meta:       lcm.TransactionMeta[i],
			FeeChanges: lcm.TransactionFeeChanges[i],
		})
	}
}
