package alidayu

import (
	"crypto/hmac"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	HTTP_URL  string = "http://gw.api.taobao.com/router/rest"
	HTTPS_URL string = "https://eco.taobao.com/router/rest"

	SUCCESS_TAG_JSON string = `"success":true`
	SUCCESS_TAG_XML  string = `<success>true</success>`
)

var (
	DefCommonParams map[string]string = map[string]string{
		"format":      "json",
		"v":           "2.0",
		"sign_method": "md5",
	}
)

type Client struct {
	AppKey    string
	AppSecret string
	UseHTTPS  bool
	http.Client
}

func (c *Client) IsValid() bool {
	if c.AppKey == "" || c.AppSecret == "" {
		return false
	}

	return true
}

func (c *Client) UpdateCommonParams(params map[string]string) {
	for k, v := range DefCommonParams {
		if _, ok := params[k]; !ok {
			params[k] = v
		}
	}

	// Set App Key.
	params["app_key"] = c.AppKey

	// Set timestamp.
	t := time.Now()
	params["timestamp"] = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func GetSortedQueryStr(params map[string]string) (sortedQueryStr string) {
	keys := []string{}
	for k, _ := range params {
		keys = append(keys, k)
	}

	str := ""
	sort.Strings(keys)
	for _, k := range keys {
		str += k + params[k]
	}
	return str
}

func (c *Client) SignMD5(params map[string]string) (signature string) {
	// Sort keys by name to generate MD5 hash.
	// See details at:
	//     http://open.taobao.com/doc2/detail.htm?articleId=101617&docType=1&treeId=1
	str := fmt.Sprintf("%s%s%s", c.AppSecret, GetSortedQueryStr(params), c.AppSecret)

	return fmt.Sprintf("%X", md5.Sum([]byte(str)))

}

func (c *Client) SignHMAC(params map[string]string) (signature string) {
	// Sort keys by name to generate HMAC_MD5 hash.
	// See details at:
	//     http://open.taobao.com/doc2/detail.htm?articleId=101617&docType=1&treeId=1
	str := GetSortedQueryStr(params)

	// HMAC_MD5
	mac := hmac.New(md5.New, []byte(c.AppSecret))
	mac.Write([]byte(str))
	return fmt.Sprintf("%X", mac.Sum(nil))
}

func (c *Client) MakeRequestBody(params map[string]string) (body io.Reader, err error) {
	if !c.IsValid() {
		return nil, errors.New("Empty App Key or App Secret.")
	}

	values := url.Values{}

	// Check "method".
	if _, ok := params["method"]; !ok {
		return nil, errors.New("No method specified.")
	}

	// Update Common Params.
	c.UpdateCommonParams(params)

	// Check "format".
	if params["format"] != "json" && params["format"] != "xml" {
		return nil, errors.New(fmt.Sprintf("format error: %v", params["format"]))
	}

	sign := ""
	switch params["sign_method"] {
	case "md5":
		sign = c.SignMD5(params)
	case "hmac":
		sign = c.SignHMAC(params)
	default:
		return nil, errors.New("Incorrect sign_method.")
	}

	params["sign"] = sign

	for k, v := range params {
		values.Set(k, v)
	}

	return strings.NewReader(values.Encode()), nil
}

func (c *Client) Post(params map[string]string) (resp *http.Response, err error) {
	var body io.Reader

	if body, err = c.MakeRequestBody(params); err != nil {
		return nil, err
	}

	urlStr := ""
	if c.UseHTTPS {
		urlStr = HTTPS_URL
	} else {
		urlStr = HTTP_URL
	}

	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return c.Do(req)
}

func (c *Client) Exec(params map[string]string) (success bool, result string, err error) {
	resp, err := c.Post(params)
	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	result = string(data)

	switch params["format"] {
	case "json":
		success = strings.Contains(result, SUCCESS_TAG_JSON)
	case "xml":
		success = strings.Contains(result, SUCCESS_TAG_XML)
	}

	return success, result, nil
}
