package clients

import (
	"doc-precise-rag/knowledge/config"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var CosClient *cos.Client

func init() {
	u, _ := url.Parse(*config.AppConfig.CosConfig.Domain)
	b := &cos.BaseURL{
		BucketURL: u,
	}

	CosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  *config.AppConfig.CosConfig.SecretId,
			SecretKey: *config.AppConfig.CosConfig.SecretKey,
		},
	})

}
