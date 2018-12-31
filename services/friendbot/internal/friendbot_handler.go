package internal

import (
	"net/http"
	"net/url"

	"github.com/kinecosystem/go/clients/horizon"
	"github.com/kinecosystem/go/strkey"
	"github.com/kinecosystem/go/support/render/hal"
	"github.com/kinecosystem/go/support/render/problem"

	"strconv"
)

// FriendbotHandler causes an account at `Address` to be created.
type FriendbotHandler struct {
	Friendbot *Bot
}

// Handle is a method that implements http.HandlerFunc
func (handler *FriendbotHandler) Handle(w http.ResponseWriter, r *http.Request) {
	result, err := handler.doHandle(r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	hal.Render(w, *result)
}

// doHandle is just a convenience method that returns the object to be rendered
func (handler *FriendbotHandler) doHandle(r *http.Request) (*horizon.TransactionSuccess, error) {
	err := handler.checkEnabled()
	if err != nil {
		return nil, err
	}

	err = r.ParseForm()
	if err != nil {
		return nil, err
	}

	address, err := handler.loadAddress(r)
	if err != nil {
		return nil, problem.MakeInvalidFieldProblem("addr", err)
	}

	amount, err := handler.loadAmount(r)
	if err != nil {
		amount = handler.Friendbot.StartingBalance
	} else {
		s, err := strconv.ParseFloat(handler.Friendbot.StartingBalance, 64)
		a, err := strconv.ParseFloat(amount, 64)

		if err != nil || a > s {
			amount = handler.Friendbot.StartingBalance
		}
	}

	isFundRequest := r.URL.Path == "/fund"

	return handler.loadResult(address, amount, isFundRequest)
}

func (handler *FriendbotHandler) checkEnabled() error {
	if handler.Friendbot != nil {
		return nil
	}

	return &problem.P{
		Type:   "friendbot_disabled",
		Title:  "Friendbot is disabled",
		Status: http.StatusForbidden,
		Detail: "Friendbot is disabled on this network. Contact the server administrator if you believe this to be in error.",
	}
}

func (handler *FriendbotHandler) loadAddress(r *http.Request) (string, error) {
	address := r.Form.Get("addr")
	unescaped, err := url.QueryUnescape(address)
	if err != nil {
		return unescaped, err
	}

	_, err = strkey.Decode(strkey.VersionByteAccountID, unescaped)
	return unescaped, err
}

func (handler *FriendbotHandler) loadAmount(r *http.Request) (string, error) {
	amount := r.Form.Get("amount")
	unescaped, err := url.QueryUnescape(amount)
	if err != nil {
		return unescaped, err
	}

	return unescaped, err
}

func (handler *FriendbotHandler) loadResult(address string, amount string, isFundRequest bool) (*horizon.TransactionSuccess, error) {
	result, err := handler.Friendbot.Pay(address, amount, isFundRequest)
	switch e := err.(type) {
	case horizon.Error:
		return result, e.Problem.ToProblem()
	case *horizon.Error:
		return result, e.Problem.ToProblem()
	}
	return result, err
}
