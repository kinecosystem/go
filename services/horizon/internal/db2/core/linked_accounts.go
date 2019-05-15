package core

import (
	sqltypes "database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/kinecosystem/go/xdr"
)

// AggregateBaalanceByAccountId returns the aggregate balance by the master account id
func (q *Q) AggregateBalanceByAccountId(dest *ControlledBalance) error {
	sql := sq.Select("SUM(balance) aggbalance").
		From("accounts").
		Where(fmt.Sprintf("accountid = '%s' OR accountid IN (SELECT accountid FROM signers WHERE publickey = '%s')",
			dest.AccountId, dest.AccountId))
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
		// translates to 404 not found
		err = sqltypes.ErrNoRows
	}
	return err
}

// ControlledBalancesByAccountId returns the list of controlled accounts with their respective balance
func (q *Q) ControlledBalancesByAccountId(controlledBalanceSlice []*ControlledBalance, aid string) error {
	// Ensure that the account exists
	IsAccountExistsQuery := sq.Select("count (*)").
		From("accounts a").
		Where(fmt.Sprintf("a.accountid='%s'", aid))
	var count int
	err := q.Get(&count, IsAccountExistsQuery)
	if err != nil {
		return err
	}
	if count == 0 {
		// translates to 404 not found
		return sqltypes.ErrNoRows
	}

	// Get the data for the account
	GetAccountsQuery := sq.Select("a.accountid id, min(a.balance) balance").
		From("accounts a, signers s").
		Where(fmt.Sprintf("s.publickey = '%s' AND (a.accountid = s.accountid OR a.accountid = s.publickey)", aid)).
		GroupBy("a.accountid")
	rows, err := q.Query(GetAccountsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var balanceInt int64
	for rows.Next() {
		var cb ControlledBalance
		err = rows.Scan(&cb.AccountId, &balanceInt)
		if err != nil {
			return err
		}
		// convert the int64 to xdr
		cb.Balance = xdr.Int64(balanceInt)
		// append the results into the slice:
		controlledBalanceSlice = append(controlledBalanceSlice, &cb)
	}
	return err
}
