package models

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type ValuteJSON struct {
	ID        string  `json:"ID"`
	NumCode   int     `json:"NumCode"`
	CharCode  string  `json:"CharCode"`
	Nominal   int     `json:"Nominal"`
	Name      string  `json:"Name"`
	Value     float64 `json:"Value"`
	VunitRate float64 `json:"VunitRate"`
}

// struct for keeping ready data
type ValCursJSON struct {
	Date   string       `json:"Date"`
	Name   string       `json:"Name"`
	Valute []ValuteJSON `json:"Valute"`
}

func (v *ValCursJSON) JsonStructFilling(xmlCurs *ValCursXML) (*ValCursJSON, error) {
	tmpDate, err := time.Parse("02.01.2006", xmlCurs.Date)
	if err != nil {
		log.Printf("time parsing error: %v", err)
		v.Date = ""
	} else {
		v.Date = tmpDate.Format(time.DateOnly)
	}
	v.Name = xmlCurs.Name
	v.Valute = make([]ValuteJSON, len(xmlCurs.Valute))
	for i, val := range xmlCurs.Valute {
		numCode, err := strconv.Atoi(val.NumCode)
		if err != nil {
			log.Printf("numCode parsing error: %v", err)
			numCode = 0
		}
		nominal, err := strconv.Atoi(val.Nominal)
		if err != nil {
			log.Printf("nominal parsing error: %v", err)
			nominal = 0
		}
		value, err := strconv.ParseFloat(strings.ReplaceAll(val.Value, ",", "."), 64) // "," на "."
		if err != nil {
			log.Printf("value parsing error: %v", err)
			value = 0.0
		}
		vunitRate, err := strconv.ParseFloat(strings.ReplaceAll(val.VunitRate, ",", "."), 64) // "," на "."
		if err != nil {
			log.Printf("vunitRate parsing error: %v", err)
			vunitRate = 0.0
		}
		// Заполнение структуры ValuteJSON
		v.Valute[i] = ValuteJSON{
			ID:        val.ID,
			NumCode:   numCode,
			CharCode:  val.CharCode,
			Nominal:   nominal,
			Name:      val.Name,
			Value:     value,
			VunitRate: vunitRate,
		}
	}

	return v, nil
}

// searching valute by charCode example(USDT, AUD)
func (v *ValCursJSON) FindCertainValute(targetVal string) (ValuteJSON, bool) {
	for _, val := range v.Valute {
		if val.CharCode == targetVal {
			return val, true
		}
	}
	return ValuteJSON{}, false
}
