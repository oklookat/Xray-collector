# Xray-collector

Builds `geoip` and `geosite` files for [Xray-core](https://github.com/XTLS/Xray-core).

`./xray-collector --help`

## geoip

Same config as [v2fly/geoip](https://github.com/v2fly/geoip).

`./xray-collector -geoip-config geoip.json -geosite-build=false`

## geosite

For build, uses modified version of [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community). Without plaintext format export support.

`./xray-collector -geosite-config geosite.json -geoip-build=false`

Config file example:

```jsonc
{
    // All downloaded geosite files and ETags will be here.
    "cachePath": "./build/geosite-cache",
    // In bytes. If max size reached, oldest files will be deleted.
    // 0 == unlimited size.
    "maxCacheSize": 100000000,
    "input": [
        {
            "type": "v2ray-data",
            "args": {
                // Directory or URL with lists.
                "uri": "./build/geosite-local",
                // Lists to get.
                "lists": [
                    // List in format of https://github.com/v2fly/domain-list-community?tab=readme-ov-file#structure-of-data
                    // If list containts includes, it must be exists in "uri". If not, will be error.
                    "censor-tracker"
                ]
            }
        },
        {
            "type": "v2ray-data",
            "args": {
                // In case of URL: URL to DIRECTORY with this lists. 
                // Example: https://raw.githubusercontent.com/v2fly/domain-list-community/refs/heads/master/data
                // + "google"
                // =: https://raw.githubusercontent.com/v2fly/domain-list-community/refs/heads/master/data/google
                "uri": "https://raw.githubusercontent.com/v2fly/domain-list-community/refs/heads/master/data", 
                // Lists to get. Includes in this lists will be downloaded and cached too.
                "lists": [
                    "google",
                    "whatsapp",
                    "category-ru"
                ]
            }
        }
    ],
    "output": [
        // In output you get all available lists:
        // ["censor-tracker", "google", "whatsapp", "category-ru"]
        {
            "type": "v2ray-dat",
            "action": "output",
            "args": {
                "path": "./build/output/geosite-client.dat",
                // You can select lists to build, includes will builded too.
                "lists": [
                    "category-ru"
                ]
            }
        },
        {
            "type": "v2ray-dat",
            "action": "output",
            "args": {
                "path": "./build/output/geosite-server.dat",
                "lists": [
                    "censor-tracker",
                    "google",
                    "whatsapp"
                ]
            }
        }
    ]
}
```
