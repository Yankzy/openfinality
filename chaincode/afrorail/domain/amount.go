package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

type RoundingMode int

const (
	RoundHalfUp RoundingMode = iota
	RoundDown
	RoundUp
)

// Amount represents a deterministic fixed-point monetary value.
// It avoids any floating-point arithmetic.
type Amount struct {
	Currency string
	Units    *big.Int
	Scale    uint32
}

var (
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrNegativeAmount   = errors.New("amount cannot be negative")
	ErrScaleMismatch    = errors.New("scale mismatch")
)

// NewAmount creates a new Amount.
func NewAmount(currency string, units int64, scale uint32) Amount {
	return Amount{
		Currency: currency,
		Units:    big.NewInt(units),
		Scale:    scale,
	}
}

// Add adds another amount to this amount. Returns a new Amount.
func (a Amount) Add(b Amount) (Amount, error) {
	if a.Currency != b.Currency {
		return Amount{}, ErrCurrencyMismatch
	}
	if a.Scale != b.Scale {
		return Amount{}, ErrScaleMismatch
	}

	result := new(big.Int).Add(a.Units, b.Units)
	return Amount{
		Currency: a.Currency,
		Units:    result,
		Scale:    a.Scale,
	}, nil
}

// Sub subtracts another amount from this amount. Returns a new Amount.
func (a Amount) Sub(b Amount) (Amount, error) {
	if a.Currency != b.Currency {
		return Amount{}, ErrCurrencyMismatch
	}
	if a.Scale != b.Scale {
		return Amount{}, ErrScaleMismatch
	}

	result := new(big.Int).Sub(a.Units, b.Units)
	if result.Sign() < 0 {
		return Amount{}, ErrNegativeAmount
	}

	return Amount{
		Currency: a.Currency,
		Units:    result,
		Scale:    a.Scale,
	}, nil
}

// Mul multiplies the amount by an integer factor.
func (a Amount) Mul(factor int64) Amount {
	result := new(big.Int).Mul(a.Units, big.NewInt(factor))
	return Amount{
		Currency: a.Currency,
		Units:    result,
		Scale:    a.Scale,
	}
}

// Div divides the amount by a divisor with a specific rounding mode.
func (a Amount) Div(divisor int64, mode RoundingMode) (Amount, error) {
	if divisor == 0 {
		return Amount{}, errors.New("division by zero")
	}

	d := big.NewInt(divisor)
	q, r := new(big.Int).QuoRem(a.Units, d, new(big.Int))

	// Adjust quotient based on rounding mode
	if r.Sign() != 0 {
		switch mode {
		case RoundUp:
			q.Add(q, big.NewInt(1))
		case RoundHalfUp:
			// check if 2*r >= d
			r2 := new(big.Int).Mul(r, big.NewInt(2))
			if r2.Cmp(d) >= 0 {
				q.Add(q, big.NewInt(1))
			}
		case RoundDown:
			// Do nothing, big.Int.QuoRem rounds towards zero (truncate)
			// Assuming positive amounts, this is RoundDown
		}
	}

	return Amount{
		Currency: a.Currency,
		Units:    q,
		Scale:    a.Scale,
	}, nil
}

// Cmp compares two amounts.
func (a Amount) Cmp(b Amount) (int, error) {
	if a.Currency != b.Currency {
		return 0, ErrCurrencyMismatch
	}
	if a.Scale != b.Scale {
		return 0, ErrScaleMismatch
	}
	return a.Units.Cmp(b.Units), nil
}

// MarshalJSON implements the json.Marshaler interface for canonical serialization.
func (a Amount) MarshalJSON() ([]byte, error) {
	unitsStr := "0"
	if a.Units != nil {
		unitsStr = a.Units.String()
	}
	return json.Marshal(struct {
		Currency string `json:"currency"`
		Units    string `json:"units"`
		Scale    uint32 `json:"scale"`
	}{
		Currency: a.Currency,
		Units:    unitsStr,
		Scale:    a.Scale,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *Amount) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Currency string `json:"currency"`
		Units    string `json:"units"`
		Scale    uint32 `json:"scale"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	units := new(big.Int)
	if _, ok := units.SetString(aux.Units, 10); !ok {
		return fmt.Errorf("invalid big.Int string: %s", aux.Units)
	}

	a.Currency = aux.Currency
	a.Units = units
	a.Scale = aux.Scale
	return nil
}
