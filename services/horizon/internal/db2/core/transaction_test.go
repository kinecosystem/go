package core

import (
	"testing"

	"github.com/kinecosystem/go/services/horizon/internal/test"
	"github.com/kinecosystem/go/xdr"
)

func TestTransactionsQueries(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()
	q := &Q{tt.CoreSession()}

	// Test TransactionsByLedger
	var txs []Transaction
	err := q.TransactionsByLedger(&txs, 2)

	if tt.Assert.NoError(err) {
		tt.Assert.Len(txs, 3)
	}

	// Test TransactionByHashAfterLedger
	var tx Transaction
	err = q.TransactionByHashAfterLedger(&tx, "cebb875a00ff6e1383aef0fd251a76f22c1f9ab2a2dffcb077855736ade2659a", 2)

	if tt.Assert.NoError(err) {
		tt.Assert.Equal(int32(3), tx.LedgerSequence)
	}

	err = q.TransactionByHashAfterLedger(&tx, "cebb875a00ff6e1383aef0fd251a76f22c1f9ab2a2dffcb077855736ade2659a", 3)

	if tt.Assert.Error(err) {
		tt.Assert.True(q.NoRows(err))
	}
}

func TestMemo(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()

	var tx Transaction

	xdr.SafeUnmarshalBase64("AAAAAMvoFDdcyQrJAcBmRdyEnW6047pvlk4MS/4r0n/1WH8VAAAAZAACnMAAAAACAAAAAAAAAAEAAAARADEuMC4xb3dlcnJpZGUgbWUAAAAAAAABAAAAAQAAAACJzogbLxrrmN7N5JVQceSxl8jkED26RGzbyyRIpwTh6wAAAAoAAAAWaSBzaG91bGQgYmUgb3dlcnJpZGRlbgAAAAAAAQAAABVpIHNob3VsZCBiZSBvd2VycmlkZW4AAAAAAAAAAAAAAacE4esAAABA0GuCIEmKyQ2DRqt5+BOIqjVlHisjY6rK1IcOtzjIKCDgSAoiv5yhYe09PohBH91TXvAQ/LZJj5hVMihfMjtgCw==", &tx.Envelope)

	tt.Assert.Equal("1.0.1owerride me", tx.Memo().String)
}

func TestSignatures(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()

	var tx Transaction

	// https://github.com/stellar/stellar-core/issues/1225
	xdr.SafeUnmarshalBase64("AAAAAMIK9djC7k75ziKOLJcvMAIBG7tnBuoeI34x+Pi6zqcZAAAAZAAZphYAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAynnCTTyw53VVRLOWX6XKTva63IM1LslPNW01YB0hz/8AAAAAAAAAAlQL5AAAAAAAAAAAAh0hz/8AAABA8qkkeKaKfsbgInyIkzXJhqJE5/Ufxri2LdxmyKkgkT6I3sPmvrs5cPWQSzEQyhV750IW2ds97xTHqTpOfuZCAnhSuFUAAAAA", &tx.Envelope)

	signatures := tx.Base64Signatures()

	tt.Assert.Equal(2, len(signatures))
	tt.Assert.Equal("8qkkeKaKfsbgInyIkzXJhqJE5/Ufxri2LdxmyKkgkT6I3sPmvrs5cPWQSzEQyhV750IW2ds97xTHqTpOfuZCAg==", signatures[0])
	tt.Assert.Equal("", signatures[1])
}

func TestFee(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()
	var tx Transaction
	xdr.SafeUnmarshalBase64("AAAAAEfdynqlBdqSRiFQ+TcvlVB0Vr025FEJhm2H8k2Drs7QAAAD5wAbppsAAAAPAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAABAAAAAEfdynqlBdqSRiFQ+TcvlVB0Vr025FEJhm2H8k2Drs7QAAAAAAAAAAAD+DxAAAAAAAAAAAGDrs7QAAAAQGYgWJlJrfeRHuGT2KZqePqlbLKUeD4KWStP6QVyWdC84trpcHH84zcHcdz+j5tWtECyITn9pvxfHXbch5x0bgs=", &tx.Envelope)

	// #1
	// Source account is whitelisted
	WhiteListData := map[string]string{
		"GBD53ST2UUC5VESGEFIPSNZPSVIHIVV5G3SFCCMGNWD7ETMDV3HNA2JV": "dSj0bg==",
	}
	tt.Assert.Equal(tx.Fee(WhiteListData), 0)
	tt.Assert.NotEqual(tx.Fee(WhiteListData), 999)

	// #2
	// Source account is NOT whitelisted
	WhiteListData = map[string]string{
		"GAMRH6ZXD2ZMXUOPUBDFHPJXWGIMVTIU26RTWV4OLJ5SQCLVFD2G552E": "dSj0bg==",
	}
	tt.Assert.Equal(tx.Fee(WhiteListData), 999)
	tt.Assert.NotEqual(tx.Fee(WhiteListData), 0)

	xdr.SafeUnmarshalBase64("AAAAAGw/UOu5NoueYpdBpxFiiWlMoooS57T7/hwA/6ISEWYAAAAD5wAeR6IAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAABAAAAAGw/UOu5NoueYpdBpxFiiWlMoooS57T7/hwA/6ISEWYAAAAAAAAAAAAD+DxAAAAAAAAAAAISEWYAAAAAQBfY4JbVVG0uCJVbwiWpOWynWMHtlNYXgEq9AwTycKh2IZyOFBP3UkdzqqqvgO4q73/9OTigmU3X8ODZ8TkmMguDrs7QAAAAQK49fDeOKvOLEKEHRE0Bhdg6sV3jGxuHRZwiutnNxJNAQonk2M79ZVyc/qqASwJawRVv/+KOVEuSDFYw5W/RWAc=", &tx.Envelope)
	// #3
	// Source account is NOT whitelisted BUT whitelisted account sign on the transaction
	WhiteListData = map[string]string{
		"GBD53ST2UUC5VESGEFIPSNZPSVIHIVV5G3SFCCMGNWD7ETMDV3HNA2JV": "dSj0bg==",
	}
	tt.Assert.Equal(tx.Fee(WhiteListData), 0)
	tt.Assert.NotEqual(tx.Fee(WhiteListData), 999)

}
