package ksakamaisdkgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/KakashiHatake324/mockjs"
)

// request akamai dynamic data
func (r *AkamaiSdkInstance) RequestDynamic(script string) error {
	r.UpdateScript(mockjs.Window.Btoa(script))
	requestData, err := structToReader(r.dynamicRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v%s", DynamicApiUrl, r.akamaiVersion), requestData)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sensor/v%s", ApiUrl, r.akamaiVersion), requestData)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/pixel", ApiUrl), requestData)
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
