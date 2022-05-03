package sentry

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var videoFileRegex = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}_\\d{2}-\\d{2}-\\d{2}-(.*?)\\.mp4$")

func ProcessVideo(directory string) error {
	event, err := ParseEvent(directory)
	if err != nil {
		return err
	}
	fmt.Printf("Event %s in %s (%s, %s)\n", event.Reason, event.City, event.EstLat, event.EstLon)
	// Merge Video
	err = mergeVideo(directory)
	if err != nil {
		return err
	}
	return nil
}

func ParseEvent(directory string) (*Event, error) {
	f, err := os.Open(path.Join(directory, "event.json"))
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	var event Event
	err = dec.Decode(&event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

type sortByName []os.DirEntry

func (s sortByName) Len() int {
	return len(s)
}

func (s sortByName) Less(i, j int) bool {
	return strings.Compare(s[i].Name(), s[j].Name()) == -1
}

func (s sortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var _ sort.Interface = (*sortByName)(nil)

func mergeVideo(directory string) error {
	// Scan directory for videos
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	sort.Sort(sortByName(entries))

	filesByType, err := FilesByType(directory, entries)
	if err != nil {
		return err
	}

	fmt.Printf("file types = %v\n", filesByType)
	for k, v := range filesByType {
		err = mergeFiles(v, k)
		if err != nil {
			return err
		}
	}

	return nil
}

func FilesByType(directory string, entries []os.DirEntry) (map[string][]string, error) {
	filesByType := map[string][]string{}
	for _, dirEntry := range entries {
		if dirEntry.Type() == 0 {
			// It's a file!
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

func mergeFiles(inFiles []string, name string) error {
	f, err := os.CreateTemp("", "sentry-viewer-*")
	if err != nil {
		return nil
	}
	defer os.Remove(f.Name())

	// Write
	for _, v := range inFiles {
		_, err = fmt.Fprintln(f, fmt.Sprintf("file %s", v))
		if err != nil {
			return err
		}
	}

	outputFileName := path.Join("output", fmt.Sprintf("%s.mp4", name))
	cmd := exec.Command("ffmpeg",
		"-f",
		"concat",
		"-safe",
		"0",
		"-i",
		f.Name(),
		"-c",
		"copy",
		outputFileName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return nil
}

func GetFileType(name string) string {
	// 2022-03-30_07-11-08-right_repeater.mp4
	if !videoFileRegex.MatchString(name) {
		return ""
	}

	matches := videoFileRegex.FindStringSubmatch(name)
	return matches[1]
}
