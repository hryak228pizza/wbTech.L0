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

func TestValidateOrder_InvalidName(t *testing.T) {
	v := newTestValidator()
	order := generator.NewOrder()
	order.Delivery.Name = ""
	err := v.ValidateOrder(&order)
	if err == nil {
		t.Error("expect invalid name, got valid")
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

func TestValidateOrder_InvalidZip(t *testing.T) {
	v := newTestValidator()
	order := generator.NewOrder()
	order.Delivery.Zip = ""
	err := v.ValidateOrder(&order)
	if err == nil {
		t.Error("expect invalid zip, got valid")
	}
}

func TestValidateOrder_InvalidCity(t *testing.T) {
	v := newTestValidator()
	order := generator.NewOrder()
	order.Delivery.City = ""
	err := v.ValidateOrder(&order)
	if err == nil {
		t.Error("expect invalid city, got valid")
	}
}

func TestValidateOrder_InvalidAddress(t *testing.T) {
	v := newTestValidator()
	order := generator.NewOrder()
	order.Delivery.Address = ""
	err := v.ValidateOrder(&order)
	if err == nil {
		t.Error("expect invalid address, got valid")
	}
}

func TestValidateOrder_InvalidRegion(t *testing.T) {
	v := newTestValidator()
	order := generator.NewOrder()
	order.Delivery.Region = ""
	err := v.ValidateOrder(&order)
	if err == nil {
		t.Error("expect invalid region, got valid")
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

