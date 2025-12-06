// Package clip provides Tesla Sentry clip file discovery and organization.
package clip

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
)

// Camera type constants for Tesla Sentry Mode videos.
const (
	CameraFront         = "front"
	CameraBack          = "back"
	CameraLeftRepeater  = "left_repeater"
	CameraRightRepeater = "right_repeater"
)

// VideoFileRegex matches Tesla Sentry video filenames.
// Format: YYYY-MM-DD_HH-MM-SS-{camera_type}.mp4
var VideoFileRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}-(.*?)\.mp4$`)

// GetFileType extracts the camera type from a Tesla Sentry video filename.
// Returns empty string if the filename doesn't match the expected pattern.
func GetFileType(name string) string {
	if !VideoFileRegex.MatchString(name) {
		return ""
	}

	matches := VideoFileRegex.FindStringSubmatch(name)
	return matches[1]
}

// FilesByType categorizes video files in a directory by their camera type.
// Returns a map where keys are camera types (front, back, left_repeater, right_repeater)
// and values are slices of absolute file paths.
func FilesByType(directory string, entries []os.DirEntry) (map[string][]string, error) {
	filesByType := map[string][]string{}

	for _, dirEntry := range entries {
		if dirEntry.Type() == 0 {
			// It's a file
			ft := GetFileType(dirEntry.Name())
			if ft != "" {
				if _, ok := filesByType[ft]; !ok {
					filesByType[ft] = []string{}
				}

				if directory == "" {
					filesByType[ft] = append(filesByType[ft], dirEntry.Name())
				} else {
					relFilePath := path.Join(directory, dirEntry.Name())
					absFilePath, err := filepath.Abs(relFilePath)
					if err != nil {
						return nil, err
					}
					filesByType[ft] = append(filesByType[ft], absFilePath)
				}
			}
		}
	}

	return filesByType, nil
}
