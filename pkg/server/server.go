package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/denysvitali/tesla-sentry-viewer/pkg/config"
)

const SoftwareName = "tesla-sentry-viewer"

var Version = "dev"

type Server struct {
	logger *logrus.Logger
	dir    string
	config *config.Config
}

var directoryRegex = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}_\\d{2}-\\d{2}-\\d{2}")

func New(directory string) (*Server, error) {
	// Backward compatibility wrapper
	cfg := config.NewConfig()
	return NewWithConfig(directory, cfg)
}

func NewWithConfig(directory string, cfg *config.Config) (*Server, error) {
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
		config: cfg,
	}

	return &s, nil
}

func (s *Server) Listen(addr string) error {
	e := gin.New()

	// Security middleware
	e.Use(gin.Recovery())
	e.Use(s.securityHeaders())

	// CORS configuration
	e.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	// Add health check endpoint
	e.GET("/health", s.healthCheck)

	// Add version endpoint
	e.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": Version,
			"status":  "healthy",
		})
	})

	e.GET("/api/v1/clips", s.getClips)
	e.GET("/api/v1/clips/:clip_id", s.validateClip, s.getClip)
	e.GET("/api/v1/clips/:clip_id/thumb", s.validateClip, s.getClipThumb)
	e.GET("/api/v1/clips/:clip_id/:file_name", s.validateClip, s.getClipFile)

	// Add options for CORS preflight
	e.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return e.Run(addr)
}

func (s *Server) securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Cache control for API responses
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"version":     Version,
		"service":     SoftwareName,
	})
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
