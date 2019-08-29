package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data Payload `json:"data"`
	Err  string  `json:"err"`
}

type ApiRequest struct {
	Req Payload `json:"data"`
}

func HttpRequestFileServer(message Payload) (res Payload, err error) {
	host := GetConfigClientFileHttp()
	//fmt.Fprintf(os.Stderr, "%v\n", host)
	if host[len(host)-1] != '/' {
		host = host + "/"
	}

	url := host + "api"

	request := ApiRequest{
		Req: message,
	}
	b, _ := json.Marshal(request)
	body := bytes.NewBuffer(b)

	req, err2 := http.NewRequest("POST", url, body)
	if err2 != nil {
		err = fmt.Errorf("http.NewRequest error: %v", err2)
		return
	}

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		panic(err2)
	}
	defer resp.Body.Close()

	var respBody Response
	if errDec := json.NewDecoder(resp.Body).Decode(&respBody); errDec != nil {
		err = fmt.Errorf("%v\n", errDec)
		return
	}

	res = respBody.Data

	return
}
