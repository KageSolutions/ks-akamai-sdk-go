package ksakamaisdkgo

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// init new sdk instance
func NewAkSdkInstance(apiKey, version, website, sensorUrl, dynamicUrl string, firstFeature, verbose bool) *AkamaiSdkInstance {
	return &AkamaiSdkInstance{
		akamaiVersion:  version,
		apiSensorUrl:   sensorUrl,
		apiDynamicUrl:  dynamicUrl,
		dynamicRequest: dynamicRequest{ApiKey: apiKey},
		sensorRequest:  akamaiRequest{ApiKey: apiKey, First: firstFeature},
		pixelRequest:   pixelRequest{ApiKey: apiKey},
		WebsiteUrl:     website,
		verbose:        verbose,
	}
}

func (p *AkamaiSdkInstance) UpdateScript(script string) {
	p.dynamicRequest.Script = script
}

func (p *AkamaiSdkInstance) UpdatePageUrl(pageURL string) {
	p.sensorRequest.PageURL = pageURL
	p.pixelRequest.PageURL = pageURL
}

func (p *AkamaiSdkInstance) UpdateForceMact(should bool) {
	p.sensorRequest.ForceMact = should
}

func (p *AkamaiSdkInstance) UpdateUserAgent(userAgent string) {
	p.sensorRequest.Ua = userAgent
	p.pixelRequest.Ua = userAgent
}

func (p *AkamaiSdkInstance) UpdateAbck(abck string) {
	p.sensorRequest.Abck = abck
}

func (p *AkamaiSdkInstance) UpdateBmsz(bmsz string) {
	p.sensorRequest.BmSz = bmsz
}

func (p *AkamaiSdkInstance) UpdatePixelScriptValue(scriptVal string) {
	p.pixelRequest.ScriptVal = scriptVal
}

func (p *AkamaiSdkInstance) UpdatePixelId(pixelId string) {
	p.pixelRequest.PixelID = pixelId
}

func (p *AkamaiSdkInstance) UpdatePixelVersion(pixelVersion string) {
	p.PixelVersion = pixelVersion
}

func structToReader(data interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonData), nil
}

func (a *AkamaiSdkInstance) IsChallenged(abck string) bool {
	Qb := a.DestructureCookie(abck)
	kb := make([]string, len(Qb)+1)

	if Qb != nil {
		for wb := 0; wb < len(Qb); wb++ {
			Wb := Qb[wb]
			if len(Wb) > 0 {
				Zb := fmt.Sprintf("%s%s", Wb[1], Wb[2])
				if len(kb) <= Wb[6].(int) {
					kb = append(kb, Zb)
				} else {
					kb[Wb[6].(int)] = Zb
				}
			}
		}
	}
	if len(kb) > 2 {
		return true
	}
	return false
}

func (a *AkamaiSdkInstance) DestructureCookie(Abck string) [][]interface{} {
	q5 := [][]interface{}{}
	Q5, _ := url.QueryUnescape(Abck)

	if Q5 != "" {
		k5 := strings.Split(Q5, "~")
		if len(k5) >= 5 {
			w5 := k5[0]
			W5 := strings.Split(k5[4], "||")
			if len(W5) > 0 {
				for Z5 := 0; Z5 < len(W5); Z5++ {
					R5 := strings.Split(W5[Z5], "-")
					if len(R5) >= 5 {
						n5, _ := strconv.Atoi(R5[0])
						r5 := R5[1]
						z5, _ := strconv.Atoi(R5[2])
						O5, _ := strconv.Atoi(R5[3])
						t5, _ := strconv.Atoi(R5[4])
						d5 := 1
						if len(R5) >= 6 {
							d5, _ = strconv.Atoi(R5[5])
						}
						x5 := []interface{}{n5, w5, r5, z5, O5, t5, d5}
						if d5 == 2 {
							q5 = append([][]interface{}{x5}, q5...)
						} else {
							q5 = append(q5, x5)
						}
					}
				}
			}
		}
	}

	return q5
}

func (t *AkamaiSdkInstance) ParseAkamaiBody(body string) error {
	webUrl, err := url.Parse(t.WebsiteUrl)
	if err != nil {
		return err
	}
	pixelValueRegex := regexp.MustCompile(`bazadebezolkohpepadr="(\d.*?)"`)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		sectionHTML, exists := s.Attr("src")
		if exists {
			if len(strings.Split(sectionHTML, "/")) >= 6 && len(strings.Split(sectionHTML, "/")) <= 10 {
				if !strings.Contains(sectionHTML, "http") && !strings.Contains(sectionHTML, ".js") && !strings.Contains(sectionHTML, ".mjs") {
					if !strings.Contains(sectionHTML, "?") {
						if len(t.AkamaiWebUrl) < len(fmt.Sprintf("https://%s%s", webUrl.Host, sectionHTML)) {
							t.AkamaiWebUrl = fmt.Sprintf("https://%s%s", webUrl.Host, sectionHTML)
						}
					}
				}
			}
		}
	})

	if t.AkamaiWebUrl == "" {
		return errors.New("couldnt find akamai web url")
	}
	if bazaCheck, err := regexp.MatchString("bazadebezolkohpepadr", body); err != nil {
		return err
	} else {
		if bazaCheck {
			t.ContainsPixel = true
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
			doc.Find("script").Each(func(i int, s *goquery.Selection) {
				src, exists := s.Attr("src")
				if exists {
					if strings.Contains(src, "akam/13/") {
						t.UpdatePixelVersion(strings.Split(src, "/")[len(strings.Split(src, "/"))-1])
					}
				}
			})

			pixelId := pixelValueRegex.FindString(body)
			Id := strings.Split(pixelId, "=")[1]
			t.UpdatePixelId(strings.ReplaceAll(Id, "\"", ""))

		}
	}
	return nil
}

func (t *AkamaiSdkInstance) ParsePixelScript(body string) error {

	var gIndex, gVal string
	exp := regexp.MustCompile(`(?m)g=_(.*),m`)

	for _, match := range exp.FindAllString(body, -1) {
		exp := regexp.MustCompile("[0-9]+")
		gIndex = exp.FindAllString(match, -1)[0]
	}

	exp2 := regexp.MustCompile(`(?m)var _=[ []"(.*)];`)
	for _, match2 := range exp2.FindAllString(body, -1) {
		rep := strings.NewReplacer("var _ = [", "", "];", "", `"`, "", "\u0020", "")
		res := rep.Replace(match2)
		arr := strings.Split(res, ",")

		intVar, err := strconv.Atoi(gIndex)
		if err != nil {
			return errors.New("error parsing pixel")
		}

		rep2 := strings.NewReplacer("\\", "", "x", "", "", "")
		gVal = rep2.Replace(arr[intVar])
		decodedString, err := hex.DecodeString(gVal)
		if err != nil {
			return errors.New("error parsing pixel")
		}
		t.UpdatePixelScriptValue(string(decodedString))
	}
	return nil
}
