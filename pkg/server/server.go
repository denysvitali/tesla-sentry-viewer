package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"regexp"
)

type Server struct {
	logger *logrus.Logger
	dir    string
}

var directoryRegex = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}_\\d{2}-\\d{2}-\\d{2}")

func New(directory string) (*Server, error) {
	fileInfo, err := os.Stat(directory)
	if err != nil {
		return nil, fmt.Errorf("unable to stat: %v", err)
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", directory)
	}

	s := Server{
		logger: logrus.New(),
		dir:    directory,
	}

	return &s, nil
}

func (s *Server) Listen(addr string) error {
	e := gin.New()
	e.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	e.GET("/api/v1/clips", s.getClips)
	e.GET("/api/v1/clips/:clip_id", s.validateClip, s.getClip)
	e.GET("/api/v1/clips/:clip_id/thumb", s.validateClip, s.getClipThumb)
	e.GET("/api/v1/clips/:clip_id/:file_name", s.validateClip, s.getClipFile)
	return e.Run(addr)
}

func (s *Server) getClips(context *gin.Context) {
	// Get all directories
	dirEntries, err := os.ReadDir(s.dir)
	if err != nil {
		s.logger.Warnf("unable to get directory entries: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Unable to read dir: %v", err),
		})
		return
	}

	var validDirectories []string
	for _, v := range dirEntries {
		if directoryRegex.MatchString(v.Name()) {
			_, err = os.Stat(path.Join(s.dir, v.Name(), "event.json"))
			if err != nil {
				s.logger.Warnf("skipping directory %s: %v", v.Name(), err)
				continue
			}
			validDirectories = append(validDirectories, v.Name())
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"events": validDirectories,
	})
}

func (s *Server) SetLogger(logger *logrus.Logger) {
	if s.logger != nil {
		s.logger = logger
	}
}

func (s *Server) validateClip(c *gin.Context) {
	clipId := c.Param("clip_id")
	if !directoryRegex.MatchString(clipId) {
		s.logger.Warnf("invalid clip id provided: %s", clipId)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid clip id",
		})
		return
	}
}
