package externalapi

import (
	"log"
	"net/url"
	"time"
)

// struct for query parameters
type MarketAPIparams struct {
	Market  string
	Date    string
	ValCode string
}

func (m *MarketAPIparams) SetParams(urlValues url.Values, market string) {
	m.Market = market // need to remake for more flexibility

	m.ValCode = urlValues.Get("val_code")

	// useless, cbr API automatic set today if parametr empty
	if urlValues.Get("val_code") == "" {
		m.Date = time.Now().Format("02.01.2006")
		log.Printf("date: %v", m.Date)
	} else {
		m.Date = urlValues.Get("date")
	}
}
