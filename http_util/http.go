package http_util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpPostJson(reqBody []byte,url string) (respBody []byte, err error) {

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
