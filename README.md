# arbitrage_service
Microservice based on https://github.com/claygod/microservice</br>/
Сервис получает данные через публичные api двух популярных криптовалютных бирж(Binance и OKEx), находит общие для обеих бирж торговые пары и отдает в виде JSON пары с разницей в цене между биржами > 1%

## Dependencies
github.com/Sirupsen/logrus
github.com/shopspring/decimal
github.com/StefanSH/arbitrage_service/crypto
github.com/claygod/BxogV2

## Build and Start

```golang
go get github.com/Sirupsen/logrus
go get github.com/shopspring/decimal
go get github.com/StefanSH/arbitrage_service/crypto
go get github.com/claygod/BxogV2
```

```golang
go build
./arbitrage_service
```
   
Main functions in /crypto/core.go </br>
default url - localhost:80/arbitrage

