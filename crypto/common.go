package crypto

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func FindWithDelta(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var formatter = &log.TextFormatter{
			FullTimestamp: true,
		}
		log.SetFormatter(formatter)
		f(w, r)
		ch1 := make(chan []*okTicker, 0)
		ch2 := make(chan []*binTicker, 0)
		go getOkexData(ch1)
		go getBinanceData(ch2)
		okexTickers := <-ch1
		binanceTickers := <-ch2
		commonTickers := findCommonTikckers(binanceTickers, okexTickers)
		deltaTickers := findTIckersWithDelta(commonTickers)
		jsonResponse, err := json.Marshal(deltaTickers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
