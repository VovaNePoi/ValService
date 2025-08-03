package models

import "encoding/xml"

type ValuteXML struct {
	XMLName   xml.Name `xml:"Valute"`
	ID        string   `xml:"ID,attr"`
	NumCode   string   `xml:"NumCode"`
	CharCode  string   `xml:"CharCode"`
	Nominal   string   `xml:"Nominal"`
	Name      string   `xml:"Name"`
	Value     string   `xml:"Value"`
	VunitRate string   `xml:"VunitRate"`
}

// struct for keeping raw data 
type ValCursXML struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Date    string      `xml:"Date,attr"`
	Name    string      `xml:"name,attr"`
	Valute  []ValuteXML `xml:"Valute"`
}
