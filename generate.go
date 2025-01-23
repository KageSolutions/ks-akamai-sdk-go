package ksakamaisdkgo

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/KakashiHatake324/mockjs"
)

// request akamai dynamic data
func (r *AkamaiSdkInstance) RequestDynamic(script string) error {
	compressed, _ := gzipEncodeHTML(script)
	r.UpdateScript(compressed)
	requestData, err := structToReader(r.dynamicRequest)
	if err != nil {
		return err
	}
	if r.verbose {
		log.Println("DYNAMIC REQUEST DATA:", r.dynamicRequest)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v%s", r.apiDynamicUrl, r.akamaiVersion), requestData)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if r.verbose {
		log.Println("DYNAMIC REQUEST RESPONSE:", string(body))
	}
	var responseData = make(map[string]interface{})
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	if _, ok := responseData["success"]; !ok {
		return errors.New("could not successfully generate akamai dynamic data, success not in body")
	}
	if !responseData["success"].(bool) {
		return errors.New("could not successfully generate akamai dynamic data")
	}
	// set the first sensor as false since it wont be the first anymore
	r.sensorRequest.DynamicData = responseData["data"].(map[string]interface{})
	return err
}

// request akamai sensor data
func (r *AkamaiSdkInstance) RequestSensor() (*AkamaiResponse, error) {
	requestData, err := structToReader(r.sensorRequest)
	if err != nil {
		return nil, err
	}
	if r.verbose {
		log.Println("SENSOR REQUEST DATA:", r.sensorRequest)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sensor/v%s", r.apiSensorUrl, r.akamaiVersion), requestData)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if r.verbose {
		log.Println("SENSOR REQUEST RESPONSE:", string(body))
	}
	var responseData AkamaiResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}
	if !responseData.Success {
		return nil, errors.New("could not successfully generate akamai sensor")
	}
	// set the first sensor as false since it wont be the first anymore
	r.sensorRequest.First = false
	r.SensorData = responseData.Data
	return &responseData, err
}

// request akamai pixel data
func (r *AkamaiSdkInstance) RequestPixel() (*AkamaiResponse, error) {
	requestData, err := structToReader(r.pixelRequest)
	if err != nil {
		return nil, err
	}
	if r.verbose {
		log.Println("PIXEL REQUEST DATA:", r.pixelRequest)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/pixel", r.apiSensorUrl), requestData)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if r.verbose {
		log.Println("PIXEL REQUEST RESPONSE:", string(body))
	}
	var responseData AkamaiResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}
	if !responseData.Success {
		return nil, errors.New("could not successfully generate akamai pixel")
	}
	r.PixelData = responseData.Data
	return &responseData, err
}

func gzipEncodeHTML(html string) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write([]byte(html)); err != nil {
		return "", err
	}

	if err := gz.Close(); err != nil {
		return "", err
	}

	return mockjs.InitWindow().Btoa(string(buf.Bytes())), nil
}
