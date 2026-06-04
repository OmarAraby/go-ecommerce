package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	productapp "github.com/OmarAraby/go-ecommerce/internal/application/services/product"
)

var _ productapp.FileStorage = (*LocalStorage)(nil)

type LocalStorage struct {
	BaseDir string // e.g. "uploads"
	BaseURL string // e.g. "/uploads"
}

func NewLocalStorage(baseDir, baseURL string) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir, BaseURL: baseURL}
}

func (s *LocalStorage) Save(_ context.Context, dir, filename string, r io.Reader) (string, error) {
	dirPath := filepath.Join(s.BaseDir, dir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}
	dst, err := os.Create(filepath.Join(dirPath, filename))
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, r); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return fmt.Sprintf("%s/%s/%s", s.BaseURL, dir, filename), nil
}

func (s *LocalStorage) Delete(_ context.Context, url string) error {
	rel := strings.TrimPrefix(url, s.BaseURL+"/")
	return os.Remove(filepath.Join(s.BaseDir, rel))
}
