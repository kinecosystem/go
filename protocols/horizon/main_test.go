package horizon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Account Tests
// An example account to be used in all the Account tests
var exampleAccount = Account{
	Data: map[string]string{
		"test":    "aGVsbG8=",
		"invalid": "a_*&^*",
	},
}

// Testing the GetData method of Account
func TestAccount_GetData(t *testing.T) {
	// Should return the decoded value if the key exists
	decoded, err := exampleAccount.GetData("test")
	assert.Nil(t, err)
	assert.Equal(t, string(decoded), "hello")

	// Should return an empty slice if key doesn't exist
	decoded, err = exampleAccount.GetData("test2")
	assert.Nil(t, err)
	assert.Equal(t, len(decoded), 0)

	// Should return error slice if value is invalid
	_, err = exampleAccount.GetData("invalid")
	assert.NotNil(t, err)
}

func TestAccount_MustGetData(t *testing.T) {
	// Should return the decoded value if the key exists
	decoded := exampleAccount.MustGetData("test")
	assert.Equal(t, string(decoded), "hello")

	// Should return an empty slice if key doesn't exist
	decoded = exampleAccount.MustGetData("test2")
	assert.Equal(t, len(decoded), 0)

	// Should panic if the value is invalid
	assert.Panics(t, func() { exampleAccount.MustGetData("invalid") })
}

// Transaction Tests
// After marshalling and unmarshalling, the resulting struct should be the exact same as the original
func TestTransactionJSONMarshal(t *testing.T) {
	transaction := Transaction{
		ID:       "12345",
		FeePaid:  10,
		MemoType: "text",
		Memo:     "",
	}
	marshaledTransaction, marshalErr := json.Marshal(transaction)
	assert.Nil(t, marshalErr)
	var result Transaction
	json.Unmarshal(marshaledTransaction, &result)
	assert.Equal(t, result, transaction)
}

//For text memos, even if memo is an empty string, the resulting JSON should
// still include memo as a field
func TestTransactionEmptyMemoText(t *testing.T) {
	transaction := Transaction{
		MemoType: "text",
		Memo:     "",
	}
	marshaledTransaction, marshalErr := json.Marshal(transaction)
	assert.Nil(t, marshalErr)
	var result struct {
		Memo *string
	}
	json.Unmarshal(marshaledTransaction, &result)
	assert.NotNil(t, result.Memo)
}

// If a transaction's memo type is None, then the memo field should be omitted from JSON
func TestTransactionMemoTypeNone(t *testing.T) {
	transaction := Transaction{
		MemoType: "none",
	}
	marshaledTransaction, marshalErr := json.Marshal(transaction)
	assert.Nil(t, marshalErr)
	var result struct {
		Memo *string
	}
	json.Unmarshal(marshaledTransaction, &result)
	assert.Nil(t, result.Memo)
}
