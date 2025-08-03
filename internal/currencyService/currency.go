package currencyservice

import (
	"fmt"
	externalapi "testTaskMod/internal/externalAPI"
	"testTaskMod/internal/models"
)

// abstract data getting from external API
type CurrencyService interface {
	GetCurrency(params externalapi.MarketAPIparams) (*models.ValCursJSON, error)
}

type currency struct {
	client externalapi.IClient
}

func NewCurrencyService(client externalapi.IClient) CurrencyService {
	return &currency{client: client}
}

func (c *currency) GetCurrency(params externalapi.MarketAPIparams) (*models.ValCursJSON, error) {
	dataJSON, err := c.client.GetJSONCurrency(params)
	if err != nil {
		return nil, fmt.Errorf("getting JSON currency error: %v", err)
	}

	return dataJSON, nil
}
