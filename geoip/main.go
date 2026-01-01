package geoip

import (
	"flag"
	"fmt"

	"github.com/oklookat/xray-collector/geoip/lib"
)

func Start(list bool, configFile string) error {
	flag.Parse()

	if list {
		lib.ListInputConverter()
		fmt.Println()
		lib.ListOutputConverter()
		return nil
	}

	instance, err := lib.NewInstance()
	if err != nil {
		return err
	}

	if err := instance.InitConfig(configFile); err != nil {
		return err
	}

	if err := instance.Run(); err != nil {
		return err
	}

	return err
}
