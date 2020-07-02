package crypto

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"

	"github.com/shopspring/decimal"
)

type Ticker interface {
	GetTickName() string
	GetPrice() string
	GetData() (interface{}, error)
}

func getTickersData(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error:%s %s doesn't response", err, url)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error:%s Failed to read response", err)
	}
	return contents, nil
}

func findCommonTikckers(firstTickers, secondTickers []*Ticker) ([]*commonTicker, error) {
	if len(firstTickers) == 0 || len(secondTickers) == 0 {
		return nil, fmt.Errorf("one or two exhages data is nil")
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		sort.Slice(firstTickers, func(i, j int) bool {
			ticksI := *firstTickers[i]
			ticksJ := *firstTickers[j]
			return ticksI.GetTickName() < ticksJ.GetTickName()
		})
		wg.Done()
	}()
	go func() {
		sort.Slice(secondTickers, func(i, j int) bool {
			ticksI := *secondTickers[i]
			ticksJ := *secondTickers[j]
			return ticksI.GetTickName() < ticksJ.GetTickName()
		})
		wg.Done()
	}()
	wg.Wait()
	return compare(firstTickers, secondTickers), nil
}

func compare(ticks1, ticks2 []*Ticker) []*commonTicker {
	common := make([]*commonTicker, 0)
	if len(ticks1) > len(ticks2) {
		for i := range ticks1 {
			t1 := *ticks1[i]
			for j := range ticks2 {
				t2 := *ticks2[j]
				if t1.GetTickName() == t2.GetTickName() {
					firstPrice, _ := decimal.NewFromString(t1.GetPrice())
					secondPrice, _ := decimal.NewFromString(t2.GetPrice())
					common = append(common, &commonTicker{
						Name:     t1.GetTickName(),
						FirstPrice: firstPrice,
						SecondPrice:  secondPrice})
				}
			}
		}
	} else {
		for i := range ticks2 {
			t2 := *ticks2[i]
			for j := range ticks1 {
				t1 := *ticks1[j]
				if t1.GetTickName() == t2.GetTickName() {
					firstPrice, _ := decimal.NewFromString(t1.GetPrice())
					secondPrice, _ := decimal.NewFromString(t2.GetPrice())
					common = append(common, &commonTicker{
						Name:     t1.GetTickName(),
						FirstPrice: firstPrice,
						SecondPrice:  secondPrice})
				}
			}
		}
	}

	return common
}

func NewTickerFromName(name string) (Ticker, error) {
	switch name {
	case "binance":
		return &BinTicker{} , nil
	case "okex":
		return &OkTicker{} , nil
	default:
		return nil, fmt.Errorf("exhange %s - not found", name)
	}
}

func typeCast(i interface{}) []*Ticker {
	switch i.(type) {
	case []*OkTicker:
		result := make([]*Ticker, 0)
		for _, tick := range i.([]*OkTicker) {
			ticker := Ticker(tick)
			result = append(result, &ticker)
		}
		return result
	case []*BinTicker:
		result := make([]*Ticker, 0)
		for _, tick := range i.([]*BinTicker) {
			ticker := Ticker(tick)
			result = append(result, &ticker)
		}
		return result
	default:
		return nil
	}

}

func findTickersWithDelta(commonTickers []*commonTicker) []*commonTicker {
	var deltaTickers []*commonTicker
	for _, ticker := range commonTickers {
		if ticker.FirstPrice.GreaterThan(ticker.SecondPrice) {
			delta := decimal.Sum(ticker.FirstPrice, decimal.New(-1, 0).Mul(ticker.SecondPrice))
			if delta.GreaterThan(decimal.New(1, 0)) {
				ticker.Delta = delta
				ticker.Percent = delta.Div(ticker.FirstPrice.Div(decimal.New(100, 0))).Round(4)
				deltaTickers = append(deltaTickers, ticker)
			}
		} else {
			delta := decimal.Sum(ticker.SecondPrice, decimal.New(-1,0).Mul(ticker.FirstPrice))
			ticker.Delta = delta
			ticker.Percent = delta.Div(ticker.SecondPrice.Div(decimal.New(100, 0))).Round(4)
			deltaTickers = append(deltaTickers, ticker)
			}
		}
	return deltaTickers
}
