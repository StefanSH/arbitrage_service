package crypto

import (
	"encoding/json"
	"strings"
)

type OkTicker struct {
	Name  string `json:"instrument_id"`
	Price string `json:"last"`
}

func (*OkTicker) GetData() (interface{}, error) {
	data, err := getTickersData("https://www.okex.com/api/spot/v3/instruments/ticker")
	if err != nil {
		return nil, err
	}
	ticks := make([]*OkTicker, 0)
	err = json.Unmarshal(data, &ticks)
	return ticks, err
}

func (t *OkTicker) GetTickName() string {
	t.Name = strings.ReplaceAll(t.Name, "-", "")
	t.Name = strings.ToUpper(t.Name)
	return t.Name
}

func (t *OkTicker) GetPrice() string {
	return t.Price
}

