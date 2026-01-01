package geosite

import (
	"encoding/json"
	"fmt"

	"github.com/oklookat/xray-collector/geosite/cache"
	"github.com/oklookat/xray-collector/geosite/config"
	"github.com/oklookat/xray-collector/geosite/v2raydata"
)

func Start(cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("config '%s' not loaded", cfgPath)
	}

	cachePath := "./.cache/geosite"
	if len(cfg.CachePath) > 0 {
		cachePath = cfg.CachePath
	}

	cache, err := cache.NewGeosite(cachePath, cfg.MaxCacheSize) // 100MB cache size
	if err != nil {
		return err
	}

	if len(cfg.Input) == 0 {
		return fmt.Errorf("config: missing input")
	}

	var inputResults []*v2raydata.ResolveInputResult

	for _, in := range cfg.Input {
		switch in.Type {

		case config.InputTypeV2RayData:
			var cfgd v2raydata.InputConfig

			if err := json.Unmarshal(in.Args, &cfgd); err != nil {
				return fmt.Errorf("config: invalid args for %s: %w", in.Type, err)
			}

			inputResult, err := v2raydata.ResolveInput(&cfgd, cache)
			if err != nil {
				return fmt.Errorf("v2raydata.ResolveInput: %w", err)
			}
			inputResults = append(inputResults, inputResult)
		}
	}

	inputResult, err := v2raydata.MergeResolveInputResults(inputResults)
	if err != nil {
		return fmt.Errorf("v2raydata.MergeResolveInputResults: %w", err)
	}

	for _, out := range cfg.Output {
		switch out.Type {

		case config.OutputTypeV2RayDat:
			var cfgd v2raydata.OutputConfig

			if err := json.Unmarshal(out.Args, &cfgd); err != nil {
				return fmt.Errorf("config: invalid args for %s: %w", out.Type, err)
			}

			if err := v2raydata.ResolveOutput(&cfgd, inputResult); err != nil {
				return fmt.Errorf("v2raydata.ResolveOutput: %w", err)
			}
		}
	}

	return err
}
