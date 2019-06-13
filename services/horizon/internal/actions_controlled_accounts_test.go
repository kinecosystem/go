package horizon

import (
	"encoding/json"
	"fmt"
	"github.com/kinecosystem/go/protocols/horizon"
	"testing"
)

func TestActionsControlledAccounts_Show(t *testing.T) {
	ht := StartHTTPTest(t, "two_signatures")
	defer ht.Finish()

	// get controlled balance for scott - should not return any records
	resp := ht.Get(
		"/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		ht.Assert.Equal(0, len(ControlledAccounts.Embeded.Records))
	}

	// get controlled balance for kp1 - should not return anything as kp1 isnt a signer anywhere else
	resp1 := ht.Get(
		"/accounts/GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp1.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp1.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		fmt.Print(resp1.Body)
		ht.Assert.Equal(0, len(ControlledAccounts.Embeded.Records))
	}

	// get controlled balance for kp2 - should return both kp1 and kp2
	resp2 := ht.Get(
		"/accounts/GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp2.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp2.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		fmt.Print(resp2.Body)
		ht.Assert.Equal(2, len(ControlledAccounts.Embeded.Records))
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Id, "GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7") // kp1
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Balance, "9999.99900")
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Id, "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y") // kp2
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Balance, "19999.99900")
	}

	// get controlled balance for kp3 - should return both kp2 and kp3
	resp3 := ht.Get(
		"/accounts/GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp3.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp3.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		fmt.Print(resp3.Body)
		ht.Assert.Equal(2, len(ControlledAccounts.Embeded.Records))
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Id, "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y") // kp2
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Balance, "19999.99900")
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Id, "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ") // kp3
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Balance, "29999.99900")
	}

	// get controlled balance for kp4 - should return both kp3 and kp4
	resp4 := ht.Get(
		"/accounts/GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp4.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp4.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		fmt.Print(resp4.Body)
		ht.Assert.Equal(2, len(ControlledAccounts.Embeded.Records))
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Id, "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ") // kp3
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Balance, "29999.99900")
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Id, "GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6") // kp4
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Balance, "40000.00000") //didnt pay any fees
	}


	// try a non exsiting account
	resp5 := ht.Get("/accounts/GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ/controlled_accounts")
	ht.Assert.Equal(404, resp5.Code)


	// get kp7's controleld accounts. kp7 is a signer in kp5, and kp5 is a signer in kp6, but that should make no difference.
	// we should only see kp7, kp5 accounts
	resp6 := ht.Get(
		"/accounts/GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp6.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp6.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		fmt.Print(resp6.Body)
		ht.Assert.Equal(2, len(ControlledAccounts.Embeded.Records))
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Id, "GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ") // kp5
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Balance, "49999.99900")
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Id, "GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL") // kp7
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Balance, "69999.99900") 
	}

	// get kp8's controleld accounts. kp8 is a signer in kp9, and kp9 is a signer in kp8, but that should make no difference.
	resp7 := ht.Get(
		"/accounts/GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp7.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp7.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		fmt.Print(resp7.Body)
		ht.Assert.Equal(2, len(ControlledAccounts.Embeded.Records))
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Id, "GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN") // kp8
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Balance, "79999.99900")
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Id, "GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT") // kp9
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Balance, "89999.99900") 
	}

	// get kp9's controleld accounts. kp9 is a signer in kp8, and kp8 is a signer in kp9, but that should make no difference.
	resp8 := ht.Get(
		"/accounts/GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT/controlled_accounts",
	)
	if ht.Assert.Equal(200, resp8.Code) {
		var ControlledAccounts horizon.ControlledAccounts
		err := json.Unmarshal(resp8.Body.Bytes(), &ControlledAccounts)
		ht.Require.NoError(err)
		fmt.Print(resp8.Body)
		ht.Assert.Equal(2, len(ControlledAccounts.Embeded.Records))
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Id, "GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN") // kp8
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[0].Balance, "79999.99900")
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Id, "GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT") // kp9
		ht.Assert.Equal(ControlledAccounts.Embeded.Records[1].Balance, "89999.99900") 
	}
}