package crypto

import "encoding/json"

type BinTicker struct {
	Name  string `json:"symbol"`
	Price string `json:"price"`
}

func (*BinTicker) GetData() (interface{}, error) {
	data, err := getTickersData("https://api.binance.com/api/v3/ticker/price")
	if err != nil {
		return nil, err
	}
	ticks := make([]*BinTicker, 0)
	err = json.Unmarshal(data, &ticks)
	return ticks, err
}


func (t BinTicker) GetTickName() string {
	return t.Name
}

func (t BinTicker) GetPrice() string {
	return t.Price
}
