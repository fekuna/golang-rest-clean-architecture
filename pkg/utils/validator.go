package utils

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	fmt.Println("validate init running lho")
	validate = validator.New()
}

// Validate struct fields
func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}
