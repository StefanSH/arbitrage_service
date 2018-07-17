package crypto

import (
	"strings"

	"github.com/shopspring/decimal"
)

type binTicker struct {
	Name  string `json:"symbol"`
	Price string `json:"price"`
}

type okTicker struct {
	Name  string `json:"symbol"`
	Price string `json:"buy"`
}

type commonTicker struct {
	Name     string
	BinPrice decimal.Decimal `json:"BinancePrice"`
	OkPrice  decimal.Decimal `json:"OkexPrice"`
	Delta    decimal.Decimal
}

func (t *okTicker) formatSymbol() {
	t.Name = strings.Replace(t.Name, "_", "", -1)
	t.Name = strings.ToUpper(t.Name)
}

type OkexData struct {
	Data []*okTicker `json:"tickers"`
}
