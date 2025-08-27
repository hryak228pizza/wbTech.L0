package validation

import (
	_ "fmt"
	"testing"

	"github.com/hryak228pizza/wbTech.L0/internal/generator"
)

func newTestValidator() *Validate {
	return NewValidator()
}

func TestValidateOrder_Valid(t *testing.T) {

	v := newTestValidator()

	order := generator.NewOrder()

	err := v.ValidateOrder(&order)
	if err != nil {
		t.Errorf("expect valid, got: %v", err)
	}
}

func TestValidateOrder_InvalidPhone(t *testing.T) {

	v := newTestValidator()

	order := generator.NewOrder()
	order.Delivery.Phone = "79522663535"

	err := v.ValidateOrder(&order)
	if err == nil {
		t.Error("expect invalid phone, got valid")
	}
}

func TestValidateOrder_InvalidEmail(t *testing.T) {

	v := newTestValidator()

	order := generator.NewOrder()
	order.Delivery.Email = "mail@@mail"

	err := v.ValidateOrder(&order)
	if err == nil {
		t.Error("expect invalid email, got valid")
	}
}

