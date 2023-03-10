package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	MethodPost    = "POST"
	MethodGet     = "GET"
	HttpPrefix    = "http://"
	ContentType   = "application/json"
	ErrorNotPlato = "没有plato节点"
)

func getPlatoHostUri(name string) (string, error) {
	platoHost := Consul.NextPlatoServer()
	if len(platoHost) == 0 {
		LoggerStd.Error(ErrorNotPlato)
		return "", errors.New(ErrorNotPlato)
	}
	url := HttpPrefix + platoHost + "/" + name
	return url, nil
}

func DecodeSlider(pictureCount int, param map[string]any) (resp map[string]any, err error) {
	uri := "DecodeSlider1"
	if pictureCount == 2 {
		uri = "DecodeSlider2"
	}
	url, err := getPlatoHostUri(uri)
	if err != nil {
		return nil, err
	}
	return HandleHttpRequest(MethodPost, url, param)
}

func HttpOCRImage(param map[string]any) (resp map[string]any, err error) {
	url, err := getPlatoHostUri("OCRImage")
	if err != nil {
		return nil, err
	}
	return HandleHttpRequest(MethodPost, url, param)
}

func HttpPost(url string, param map[string]any) (resp map[string]any, err error) {
	return HandleHttpRequest(MethodPost, url, param)
}

func HttpGet(url string, param map[string]any) (resp map[string]any, err error) {
	return HandleHttpRequest(MethodGet, url, param)
}

func HttpGetFile(url string, param map[string]any, filePath string) (err error) {
	byteArr, err := HandleHttpRequestByte(MethodGet, url, param)
	if err != nil {
		return err
	}
	//normalize bytes to file
	_, err = Write0(filePath, byteArr)
	if err != nil {
		return err
	}
	return nil
}

func HandleHttpRequest(method string, url string, param map[string]any) (resp map[string]any, err error) {
	body, err := HandleHttpRequestByte(method, url, param)
	//Sugar.Infow("HandleHttpRequest", "url", url, "param", param, "response", body)
	if err != nil {
		return nil, err
	}
	resp = make(map[string]any)
	err = json.Unmarshal(body, &resp)
	return resp, err
}

func HandleHttpRequestByte(method string, url string, param map[string]any) (resp []byte, err error) {
	if param == nil {
		param = make(map[string]any)
	}
	jsonParam, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonParam))
	req.Header.Add("Content-Type", ContentType)
	response, err := client.Do(req)
	if err != nil {
		//Sugar.Errorw("HandleHttpRequest", "url", url, "param", param, "error", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return body, err
}
