package printer

import (
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()

    validate.RegisterValidation("itemtype", func(fl validator.FieldLevel) bool {
        val := ItemType(fl.Field().String())
        switch val {
        case Text, QRCode, Blank, Line, Cut:
            return true
        }
        return false
    })

    validate.RegisterValidation("alignment", func(fl validator.FieldLevel) bool {
        val := Alignment(fl.Field().String())
        switch val {
        case AlignLeft, AlignCenter, AlignRight:
            return true
        }
        return false
    })
}

func ValidatePrintRequest(pr PrintRequest) error {
    return validate.Struct(pr)
}