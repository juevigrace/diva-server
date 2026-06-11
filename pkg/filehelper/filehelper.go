package filehelper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/pkg/errs"
)

const DEFAULT_MAX_FILE_SIZE = 5 << 20

var defaultAllowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

type FileHelper struct {
	RootDir string

	AllowedFileTypes map[string]string
	MaxFileSize      int
}

func NewFileHelper(
	rootDir string,
	allowedFileTypes map[string]string,
	maxFileSize int,
) *FileHelper {
	var dir string = "./uploads"
	if rootDir != "" {
		dir = rootDir
	}

	var maxfs = DEFAULT_MAX_FILE_SIZE
	if maxFileSize > 0 {
		maxfs = maxFileSize
	}

	var allowed map[string]string = defaultAllowedImageTypes
	if len(allowedFileTypes) > 0 {
		allowed = allowedFileTypes
	}

	return &FileHelper{
		RootDir:          dir,
		AllowedFileTypes: allowed,
		MaxFileSize:      maxfs,
	}
}

func (h *FileHelper) SaveImage(
	file multipart.File,
	size int64,
	contentType, saveDir string,
) (string, error) {
	if file == nil {
		return "", errs.ErrFileRequired
	}
	defer file.Close()

	ext, ok := h.AllowedFileTypes[contentType]
	if !ok {
		return "", errs.ErrUnsupportedImage
	}

	if size > int64(h.MaxFileSize) {
		return "", errs.ErrFileTooLarge
	}

	fullPath := filepath.Join(h.RootDir, saveDir)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", err
	}

	filename := uuid.New().String() + ext
	dest := filepath.Join(fullPath, filename)

	dst, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s/%s", strings.TrimPrefix(saveDir, "/"), filename), nil
}

func (h *FileHelper) DeleteFile(path string) error {
	if path == "" {
		return nil
	}
	return os.Remove(filepath.Join(h.RootDir, strings.TrimPrefix(path, "/")))
}
