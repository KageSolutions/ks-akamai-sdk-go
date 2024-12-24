package ksakamaisdkgo

type AkamaiSdkInstance struct {
	akamaiVersion  string
	dynamicRequest dynamicRequest
	sensorRequest  akamaiRequest
	pixelRequest   pixelRequest
	apiSensorUrl   string
	apiDynamicUrl  string

	WebsiteUrl    string
	AkamaiWebUrl  string
	SensorData    string
	PixelData     string
	PixelVersion  string
	ContainsPixel bool
	verbose       bool
}

type dynamicRequest struct {
	ApiKey string `json:"apiKey"`
	Script string `json:"script"`
}

type akamaiRequest struct {
	ApiKey      string                 `json:"apiKey"`
	Ua          string                 `json:"ua"`
	PageURL     string                 `json:"pageUrl"`
	Abck        string                 `json:"_abck"`
	BmSz        string                 `json:"bm_sz"`
	First       bool                   `json:"first"`
	ForceMact   bool                   `json:"forceMact"`
	DynamicData map[string]interface{} `json:"dynamic,omitempty"`
}

type pixelRequest struct {
	ApiKey    string `json:"apiKey"`
	Ua        string `json:"ua"`
	PageURL   string `json:"pageUrl"`
	ScriptVal string `json:"scriptVal"`
	PixelID   string `json:"pixelId"`
}

type AkamaiResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}
