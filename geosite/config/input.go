package config

import "encoding/json"

type InputType string

const (
	// https://github.com/v2fly/domain-list-community?tab=readme-ov-file#structure-of-data
	InputTypeV2RayData InputType = "v2ray-data"
)

type Input struct {
	Type InputType       `json:"type"`
	Args json.RawMessage `json:"args"`
}
