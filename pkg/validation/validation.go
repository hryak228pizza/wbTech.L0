package validation

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hryak228pizza/wbTech.L0/internal/model"
	"github.com/nyaruka/phonenumbers"
)

var (
	phoneRegexp = regexp.MustCompile(`^(?:\+7|8)[-\s]?\d{3}[-\s]?\d{3}[-\s]?\d{2}[-\s]?\d{2}$`)
	emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	zipRegexp   = regexp.MustCompile(`^\d+$`)
	nameRegexp  = regexp.MustCompile(`^[a-zA-Z]+$`)
)

type Validate struct {
	validate *validator.Validate
}

// NewValidator creates new validator with registrations
func NewValidator() *Validate {

	// init validator
	newValidator := validator.New()
	// registrate reciever name validation
	newValidator.RegisterValidation("name", func(fl validator.FieldLevel) bool {
		name := strings.TrimSpace(fl.Field().String())
		if name == "" {
			return false
		}
		return nameRegexp.MatchString(name)
	})
	// registrate reciever zip code validation
	newValidator.RegisterValidation("zip", func(fl validator.FieldLevel) bool {
		zip := strings.TrimSpace(fl.Field().String())
		if zip == "" {
			return false
		}
		return zipRegexp.MatchString(zip)
	})
	// registrate reciever city validation
	newValidator.RegisterValidation("city", func(fl validator.FieldLevel) bool {
		city := strings.TrimSpace(fl.Field().String())
		return city != ""
	})
	// registrate reciever address validation
	newValidator.RegisterValidation("address", func(fl validator.FieldLevel) bool {
		address := strings.TrimSpace(fl.Field().String())
		return address != ""
	})
	// registrate reciever region validation
	newValidator.RegisterValidation("region", func(fl validator.FieldLevel) bool {
		region := strings.TrimSpace(fl.Field().String())
		return region != ""
	})
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
		return emailRegexp.MatchString(email)
	})
	return &Validate{validate: newValidator}
}

// IsValid returns true if object is valid
func (v *Validate) ValidateOrder(order *model.Order) error {

	// validate order data
	return v.validate.Struct(order)
}
