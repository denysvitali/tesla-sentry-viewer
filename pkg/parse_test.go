package sentry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileType(t *testing.T) {
	assert.Equal(t, "right_repeater", GetFileType("2022-03-30_07-11-08-right_repeater.mp4"))
	assert.Equal(t, "front", GetFileType("2022-02-11_07-00-00-front.mp4"))
	assert.Equal(t, "", GetFileType("1-02-11_07-00-00-front.mp4"))
}
