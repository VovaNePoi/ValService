package externalapi

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testTaskMod/internal/models"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type IClient interface {
	GetJSONCurrency(params MarketAPIparams) (*models.ValCursJSON, error)
}

// client realiztion for requests to cbr API by http
type clientStruct struct {
	client    *http.Client
	userAgent string
}

func NewClient(client *http.Client, userAgent string) IClient {
	return &clientStruct{
		client:    client,
		userAgent: userAgent,
	}
}

// build and send request to API with query params, using client/user agent
// response from API returning in xml with encode windows1251, so need to change
func (c *clientStruct) GetJSONCurrency(params MarketAPIparams) (*models.ValCursJSON, error) {
	url := c.urlCreate(params)
	log.Printf("creating get request to the: %v", url)

	// creating browser user-agent, couse cbr dont allow to get data with default agent
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to market error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request denied, status code: %v", resp.StatusCode)
	}
	log.Printf("status code: %v", resp.StatusCode)

	// contentType := resp.Header.Get("Content-Type") // check if not utf-8 txt/xml
	// log.Printf("content-type: %v", contentType)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response reading error: %v", err)
	}

	valuteJSON, err := c.parseRespToJSON(body)
	if err != nil {
		log.Printf("valute: %v", valuteJSON)
		return nil, fmt.Errorf("response parsing error: %v", err)
	}

	// for now if no such valute as params.currency, will response all valutes
	// but after better to response for client message about error currency name
	dataJSON, findCheck := valuteJSON.FindCertainValute(params.ValCode)
	if findCheck {
		valuteJSON.Valute = make([]models.ValuteJSON, 0)
		valuteJSON.Valute = append(valuteJSON.Valute, dataJSON)
	}

	return valuteJSON, nil
}

// taking querry params and market domen, return market API url
// check if params empty, don't add to url, or will be 500 status code
func (c *clientStruct) urlCreate(params MarketAPIparams) string {
	url := fmt.Sprintf("%v", params.Market)
	if params.Date != "" {
		url = url + fmt.Sprintf("?date_req=%v", params.Date)
		if params.ValCode != "" {
			url = url + fmt.Sprintf("&val=%v", params.ValCode)
		}
	} else if params.ValCode != "" {
		url = url + fmt.Sprintf("?val=%v", params.ValCode)
	}
	return url
}

// Create JSON response, by change XML also parsing xml from win1251 to utf8
func (c *clientStruct) parseRespToJSON(respBodyWin1251 []byte) (*models.ValCursJSON, error) {
	log.Print("parser starts")
	var valCursXML *models.ValCursXML
	valCursJSON := &models.ValCursJSON{}

	decoder := charmap.Windows1251.NewDecoder()
	respBodyUTF8, _, err := transform.Bytes(decoder, respBodyWin1251)
	if err != nil {
		return nil, fmt.Errorf("transforming win1251 to UTf-8 err: %v", err) // mb better against nil return &ValCursJSOn{}
	}

	// change encode header, if not, umarshal error
	respBodyUTF8Str := string(respBodyUTF8)
	respBodyUTF8Str = strings.Replace(respBodyUTF8Str, `encoding="windows-1251"`, `encoding="UTF-8"`, 1)
	respBodyUTF8 = []byte(respBodyUTF8Str)

	err = xml.Unmarshal(respBodyUTF8, &valCursXML)
	if err != nil {
		return nil, fmt.Errorf("xml data unmarshal err: %v", err)
	}

	// fill JSON from XML
	valCursJSON.JsonStructFilling(valCursXML)

	return valCursJSON, nil
}
