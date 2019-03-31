package core

import (
	"fmt"
		sq "github.com/Masterminds/squirrel"
		"github.com/kinecosystem/go/xdr"
		"github.com/kinecosystem/go/amount"
)

// AggregateBaalanceByAccountId returns the aggregate balance by the master account id
func (q *Q) AggregateBalanceByAccountId(dest interface{}, accountid string) (string, error) {
	sql := sq.Select("SUM(balance) aggbalance").From("accounts").Where(fmt.Sprintf("accountid = '%s' OR accountid IN (SELECT accountid FROM signers WHERE publickey = '%s')",accountid, accountid))
	result := struct {
	Balance int64  `db:"aggbalance"`
	}{}
	err := q.Get(&result, sql)
	return amount.String(xdr.Int64(result.Balance)), err
}