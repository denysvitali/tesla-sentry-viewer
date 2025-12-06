package clip

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "front camera",
			filename: "2022-03-30_07-11-08-front.mp4",
			expected: CameraFront,
		},
		{
			name:     "back camera",
			filename: "2022-03-30_07-11-08-back.mp4",
			expected: CameraBack,
		},
		{
			name:     "left repeater camera",
			filename: "2022-03-30_07-11-08-left_repeater.mp4",
			expected: CameraLeftRepeater,
		},
		{
			name:     "right repeater camera",
			filename: "2022-03-30_07-11-08-right_repeater.mp4",
			expected: CameraRightRepeater,
		},
		{
			name:     "invalid filename",
			filename: "random-file.mp4",
			expected: "",
		},
		{
			name:     "not mp4",
			filename: "2022-03-30_07-11-08-front.txt",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileType(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilesByType(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test video files
	files := []string{
		"2022-03-30_07-11-08-front.mp4",
		"2022-03-30_07-11-09-front.mp4",
		"2022-03-30_07-11-08-back.mp4",
		"2022-03-30_07-11-08-left_repeater.mp4",
		"event.json", // Should be ignored
	}

	for _, f := range files {
		err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644)
		require.NoError(t, err)
	}

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)

	result, err := FilesByType(tmpDir, entries)
	require.NoError(t, err)

	assert.Len(t, result[CameraFront], 2)
	assert.Len(t, result[CameraBack], 1)
	assert.Len(t, result[CameraLeftRepeater], 1)
	assert.NotContains(t, result, "event.json")
}

func TestFilesByType_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)

	result, err := FilesByType(tmpDir, entries)
	require.NoError(t, err)
	assert.Empty(t, result)
}
