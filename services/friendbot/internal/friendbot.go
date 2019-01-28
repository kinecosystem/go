package internal

import (
	"strconv"
	"sync"

	b "github.com/kinecosystem/go/build"
	"github.com/kinecosystem/go/clients/horizon"
	"github.com/kinecosystem/go/keypair"
	"github.com/kinecosystem/go/support/errors"
)

// TxResult is the result from the asynchronous submit transaction method over a channel
type TxResult struct {
	maybeTransactionSuccess *horizon.TransactionSuccess
	maybeErr                error
}

// Bot represents the friendbot subsystem.
type Bot struct {
	Horizon           *horizon.Client
	Secret            string
	Network           string
	StartingBalance   string
	SubmitTransaction func(bot *Bot, channel chan TxResult, signed string)

	// uninitialized
	sequence             uint64
	forceRefreshSequence bool
	lock                 sync.Mutex
}

<<<<<<< HEAD
// Pay funds to account.
//
// amount is the amount to fund, capped at Bot.StartingBalance.
// Setting isFundRequest to true will submit a transaction with a Payment instead of CreateAccount operation.
func (bot *Bot) Pay(destAddress string, amount string, isFundRequest bool) (*horizon.TransactionSuccess, error) {
	channel := make(chan interface{})
	shouldReadChannel, result, err := bot.lockedPay(channel, destAddress, amount, isFundRequest)
	if !shouldReadChannel {
		return result, err
=======
// Pay funds the account at `destAddress`
func (bot *Bot) Pay(destAddress string) (*horizon.TransactionSuccess, error) {
	channel := make(chan TxResult)
	err := bot.lockedPay(channel, destAddress)
	if err != nil {
		return nil, err
>>>>>>> horizon-v0.15.3
	}

	v := <-channel
	return v.maybeTransactionSuccess, v.maybeErr
}

<<<<<<< HEAD
func (bot *Bot) lockedPay(channel chan interface{}, destAddress string, amount string, isFundRequest bool) (bool, *horizon.TransactionSuccess, error) {
=======
func (bot *Bot) lockedPay(channel chan TxResult, destAddress string) error {
>>>>>>> horizon-v0.15.3
	bot.lock.Lock()
	defer bot.lock.Unlock()

	err := bot.checkSequenceRefresh()
	if err != nil {
		return err
	}

	signed, err := bot.makeTx(destAddress, amount, isFundRequest)
	if err != nil {
		return err
	}

	go bot.SubmitTransaction(bot, channel, signed)
	return nil
}

// AsyncSubmitTransaction should be passed into the bot
func AsyncSubmitTransaction(bot *Bot, channel chan TxResult, signed string) {
	result, err := bot.Horizon.SubmitTransaction(signed)
	if err != nil {
		switch e := err.(type) {
		case *horizon.Error:
			bot.checkHandleBadSequence(e)
		}

		channel <- TxResult{
			maybeTransactionSuccess: nil,
			maybeErr:                err,
		}
	} else {
		channel <- TxResult{
			maybeTransactionSuccess: &result,
			maybeErr:                nil,
		}
	}
}

func (bot *Bot) checkHandleBadSequence(err *horizon.Error) {
	resCode, e := err.ResultCodes()
	isTxBadSeqCode := e == nil && resCode.TransactionCode == "tx_bad_seq"
	if !isTxBadSeqCode {
		return
	}
	bot.forceRefreshSequence = true
}

// establish initial sequence if needed
func (bot *Bot) checkSequenceRefresh() error {
	if bot.sequence != 0 && !bot.forceRefreshSequence {
		return nil
	}
	return bot.refreshSequence()
}

func (bot *Bot) makeTx(destAddress string, amount string, isFundRequest bool) (string, error) {
	ops := []b.TransactionMutator{
		b.SourceAccount{AddressOrSeed: bot.Secret},
		b.Sequence{Sequence: bot.sequence + 1},
		b.Network{Passphrase: bot.Network},
	}

	if isFundRequest {
		ops = append(ops, b.Payment(
			b.Destination{AddressOrSeed: destAddress},
			b.NativeAmount{Amount: amount}))
	} else {
		ops = append(ops, b.CreateAccount(
			b.Destination{AddressOrSeed: destAddress},
			b.NativeAmount{Amount: amount}))
	}

	txn, err := b.Transaction(ops...)

	if err != nil {
		return "", errors.Wrap(err, "Error building a transaction")
	}

	txs, err := txn.Sign(bot.Secret)
	if err != nil {
		return "", errors.Wrap(err, "Error signing a transaction")
	}

	base64, err := txs.Base64()

	// only increment the in-memory sequence number if we are going to submit the transaction, while we hold the lock
	if err == nil {
		bot.sequence++
	}
	return base64, err
}

// refreshes the sequence from the bot account
func (bot *Bot) refreshSequence() error {
	botAccount, err := bot.Horizon.LoadAccount(bot.address())
	if err != nil {
		bot.sequence = 0
		return err
	}

	seq, err := strconv.ParseInt(botAccount.Sequence, 10, 64)
	if err != nil {
		bot.sequence = 0
		return err
	}

	bot.sequence = uint64(seq)
	bot.forceRefreshSequence = false
	return nil
}

func (bot *Bot) address() string {
	kp := keypair.MustParse(bot.Secret)
	return kp.Address()
}
