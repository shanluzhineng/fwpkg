package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type GenReqWithBodyTypeFunc func(pathUrl string, method string) (*http.Request, error)

func GenReqWithFormBody(body map[string]interface{}, headers map[string]string) GenReqWithBodyTypeFunc {
	return func(pathUrl, method string) (*http.Request, error) {
		formBody := url.Values{}
		for key, value := range body {
			formBody.Set(key, fmt.Sprintf("%v", value))
		}
		req, err := http.NewRequest(method, pathUrl, strings.NewReader(formBody.Encode()))
		if err != nil {
			return nil, err
		}
		//req.Form = formBody
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if len(headers) > 0 {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		return req, nil
	}
}

func GenReqWithJsonBody(body map[string]interface{}, headers map[string]string) GenReqWithBodyTypeFunc {
	return func(pathUrl, method string) (*http.Request, error) {
		var buf *bytes.Buffer = nil
		if body != nil {
			jsonByte, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewBuffer(jsonByte)
		}
		req, err := http.NewRequest(method, pathUrl, buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if len(headers) > 0 {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		return req, nil
	}
}

func GenReqForGet(headers map[string]string) GenReqWithBodyTypeFunc {
	return func(pathUrl, method string) (*http.Request, error) {
		req, err := http.NewRequest(method, pathUrl, nil)
		if err != nil {
			return nil, err
		}
		if len(headers) > 0 {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		return req, nil
	}
}

/**
 * 封装请求及body解析。json参数用 map 传递，response body 也使用map读取。非200会报错，err中包含异常body。
 */
func DoHttpRequest(url string, method string, token string, bodyTypeFunc GenReqWithBodyTypeFunc) (map[string]interface{}, error) {
	var (
		err error
		req *http.Request
	)
	//if body != nil {
	//	bodyByte, err := json.Marshal(body)
	//	if err != nil {
	//		return nil, err
	//	}
	//	reader = bytes.NewReader(bodyByte)
	//	req, err = http.NewRequest(strings.ToUpper(method), url, reader)
	//} else {
	//}
	method = strings.ToUpper(method)
	if bodyTypeFunc != nil {
		req, err = bodyTypeFunc(url, method)
		if err != nil {
			return nil, err
		}
	}
	if req == nil {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}
	//req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("readBody failed: %s", err.Error())
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %s, body: %s", response.Status, string(respBody))
	}
	// log.Logger.Debug(fmt.Sprintf("response body: %s, url: %s", string(respBody), url))
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func PostJson(url string, body map[string]interface{}) (map[string]interface{}, error) {
	return DoHttpRequest(url, http.MethodPost, "", GenReqWithJsonBody(body, nil))
}

func PostJsonWithHeader(url string, body map[string]interface{}, headerMap map[string]string) (map[string]interface{}, error) {
	return DoHttpRequest(url, http.MethodPost, "", GenReqWithJsonBody(body, headerMap))
}

func PostForm(url string, body map[string]interface{}) (map[string]interface{}, error) {
	return DoHttpRequest(url, http.MethodPost, "", GenReqWithFormBody(body, nil))
}

func PostFormWithHeader(url string, body map[string]interface{}, headerMap map[string]string) (map[string]interface{}, error) {
	return DoHttpRequest(url, http.MethodPost, "", GenReqWithFormBody(body, headerMap))
}

func GetUrl(url string) (map[string]interface{}, error) {
	return DoHttpRequest(url, http.MethodGet, "", nil)
}

func GetUrlWithHeader(url string, headerMap map[string]string) (map[string]interface{}, error) {
	return DoHttpRequest(url, http.MethodGet, "", GenReqForGet(headerMap))
}

// GetIP returns request real ip.
func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", fmt.Errorf("no valid ip found")
}
