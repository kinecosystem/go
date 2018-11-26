package amount_test

import (
	"testing"

	"github.com/kinecosystem/go/amount"
	"github.com/kinecosystem/go/xdr"
)

var Tests = []struct {
	S     string
	I     xdr.Int64
	valid bool
}{
	{"100.00000", 10000000, true},
	{"-100.00000", -10000000, true},
	{"100.00001", 10000001, true},
	{"123.00001", 12300001, true},
	{"123.000001", 0, false},
	{"922337203685.47758", 92233720368547758, true},
	{"922337203685477.58", 0, false},
	{"92233720368600", 0, false},
	{"-922337203685.47758", -92233720368547758, true},
	{"-922337203685477.58", 0, false},
	{"-92233720368600", 0, false},
	{"100000000000000.000", 0, false},
	{"100000000000000", 0, false},
}

func TestParse(t *testing.T) {
	for _, v := range Tests {
		o, err := amount.Parse(v.S)
		if !v.valid && err == nil {
			t.Errorf("expected err for input %s", v.S)
			continue
		}
		if v.valid && err != nil {
			t.Errorf("couldn't parse %s: %v", v.S, err)
			continue
		}

		if o != v.I {
			t.Errorf("%s parsed to %d, not %d", v.S, o, v.I)
		}
	}
}

func TestString(t *testing.T) {
	for _, v := range Tests {
		if !v.valid {
			continue
		}

		o := amount.String(v.I)

		if o != v.S {
			t.Errorf("%d stringified to %s, not %s", v.I, o, v.S)
		}
	}
}
