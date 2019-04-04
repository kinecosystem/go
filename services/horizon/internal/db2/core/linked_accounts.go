package core

import (
	sqltypes "database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/kinecosystem/go/xdr"
)

// AggregateBaalanceByAccountId returns the aggregate balance by the master account id
func (q *Q) AggregateBalanceByAccountId(dest *ControlledBalance) error {
	sql := sq.Select("SUM(balance) aggbalance").From("accounts").Where(fmt.Sprintf("accountid = '%s' OR accountid IN (SELECT accountid FROM signers WHERE publickey = '%s')", dest.AccountId, dest.AccountId))
	result := struct {
		NullableBalance sqltypes.NullInt64 `db:"aggbalance"`
	}{}
	err := q.Get(&result, sql)
	if err != nil {
		return err
	}
	if result.NullableBalance.Valid {
		dest.Balance = xdr.Int64(result.NullableBalance.Int64)
	} else {
		err = sqltypes.ErrNoRows // translates to 404 not found
	}
	return err
}

// ControlledBalancesByAccountId returns the list of controlled accounts with their respective balance
func (q *Q) ControlledBalancesByAccountId(dest *[]*ControlledBalance, aid string) error {
	sql := sq.Select("a.accountid id, min(a.balance) balance").
		From("accounts a, signers s").
		Where(fmt.Sprintf("s.publickey = '%s' AND (a.accountid = s.accountid OR a.accountid = s.publickey)", aid)).
		GroupBy("a.accountid")
	rows, err := q.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	var balance_int int64
	for rows.Next() {
		cb := ControlledBalance{}
		err = rows.Scan(&cb.AccountId, &balance_int)
		if err != nil {
			return err
		}
		// convert the int64 to xdr
		cb.Balance = xdr.Int64(balance_int)
		// append the results into the slice:
		*dest = append(*dest, &cb)
	}
	return err
}
