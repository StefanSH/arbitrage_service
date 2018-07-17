package crypto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/shopspring/decimal"
)

func getTickersData(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Errorf("Error:%s %s doesn't response", err, url)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error:%s Failed to read response", err)
	}
	return contents
}

func findCommonTikckers(binanceTickers []*binTicker, okexTickers []*okTicker) []*commonTicker {
	var common []*commonTicker
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sort.Slice(binanceTickers, func(i, j int) bool {
			return binanceTickers[i].Name < binanceTickers[j].Name
		})
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		sort.Slice(okexTickers, func(i, j int) bool {
			return okexTickers[i].Name < okexTickers[j].Name
		})
		wg.Done()
	}()
	wg.Wait()
	for i := range binanceTickers {
		for j := range okexTickers {
			if binanceTickers[i].Name == okexTickers[j].Name {
				decBinPrice, _ := decimal.NewFromString(binanceTickers[i].Price)
				decOkPrice, _ := decimal.NewFromString(okexTickers[j].Price)
				common = append(common, &commonTicker{
					Name:     binanceTickers[i].Name,
					BinPrice: decBinPrice,
					OkPrice:  decOkPrice})
			}
		}
	}
	return common
}

func getOkexData(out chan []*okTicker) {
	var temp *OkexData
	data := getTickersData("https://www.okex.com/api/v1/tickers.do")
	err := json.Unmarshal(data, &temp)
	if err != nil {
		log.Errorf("Error:%s Failed to Unmarshal json", err)
	}
	for _, ticker := range temp.Data {
		ticker.formatSymbol()
	}
	okexTickers := temp.Data
	out <- okexTickers
}

func getBinanceData(out chan []*binTicker) {
	var temp []*binTicker
	data := getTickersData("https://api.binance.com/api/v3/ticker/price")
	err := json.Unmarshal(data, &temp)
	if err != nil {
		fmt.Println("Failed to unmarshal binance", err)
	}
	out <- temp
}

func findTIckersWithDelta(commonTickers []*commonTicker) []*commonTicker {
	var deltaTickers []*commonTicker
	for _, ticker := range commonTickers {
		if ticker.BinPrice.GreaterThan(ticker.OkPrice) {
			delta := decimal.Sum(ticker.BinPrice.Div(ticker.OkPrice), decimal.New(-1, 0)).Mul(decimal.New(100, 0))
			if delta.GreaterThan(decimal.New(1, 0)) {
				ticker.Delta = delta.Round(4)
				deltaTickers = append(deltaTickers, ticker)
			}
		} else {
			delta := decimal.Sum(ticker.OkPrice.Div(ticker.BinPrice), decimal.New(-1, 0)).Mul(decimal.New(100, 0))
			if delta.GreaterThan(decimal.New(1, 0)) {
				ticker.Delta = delta.Round(4)
				deltaTickers = append(deltaTickers, ticker)
			}
		}
	}
	return deltaTickers
}
