package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordValidator_Validate(t *testing.T) {
	validator := NewPasswordValidator()
	assert.Nil(t, validator.Validate("Password1"))
	assert.Nil(t, validator.Validate("regularpass123"))
	assert.NotNil(t, validator.Validate("short1"))
	assert.NotNil(t, validator.Validate("alllettersonly"))
	assert.NotNil(t, validator.Validate("12345678"))
}
