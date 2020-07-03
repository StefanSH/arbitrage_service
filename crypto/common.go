package crypto

import (
	"encoding/json"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/shopspring/decimal"
)

type commonTicker struct {
	Name        string
	FirstPrice  decimal.Decimal `json:"BinancePrice"`
	SecondPrice decimal.Decimal `json:"OkexPrice"`
	Delta       decimal.Decimal
	Percent     decimal.Decimal
}

func FindWithDelta(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var formatter = &log.TextFormatter{
			FullTimestamp: true,
		}
		log.SetFormatter(formatter)
		f(w, r)

		exchanges := make(map[string]chan []Ticker, 0)
		q := r.URL.Query()

		firstEx := q.Get("first")
		secondEx := q.Get("second")
		if len(firstEx) == 0 || len(secondEx) == 0 {
			http.Error(w, "bad request params", http.StatusBadRequest)
			return
		}
		exchanges[firstEx] = nil
		exchanges[secondEx] = nil

		wg := &sync.WaitGroup{}
		m := &sync.Mutex{}
		wg.Add(2)
		for exchange, _ := range exchanges {
			go func(exchange string) {
				ticker, err := NewTickerFromName(exchange)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				ticksData, err := ticker.GetData()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				ch := make(chan []Ticker, 1)
				tickers := make([]Ticker, 0)
				tickers = append(tickers, typeCast(ticksData)...)
				if len(tickers) == 0 {
					http.Error(w, "invalid data type", http.StatusInternalServerError)
					return
				}
				m.Lock()
				exchanges[exchange] = ch
				m.Unlock()
				ch <- tickers
				wg.Done()
			}(exchange)
		}
		wg.Wait()
		commonTickers, err := findCommonTikckers(<-exchanges[firstEx], <-exchanges[secondEx])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		deltaTickers := findTickersWithDelta(commonTickers)
		jsonResponse, err := json.Marshal(deltaTickers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
