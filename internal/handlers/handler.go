package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	currencyservice "testTaskMod/internal/currencyService"
	externalapi "testTaskMod/internal/externalAPI"
)

// cbr API URL. TODO: config will be better to change API without complation
const (
	market = "https://cbr.ru/scripts/XML_daily.asp"
)

type Handler struct {
	currService currencyservice.CurrencyService
}

func NewHandelr(curService currencyservice.CurrencyService) *Handler {
	return &Handler{
		currService: curService,
	}
}

// getting query params for market API request
// sending response to requesting service
func (h *Handler) HandleCurrencyService(w http.ResponseWriter, r *http.Request) {
	// get query params and saving in structure for using in url builder
	urlValues := r.URL.Query()
	queryParams := &externalapi.MarketAPIparams{}
	queryParams.SetParams(urlValues, market)
	log.Printf("handleCurrService, params: %v, %v", queryParams.ValCode, queryParams.Date)

	jsonData, err := h.currService.GetCurrency(*queryParams)
	if err != nil {
		log.Printf("geting currency error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	err = enc.Encode(jsonData)
	if err != nil {
		log.Printf("encoding error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// listen for get requests, starting server on different adresses, default = localhost:8080
func (h *Handler) StartServer(adress string) error {
	http.HandleFunc("/currency/", h.HandleCurrencyService)
	log.Printf("Server start on http://localhost:8080/currency/")
	return http.ListenAndServe(adress, nil)
}
