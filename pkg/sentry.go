// Package sentry provides backward compatibility for the original API.
//
// Deprecated: Import the specific packages instead:
//   - github.com/denysvitali/tesla-sentry-viewer/pkg/event
//   - github.com/denysvitali/tesla-sentry-viewer/pkg/clip
package sentry

import (
	"os"

	"github.com/denysvitali/tesla-sentry-viewer/pkg/clip"
	"github.com/denysvitali/tesla-sentry-viewer/pkg/event"
)

// Event is an alias for backward compatibility.
//
// Deprecated: Use event.Event instead.
type Event = event.Event

// ParseEvent wraps event.ParseEvent for backward compatibility.
//
// Deprecated: Use event.ParseEvent instead.
func ParseEvent(directory string) (*Event, error) {
	return event.ParseEvent(directory)
}

// FilesByType wraps clip.FilesByType for backward compatibility.
//
// Deprecated: Use clip.FilesByType instead.
func FilesByType(directory string, entries []os.DirEntry) (map[string][]string, error) {
	return clip.FilesByType(directory, entries)
}

// GetFileType wraps clip.GetFileType for backward compatibility.
//
// Deprecated: Use clip.GetFileType instead.
func GetFileType(name string) string {
	return clip.GetFileType(name)
}
