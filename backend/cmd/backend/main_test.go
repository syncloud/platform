package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerGood(t *testing.T) {

	backend := Backend()

	assert.NotNil(t, backend.Master)

}
