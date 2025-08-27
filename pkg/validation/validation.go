package validation

import (
	"net/mail"
	"regexp"

	"github.com/hryak228pizza/wbTech.L0/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

var (
	phoneRegexp = regexp.MustCompile(`^(?:\+7|8)[-\s]?\d{3}[-\s]?\d{3}[-\s]?\d{2}[-\s]?\d{2}$`)
    emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

type Validate struct {
	validate *validator.Validate
}

// NewValidator creates new validator with registrations
func NewValidator() *Validate {

	// init validator
	newValidator := validator.New()
	// registrate phone validation
    newValidator.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
        phone := fl.Field().String()
        if _, err := phonenumbers.Parse(phone, "RU"); err != nil {
			return false
		}		
        return phoneRegexp.MatchString(phone)
    })
	// registrate email validation
	newValidator.RegisterValidation("email", func(fl validator.FieldLevel) bool {
        email := fl.Field().String()
		if _, err := mail.ParseAddress(email); err != nil {
			return false
		}
        return emailRegex.MatchString(email)
    })	
	return &Validate{ validate: newValidator }
}

// IsValid returns true if object is valid
func (v *Validate) ValidateOrder(order *model.Order) error {

	// validate order data
	return v.validate.Struct(order)
}
