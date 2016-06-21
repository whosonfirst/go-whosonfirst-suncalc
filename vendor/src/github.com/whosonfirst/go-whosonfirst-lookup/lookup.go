package lookup

import (
	"github.com/whosonfirst/go-whosonfirst-geojson"
)

type WOFLookupProvider interface {
	GetFeatureById(wofid int) (*geojson.WOFFeature, error)
}
