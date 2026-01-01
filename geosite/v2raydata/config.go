package v2raydata

type InputConfig struct {
	// Path to directory or URL.
	//
	// Example:
	//
	// 1. "https://raw.githubusercontent.com/v2fly/domain-list-community/refs/heads/master/data"
	//
	// 2. "./data"
	URI string `json:"uri"`

	// Data filenames to add, files must have parent dir/url of URI.
	//
	// Example ["google", "whatsapp"].
	//
	// Includes will be resolved and added to output.
	Lists []string `json:"lists"`
}

type OutputConfig struct {
	// Output path. Example: "./output/geosite.dat"
	Path string `json:"path"`

	// Lists to add from input. Example: ["google", "whatsapp"]. Includes will be added automatically.
	Lists []string `json:"lists"`
}
