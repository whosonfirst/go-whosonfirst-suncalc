package providers

import (
	"errors"
	"fmt"		
	"github.com/whosonfirst/go-whosonfirst-geojson"
	"github.com/whosonfirst/go-whosonfirst-lookup"
	"github.com/whosonfirst/go-whosonfirst-utils"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

type WOFFSProvider struct {
	lookup.WOFLookupProvider
	root     string
	is_local bool
	ua       *http.Client
}

func NewWOFFSProvider(root string) (*WOFFSProvider, error) {

	re_file, _ := regexp.Compile(`^file\:\/\/(.*)$`)

	is_local := false
	ua := &http.Client{}

	if re_file.Match([]byte(root)) {

		match := re_file.FindStringSubmatch(root)

		root = "file:///"
		root_trimmed := match[1]

		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir(root_trimmed)))
		ua = &http.Client{Transport: t}
	}

	l := WOFFSProvider{
		root:     root,
		is_local: is_local,
		ua:       ua,
	}

	return &l, nil
}

func (l *WOFFSProvider) Id2AbsPath(id int) string {

	rel_path := utils.Id2RelPath(id)
	return l.root + rel_path
}

// sudo cache me

func (l *WOFFSProvider) GetFeatureById(wofid int) (*geojson.WOFFeature, error) {

	uri := l.Id2AbsPath(wofid)

	fmt.Println(uri)
	
	r, err := l.ua.Get(uri)

	if err != nil && err != io.EOF {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, errors.New("404")
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	feature, err := geojson.UnmarshalFeature(body)

	if err != nil {
		return nil, err
	}

	return feature, nil
}
