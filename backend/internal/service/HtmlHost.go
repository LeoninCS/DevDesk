package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrHTMLTooLarge    = errors.New("html file is too large")
	ErrHTMLInvalidExt  = errors.New("only .html or .htm files are allowed")
	ErrHTMLMissingFile = errors.New("html file is required")
)

const DefaultMaxHTMLSize = 1 * 1024 * 1024

type HTMLHostConfig struct {
	BaseDir      string
	MaxSizeBytes int64
}

type HTMLHostService struct {
	baseDir string
	maxSize int64
}

func NewHTMLHostService(cfg HTMLHostConfig) (*HTMLHostService, error) {
	if cfg.BaseDir == "" {
		cfg.BaseDir = "hosted_html"
	}
	if cfg.MaxSizeBytes <= 0 {
		cfg.MaxSizeBytes = DefaultMaxHTMLSize
	}

	if err := os.MkdirAll(cfg.BaseDir, 0o755); err != nil {
		return nil, err
	}

	return &HTMLHostService{
		baseDir: cfg.BaseDir,
		maxSize: cfg.MaxSizeBytes,
	}, nil
}

func (s *HTMLHostService) BaseDir() string {
	return s.baseDir
}

func (s *HTMLHostService) MaxSizeBytes() int64 {
	return s.maxSize
}

func (s *HTMLHostService) SaveHTML(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", ErrHTMLMissingFile
	}

	if fileHeader.Size > s.maxSize {
		return "", ErrHTMLTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".html" && ext != ".htm" {
		return "", ErrHTMLInvalidExt
	}

	filename := GetHash(12) + ".html"
	dst := filepath.Join(s.baseDir, filename)

	if err := s.writeFileWithLimit(fileHeader, dst); err != nil {
		return "", err
	}

	return filename, nil
}

func (s *HTMLHostService) writeFileWithLimit(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	tmpFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	limit := s.maxSize + 1
	n, err := io.Copy(tmpFile, io.LimitReader(src, limit))
	if err != nil {
		os.Remove(dst)
		return err
	}
	if n > s.maxSize {
		os.Remove(dst)
		return ErrHTMLTooLarge
	}
	return nil
}
