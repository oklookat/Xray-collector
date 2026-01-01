package config

import "encoding/json"

type OutputType string

const (
	// Compiled v2ray data files.
	OutputTypeV2RayDat OutputType = "v2ray-dat"
)

type OutputAction string

const (
	OutputActionOutput OutputAction = "output"
)

type Output struct {
	Type   OutputType      `json:"type"`
	Action OutputAction    `json:"action"`
	Args   json.RawMessage `json:"args,omitempty"`
}
