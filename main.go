package main

import (
	"net/http"
	currencyservice "testTaskMod/internal/currencyService"
	externalapi "testTaskMod/internal/externalAPI"
	"testTaskMod/internal/handlers"
)

// to cfg also
const (
	userAg = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) YaBrowser/23.5.3.725 Yowser/2.5 Chrome/113.0.5672.154 Safari/537.36"
	adress = ":8080"
)

func main() {
	client := &http.Client{}
	cbrClient := externalapi.NewClient(client, userAg)
	cbrCurrencyService := currencyservice.NewCurrencyService(cbrClient)
	cbrHandler := handlers.NewHandelr(cbrCurrencyService)
	cbrHandler.StartServer(adress)
}
