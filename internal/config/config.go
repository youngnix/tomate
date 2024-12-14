package config

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"github.com/BurntSushi/toml"
)

type Cycle struct {
	Title        string             `toml:"title"`
	Duration     time.Duration      `toml:"duration"`
	Notification *fyne.Notification `toml:"notification"`
}

var Config struct {
	Cycles []*Cycle `toml:"cycles"`
	Notify bool     `toml:"notify"`
	Bell   bool     `toml:"bell"`
}

func init() {
	_, err := toml.DecodeFile("./config.toml", &Config)

	if err != nil {
		log.Fatal(err)
	}
}
