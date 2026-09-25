package domain

import (
	"encoding/json"
	"testing"
)

func TestAmount_Add(t *testing.T) {
	a1 := NewAmount("NGN", 1000, 2)
	a2 := NewAmount("NGN", 500, 2)

	res, err := a1.Add(a2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Units.Int64() != 1500 {
		t.Errorf("expected 1500, got %d", res.Units.Int64())
	}
}

func TestAmount_Add_CurrencyMismatch(t *testing.T) {
	a1 := NewAmount("NGN", 1000, 2)
	a2 := NewAmount("MAD", 500, 2)

	_, err := a1.Add(a2)
	if err != ErrCurrencyMismatch {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestAmount_Sub(t *testing.T) {
	a1 := NewAmount("NGN", 1000, 2)
	a2 := NewAmount("NGN", 500, 2)

	res, err := a1.Sub(a2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Units.Int64() != 500 {
		t.Errorf("expected 500, got %d", res.Units.Int64())
	}
}

func TestAmount_Sub_NegativeAmount(t *testing.T) {
	a1 := NewAmount("NGN", 500, 2)
	a2 := NewAmount("NGN", 1000, 2)

	_, err := a1.Sub(a2)
	if err != ErrNegativeAmount {
		t.Errorf("expected ErrNegativeAmount, got %v", err)
	}
}

func TestAmount_Mul(t *testing.T) {
	a1 := NewAmount("NGN", 1000, 2)
	res := a1.Mul(3)
	if res.Units.Int64() != 3000 {
		t.Errorf("expected 3000, got %d", res.Units.Int64())
	}
}

func TestAmount_Div(t *testing.T) {
	a1 := NewAmount("NGN", 1000, 2)
	
	// Exact division
	res, err := a1.Div(2, RoundDown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Units.Int64() != 500 {
		t.Errorf("expected 500, got %d", res.Units.Int64())
	}

	// 1000 / 3 = 333.33...
	// RoundDown -> 333
	res, _ = a1.Div(3, RoundDown)
	if res.Units.Int64() != 333 {
		t.Errorf("expected 333, got %d", res.Units.Int64())
	}

	// RoundUp -> 334
	res, _ = a1.Div(3, RoundUp)
	if res.Units.Int64() != 334 {
		t.Errorf("expected 334, got %d", res.Units.Int64())
	}

	// RoundHalfUp -> 333 (because remainder 1 is less than half of 3)
	res, _ = a1.Div(3, RoundHalfUp)
	if res.Units.Int64() != 333 {
		t.Errorf("expected 333, got %d", res.Units.Int64())
	}

	// 1000 / 6 = 166.66...
	// RoundHalfUp -> 167 (because remainder 4 is > half of 6)
	res, _ = a1.Div(6, RoundHalfUp)
	if res.Units.Int64() != 167 {
		t.Errorf("expected 167, got %d", res.Units.Int64())
	}
}

func TestAmount_Cmp(t *testing.T) {
	a1 := NewAmount("NGN", 1000, 2)
	a2 := NewAmount("NGN", 500, 2)

	cmp, _ := a1.Cmp(a2)
	if cmp <= 0 {
		t.Errorf("expected a1 > a2")
	}
}

func TestAmount_JSON(t *testing.T) {
	a1 := NewAmount("NGN", 1000, 2)
	
	data, err := json.Marshal(a1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedJSON := `{"currency":"NGN","units":"1000","scale":2}`
	if string(data) != expectedJSON {
		t.Errorf("expected %s, got %s", expectedJSON, string(data))
	}

	var a2 Amount
	err = json.Unmarshal(data, &a2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a2.Currency != "NGN" || a2.Units.Int64() != 1000 || a2.Scale != 2 {
		t.Errorf("unmarshal failed, got: %+v", a2)
	}
}
