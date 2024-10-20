package ksakamaisdkgo

type AkamaiSdkInstance struct {
	akamaiVersion  string
	dynamicRequest dynamicRequest
	sensorRequest  akamaiRequest
	pixelRequest   pixelRequest

	WebsiteUrl    string
	AkamaiWebUrl  string
	ContainsPixel bool
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
