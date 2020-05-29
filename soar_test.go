package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalSoar(t *testing.T) {
	str, err := LocalSoar()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(str) > 0)
}
func TestTempSQLFile(t *testing.T) {
	str, err := tempSQLFile("select * from car;")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, len(str) > 0)
	err = removeTempSQLFile(str)
	assert.Equal(t, nil, err)
}
