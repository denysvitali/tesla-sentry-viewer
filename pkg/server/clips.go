package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/denysvitali/tesla-sentry-viewer/pkg/clip"
	"github.com/denysvitali/tesla-sentry-viewer/pkg/event"
	"github.com/gin-gonic/gin"
)

type ClipResponse struct {
	Event     event.Event         `json:"event"`
	ClipFiles map[string][]string `json:"clipFiles"`
}

func (s *Server) getClip(c *gin.Context) {
	clipId := c.Param("clip_id")
	if !directoryRegex.MatchString(clipId) {
		s.logger.Warnf("invalid clip id provided: %s", clipId)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid clip id",
			"code":  "INVALID_CLIP_ID",
		})
		return
	}

	// Check if clip exists
	clipPath := path.Join(s.dir, clipId)
	fileInfo, err := os.Stat(clipPath)
	if err != nil {
		s.logger.Warnf("clip not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("clip %s not found", clipId),
		})
		return
	}

	if !fileInfo.IsDir() {
		s.logger.Warnf("clip is not a directory")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("clip %s is invalid", clipId),
		})
		return
	}

	// Get event
	evt, err := event.ParseEvent(clipPath)
	if err != nil {
		s.logger.Warnf("unable to parse clip directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to parse clip directory",
		})
		return
	}

	dirEntries, err := os.ReadDir(clipPath)
	if err != nil {
		s.logger.Warnf("unable to read dir: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to read dir",
		})
		return
	}

	filesByType, err := clip.FilesByType("", dirEntries)
	if err != nil {
		s.logger.Warnf("unable to get files by type: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to process files by type",
		})
		return
	}

	clipResponse := ClipResponse{
		Event:     *evt,
		ClipFiles: filesByType,
	}

	c.JSON(http.StatusOK, clipResponse)
}

func (s *Server) getClipFile(c *gin.Context) {
	clipId := c.Param("clip_id")
	fileName := c.Param("file_name")
	if !strings.HasSuffix(fileName, ".mp4") {
		s.logger.Warnf("invalid file type requested: %s", fileName)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the clip file must be an .mp4",
			"code":  "INVALID_FILE_TYPE",
		})
		return
	}

	filePath := path.Join(s.dir, clipId, fileName)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		s.logger.Warnf("clip file doesn't exist: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "clip file doesn't exist"})
		return
	}

	if fileInfo.IsDir() {
		s.logger.Warn("clip file is dir")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid clip file"})
		return
	}

	c.File(filePath)
}

func (s *Server) getClipThumb(c *gin.Context) {
	clipId := c.Param("clip_id")
	if !directoryRegex.MatchString(clipId) {
		s.logger.Warnf("invalid clip id provided for thumb: %s", clipId)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid clip id",
			"code":  "INVALID_CLIP_ID",
		})
		return
	}

	filePath := path.Join(s.dir, clipId, "thumb.png")
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		s.logger.Warnf("thumb file doesn't exist: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "thumb file doesn't exist"})
		return
	}

	if fileInfo.IsDir() {
		s.logger.Warn("thumb file is dir")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid thumb file"})
		return
	}

	c.File(filePath)
}
