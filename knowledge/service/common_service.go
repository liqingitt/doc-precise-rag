package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"doc-precise-rag/knowledge/config"

	"github.com/jmoiron/sqlx"
	"github.com/tencentyun/cos-go-sdk-v5"
)

const (
	uploadPrefix = "doc-anchor/pdf"
	urlExpire    = 15 * time.Minute
)

var ErrInvalidUpload = errors.New("invalid upload request")

type CommonService struct {
	db  *sqlx.DB
	cos *cos.Client
}

type UploadURLResult struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	FileURL   string `json:"file_url"`
	Filename  string `json:"filename"`
}

func NewCommonService(cos *cos.Client) *CommonService {
	return &CommonService{
		cos: cos,
	}
}

func (s *CommonService) GetUploadURL(ctx context.Context, filename string) (*UploadURLResult, error) {
	safeName, err := sanitizePDFName(filename)
	if err != nil {
		return nil, err
	}

	objectKey := fmt.Sprintf("%s/%s_%s", uploadPrefix, newUUID(), safeName)
	signedURL, err := s.cos.Object.GetPresignedURL2(ctx, http.MethodPut, objectKey, urlExpire, nil)
	if err != nil {
		return nil, err
	}

	domain := strings.TrimRight(*config.AppConfig.CosConfig.Domain, "/")
	return &UploadURLResult{
		UploadURL: signedURL.String(),
		ObjectKey: objectKey,
		FileURL:   domain + "/" + objectKey,
		Filename:  safeName,
	}, nil
}

func sanitizePDFName(name string) (string, error) {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" || name == "." || name == ".." {
		return "", fmt.Errorf("%w: 文件名无效", ErrInvalidUpload)
	}
	if !strings.EqualFold(filepath.Ext(name), ".pdf") {
		return "", fmt.Errorf("%w: 仅支持 PDF 文件", ErrInvalidUpload)
	}

	stem := strings.TrimSuffix(name, filepath.Ext(name))
	var b strings.Builder
	for _, r := range stem {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '-', r == '_', r == '.':
			b.WriteRune(r)
		case r == ' ':
			b.WriteByte('_')
		}
	}

	safe := strings.Trim(b.String(), "._")
	if safe == "" {
		safe = "file"
	}
	if len([]rune(safe)) > 80 {
		safe = string([]rune(safe)[:80])
	}
	return safe + ".pdf", nil
}

func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
