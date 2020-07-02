# arbitrage_service
Based on https://github.com/claygod/microservice.
Сервис получает данные через публичные api двух популярных криптовалютных бирж(Binance и OKEx), находит общие для обеих бирж торговые пары и отдает в виде JSON пары с разницей в цене между биржами

## Run

```golang
git clone github.com/StefanSH/arbitrage_service
go get
go run ./
```

```GET
go build
./arbitrage_service
```
   
Main functions in /crypto/core.go </br>
test url - localhost:80/arbitrage?first=binance&second=okex

