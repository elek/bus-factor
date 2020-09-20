package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadAliasFile(t *testing.T) {
	d, err := readAliasFile("testdata/test_alias")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(d))
	assert.Equal(t, "bsd", d["asd"])
}
