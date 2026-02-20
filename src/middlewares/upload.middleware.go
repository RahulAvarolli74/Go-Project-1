package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"recipe-api/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/jpg":  true,
}

func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			c.Next()
			return
		}
		defer file.Close()

		contentType := header.Header.Get("Content-Type")
		if !AllowedImageTypes[contentType] {
			utils.ErrorResponse(c, http.StatusBadRequest,
				fmt.Sprintf("Invalid file type: %s. Only JPEG and PNG are allowed", contentType))
			c.Abort()
			return
		}

		maxSize := int64(10 << 20)
		if ms := os.Getenv("MAX_UPLOAD_SIZE"); ms != "" {
			if parsed, err := strconv.ParseInt(ms, 10, 64); err == nil {
				maxSize = parsed << 20
			}
		}
		if header.Size > maxSize {
			utils.ErrorResponse(c, http.StatusBadRequest,
				fmt.Sprintf("File too large. Maximum size is %dMB", maxSize>>20))
			c.Abort()
			return
		}

		uploadDir := os.Getenv("UPLOAD_DIR")
		if uploadDir == "" {
			uploadDir = "./public/temp"
		}
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create upload directory")
			c.Abort()
			return
		}

		ext := strings.ToLower(filepath.Ext(header.Filename))
		tempFilename := fmt.Sprintf("raw_%s%s", uuid.New().String()[:8], ext)
		tempPath := filepath.Join(uploadDir, tempFilename)

		if err := c.SaveUploadedFile(header, tempPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save uploaded file")
			c.Abort()
			return
		}

		c.Set("uploadedFilePath", tempPath)
		c.Next()
	}
}
