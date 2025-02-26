package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getLastSegmentAsInt(t *testing.T) {
	assert.Equal(t, 151, getLastSegmentAsInt("http://example.com/151"))
	assert.Equal(t, 151, getLastSegmentAsInt("http://example.com/151"))
}
