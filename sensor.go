package ksakamaisdkgo

func (t *AkamaiSdkInstance) GenerateSensor(max int, postReq func() error) error {
	for l := 0; l < max; l++ {
		if _, err := t.RequestSensor(); err != nil {
			return err
		}

		if err := postReq(); err != nil {
			return err
		}
	}
	return nil
}

func (t *AkamaiSdkInstance) HandlePixel(postReq func() error) error {
	if _, err := t.RequestPixel(); err != nil {
		return err
	}

	if err := postReq(); err != nil {
		return err
	}
	return nil
}
