package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func (c *HttpClient) Get(url string, dest any, headers http.Header) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = headers
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("请求失败，状态码：" + strconv.Itoa(resp.StatusCode))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if dest == nil {
		return nil
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}

	return nil
}

func (c *HttpClient) Post(url string, body any, dest any, headers http.Header) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = headers
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("请求失败，状态码：" + strconv.Itoa(resp.StatusCode))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if dest == nil {
		return nil
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}

	return nil
}

func (c *HttpClient) GetDownload(url string, dstPath string, headers http.Header) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = headers
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("请求失败，状态码：" + strconv.Itoa(resp.StatusCode))
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	return err

}

var Request *HttpClient

func init() {
	Request = &HttpClient{client: &http.Client{Timeout: 30 * time.Second}}
}
