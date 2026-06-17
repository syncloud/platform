package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordValidator_Validate(t *testing.T) {
	validator := NewPasswordValidator()
	assert.NoError(t, validator.Validate("Password1"))
	assert.NoError(t, validator.Validate("regularpass123"))
	assert.Error(t, validator.Validate("short1"))
	assert.Error(t, validator.Validate("alllettersonly"))
	assert.Error(t, validator.Validate("12345678"))
}
