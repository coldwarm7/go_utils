package http_util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func HTTPPostJson(reqBody []byte,url string) (respBody []byte, err error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client do err: ", err)
		return
	}
	defer resp.Body.Close()
	respBody, _ = ioutil.ReadAll(resp.Body)

	return
}

func HTTPGet(reqURL string, args map[string]string) ([]byte, error) {
	URL, err := url.Parse(strings.Trim(reqURL, "/"))
	if err != nil {
		return nil, err
	}

	query := URL.Query()
	if nil != args {
		for k, v := range args {
			query.Add(k, v)
		}
	}

	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Http statusCode:%d", resp.StatusCode)
	}

	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return result, nil
}