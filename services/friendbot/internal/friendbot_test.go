package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sync"
)

// REGRESSION:  ensure that we can craft a transaction
func TestFriendbot_makeTx(t *testing.T) {
	fb := &Bot{
		Secret:          "SAQWC7EPIYF3XGILYVJM4LVAVSLZKT27CTEI3AFBHU2VRCMQ3P3INPG5",
		Network:         "Test SDF Network ; September 2015",
		StartingBalance: "10000",
		sequence:        2,
	}

	txn, err := fb.makeTx("GDJIN6W6PLTPKLLM57UW65ZH4BITUXUMYQHIMAZFYXF45PZVAWDBI77Z", fb.StartingBalance)
	if !assert.NoError(t, err) {
		return
	}
	expectedTxn := "AAAAAPuYf7x7KGvFX9fjCR9WIaoTX3yHJYwX6ZSx6w76HPjEAAAAZAAAAAAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA" +
		"0ob63nrm9S1s7+lvdyfgUTpejMQOhgMlxcvOvzUFhhQAAAAAAJiWgAAAAAAAAAAB+hz4xAAAAEAxJX1arTWrUylEZNhvMezff4pLc5CcgIfQcF" +
		"//EXoudKwXSgQGFNivsgVRKL7djlsUkr1m8nMIgZ7lqibe46EB"
	assert.Equal(t, expectedTxn, txn)

	// ensure we're race free. NOTE:  presently, gb can't
	// run with -race on... we'll confirm this works when
	// horizon is in the monorepo
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		_, err := fb.makeTx("GDJIN6W6PLTPKLLM57UW65ZH4BITUXUMYQHIMAZFYXF45PZVAWDBI77Z", fb.StartingBalance)
		// don't assert on the txn value here because the ordering is not guaranteed between these 2 goroutines
		assert.NoError(t, err)
		wg.Done()
	}()
	go func() {
		_, err := fb.makeTx("GDJIN6W6PLTPKLLM57UW65ZH4BITUXUMYQHIMAZFYXF45PZVAWDBI77Z", fb.StartingBalance)
		assert.NoError(t, err)
		wg.Done()
	}()
	wg.Wait()
}
