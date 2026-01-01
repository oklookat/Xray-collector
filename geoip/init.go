package geoip

import (
	_ "github.com/oklookat/xray-collector/geoip/plugin/dbip"
	_ "github.com/oklookat/xray-collector/geoip/plugin/maxmind"
	_ "github.com/oklookat/xray-collector/geoip/plugin/plaintext"
	_ "github.com/oklookat/xray-collector/geoip/plugin/special"
	_ "github.com/oklookat/xray-collector/geoip/plugin/v2ray"
)
