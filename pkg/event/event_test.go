package event

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEvent(t *testing.T) {
	// Create a temporary directory with an event.json file
	tmpDir := t.TempDir()
	eventJSON := `{
		"timestamp": "2022-03-30T07:11:08",
		"city": "San Francisco",
		"est_lat": "37.7749",
		"est_lon": "-122.4194",
		"reason": "sentry_aware_object_detection",
		"camera": "front"
	}`

	err := os.WriteFile(filepath.Join(tmpDir, "event.json"), []byte(eventJSON), 0644)
	require.NoError(t, err)

	event, err := ParseEvent(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "2022-03-30T07:11:08", event.Timestamp)
	assert.Equal(t, "San Francisco", event.City)
	assert.Equal(t, "37.7749", event.EstLat)
	assert.Equal(t, "-122.4194", event.EstLon)
	assert.Equal(t, "sentry_aware_object_detection", event.Reason)
	assert.Equal(t, "front", event.Camera)
}

func TestParseEvent_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := ParseEvent(tmpDir)
	assert.Error(t, err)
}

func TestParseEvent_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "event.json"), []byte("invalid json"), 0644)
	require.NoError(t, err)

	_, err = ParseEvent(tmpDir)
	assert.Error(t, err)
}
