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
	// HTTP_URL is for HTTP REST API URL.
	HTTP_URL string = "http://gw.api.taobao.com/router/rest"
	// HTTPS_URL is for HTTPS REST API URL.
	HTTPS_URL string = "https://eco.taobao.com/router/rest"

	// SUCCESS_TAG_JSON is success tag for JSON in response string.
	SUCCESS_TAG_JSON string = `"success":true`
	// SUCCESS_TAG_XML is sucess tag for XML in response string
	SUCCESS_TAG_XML string = `<success>true</success>`
)

var (
	// DefCommonParams is default common parameters of alidayu APIs
	DefCommonParams = map[string]string{
		"format":      "json", // Response Format("json" or "xml")
		"v":           "2.0",  // API version("2.0")
		"sign_method": "md5",  // Sign method("md5" or "hmac")
	}
)

// Client contains app key and secret and provides method to post HTTP request like http.Request.
type Client struct {
	AppKey    string // App Key
	AppSecret string // App Secret
	UseHTTPS  bool   // Use HTTPS URL or not
	http.Client
}

// IsValid checks if a client is valid.
func (c *Client) IsValid() bool {
	if c.AppKey == "" || c.AppSecret == "" {
		return false
	}

	return true
}

// UpdateCommonParams updates the given parameters by merge default common parameters.
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

// GetSortedQueryStr sorts the keys of parameters and generate sorted query string which can be used to sign(MD5 or HMAC_MD5).
func GetSortedQueryStr(params map[string]string) (sortedQueryStr string) {
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}

	str := ""
	sort.Strings(keys)
	for _, k := range keys {
		str += k + params[k]
	}
	return str
}

// SignMD5 gets the sorted query string and sign it using MD5.
func (c *Client) SignMD5(params map[string]string) (signature string) {
	// Sort keys by name to generate MD5 hash.
	// See details at:
	//     http://open.taobao.com/doc2/detail.htm?articleId=101617&docType=1&treeId=1
	str := fmt.Sprintf("%s%s%s", c.AppSecret, GetSortedQueryStr(params), c.AppSecret)

	return fmt.Sprintf("%X", md5.Sum([]byte(str)))

}

// SignHMAC gets the sorted query string and sign it using HMAC_MD5.
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

// MakeRequestBody makes the HTTP request body by given parameters for each REST API.
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

// Post does the HTTP post for each REST API.
//
//   Params:
//     params: map that contains parameters of REST API. See official docs to fill the parameters.
//   Returns:
//     resp: HTTP Response. Do not forget to call resp.Body.Close() after use.
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

// Exec executes the REST API and get the response. It's a wrapper of Post().
//
//   Params:
//     params: map that contains parameters of REST API. See official docs to fill the parameters.
//   Returns:
//     success: If REST API succeeds.
//     result: Raw response string of REST API.
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
