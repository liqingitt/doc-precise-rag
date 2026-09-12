package clients

import (
	"doc-precise-rag/knowledge/config"
	"errors"
	"net/http"
)

type Mineru struct {
	baseURL          string
	parserFileURL    string
	getTaskStatusURL string
	headers          http.Header
}

type ParserFileStatusData struct {
	TaskID     *string `json:"task_id"`
	State      *string `json:"state"`
	FullZipURL *string `json:"full_zip_url"`
	ErrorMsg   *string `json:"err_msg"`
}

type ParserFileStatusResp struct {
	Code    *int                  `json:"code"`
	Msg     *string               `json:"message"`
	TraceID *string               `json:"trace_id"`
	Data    *ParserFileStatusData `json:"data"`
}

type CreateParserFileTaskData struct {
	TaskID *string `json:"task_id"`
}

type CreateParserFileTaskResp struct {
	Code    *int                      `json:"code"`
	Msg     *string                   `json:"message"`
	TraceID *string                   `json:"trace_id"`
	Data    *CreateParserFileTaskData `json:"data"`
}

func (c *Mineru) CreateParserFileTask(fileURL string) (*CreateParserFileTaskResp, error) {
	var resp CreateParserFileTaskResp
	err := Request.Post(c.parserFileURL, map[string]any{
		"url":            fileURL,
		"model_version":  "vlm",
		"is_ocr":         true,
		"enable_formula": true,
		"enable_table":   true,
	}, &resp, c.headers)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Mineru) GetParserFileStatus(taskID string) (*ParserFileStatusResp, error) {
	var resp ParserFileStatusResp
	err := Request.Get(c.getTaskStatusURL+"/"+taskID, &resp, c.headers)
	if err != nil {
		return nil, err
	}

	if *resp.Code != 0 {
		return nil, errors.New(*resp.Msg)
	}
	if *resp.Code != 0 {
		return nil, errors.New(*resp.Msg)
	}
	return &resp, nil
}

var MineruClient *Mineru

func init() {
	MineruClient = &Mineru{
		baseURL:          *config.AppConfig.AiConfig.MineruConfig.BaseURL,
		parserFileURL:    *config.AppConfig.AiConfig.MineruConfig.BaseURL + "/api/v4/extract/task",
		getTaskStatusURL: *config.AppConfig.AiConfig.MineruConfig.BaseURL + "/api/v4/extract/task",
		headers: http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + *config.AppConfig.AiConfig.MineruConfig.Token},
		},
	}
}
