package main

import (
	"arbitrage_service/crypto"
	"os"

	b "github.com/claygod/BxogV2"
	log "github.com/sirupsen/logrus"
)

// Main
func main() {
	conf, err := NewTuner("config.toml")
	var formatter = &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	if err != nil {
		panic(err)
	}

	hr := NewHandler(conf)

	m := b.New()
	m.Add("/arbitrage", crypto.FindWithDelta(hr.Default))
	log.Infof("Start server on localhost:%s", conf.Main.Port)
	m.Start(":" + conf.Main.Port)

}
