package main

import (
	"errors"
	"flag"
	"io/fs"
	"log/slog"
	"os"

	"github.com/oklookat/xray-collector/geoip"
	"github.com/oklookat/xray-collector/geosite"
)

func main() {
	var (
		buildGeosite bool
		geositeCfg   string

		buildGeoip bool
		geoipCfg   string
		geoipList  bool
	)

	flag.BoolVar(&buildGeosite, "geosite-build", true, "Build geosite")
	flag.BoolVar(&buildGeoip, "geoip-build", true, "Build geoip")

	flag.StringVar(&geositeCfg, "geosite-config", "geosite.json", "Path to geosite config file")
	flag.StringVar(&geoipCfg, "geoip-config", "geoip.json", "Path to geoip config file")

	flag.BoolVar(&geoipList, "geoip-list", false, "List all available input and output formats")

	flag.Parse()

	if buildGeoip || geoipList {
		if _, err := os.Stat(geoipCfg); errors.Is(err, fs.ErrNotExist) {
			slog.Warn("geoip: config not found")
		} else {
			slog.Info("geoip: start")
			chk(geoip.Start(geoipList, geoipCfg), "geoip")
		}
	}

	if buildGeosite {
		if _, err := os.Stat(geositeCfg); errors.Is(err, fs.ErrNotExist) {
			slog.Warn("geosite: config not found")
		} else {
			slog.Info("geosite: start")
			chk(geosite.Start(geositeCfg), "geosite")
		}
	}
}

func chk(err error, issuer string) {
	if err != nil {
		slog.Error(issuer, "error", err.Error())
		os.Exit(1)
	}
	slog.Info(issuer + ": ok")
}
