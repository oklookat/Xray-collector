package v2raydata

import (
	"net/url"
	"testing"

	"github.com/oklookat/xray-collector/geosite/cache"
)

func TestNew(t *testing.T) {
	const testCache = "./test/cache"
	cach, err := cache.NewGeosite(testCache, 100000000)
	if err != nil {
		t.Fatal(err)
	}
	// defer func ()  {
	// 	os.RemoveAll(testCache)
	// }()

	rootURL, err := url.Parse("https://raw.githubusercontent.com/v2fly/domain-list-community/refs/heads/master/data")
	if err != nil {
		t.Fatal(err)
	}

	remList := newInputRemoteList(cach, rootURL)
	if _, err := remList.Resolve("category-ru"); err != nil {
		t.Fatal(err)
	}
}
