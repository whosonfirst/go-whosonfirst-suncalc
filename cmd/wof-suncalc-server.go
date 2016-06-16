package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-httpony/cors"
	"github.com/whosonfirst/go-httpony/tls"
	"github.com/whosonfirst/go-whosonfirst-geojson"
	"github.com/whosonfirst/go-whosonfirst-utils"
	"github.com/whosonfirst/suncalc-go"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

// sudo put me in a package somewhere... specifically make me an interface
// so lookups can be remote, local, in a DB whatever (20160616/thisisaaronland)

type WOFLookup struct {
	Root     string
	is_local bool
	ua       *http.Client
}

func NewWOFLookup(root string) (*WOFLookup, error) {

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

	l := WOFLookup{
		Root:     root,
		is_local: is_local,
		ua:       ua,
	}

	return &l, nil
}

func (l *WOFLookup) Id2AbsPath(id int) string {

	rel_path := utils.Id2RelPath(id)
	return l.Root + rel_path
}

// sudo cache me

func (l *WOFLookup) GetFeatureById(wofid int) (*geojson.WOFFeature, error) {

	uri := l.Id2AbsPath(wofid)

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

// sudo cache me

func (l *WOFLookup) GetLatLonById(wofid int) (float64, float64, error) {

	var lat float64
	var lon float64

	feature, err := l.GetFeatureById(wofid)

	if err != nil {
		return lat, lon, err
	}

	fbody := feature.Body()

	var ok bool

	lat, ok = fbody.Path("properties.geom:latitude").Data().(float64)

	if !ok {
		return lat, lon, errors.New("Missing latitude")
	}

	lon, ok = fbody.Path("properties.geom:longitude").Data().(float64)

	if !ok {
		return lat, lon, errors.New("Missing longitude")
	}

	return lat, lon, nil
}

func main() {

	var host = flag.String("host", "localhost", "Hostname to listen on")
	var port = flag.Int("port", 8080, "Port to listen on")
	var cors_enable = flag.Bool("cors", false, "...")
	var cors_allow = flag.String("cors-allow", "*", "...")
	var wof_root = flag.String("wof-root", "https://whosonfirst.mapzen.com/data/", "...")
	var tls_enable = flag.Bool("tls", false, "Serve requests over TLS") // because CA warnings in browsers...
	var tls_cert = flag.String("tls-cert", "", "Path to an existing TLS certificate. If absent a self-signed certificate will be generated.")
	var tls_key = flag.String("tls-key", "", "Path to an existing TLS key. If absent a self-signed key will be generated.")

	flag.Parse()

	lookup, err := NewWOFLookup(*wof_root)

	if err != nil {
		panic(err)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	request_handler := func() http.Handler {

		re_ymd, _ := regexp.Compile(`^(\d{4})-?(\d{2})-?(\d{2})$`)

		fn := func(rsp http.ResponseWriter, req *http.Request) {

			var t time.Time
			var lat float64
			var lon float64

			query := req.URL.Query()

			str_wofid := query.Get("wofid")
			str_ymd := query.Get("ymd")

			str_lat := query.Get("latitude")
			str_lon := query.Get("longitude")

			if str_wofid != "" {

				wof_id, err := strconv.Atoi(str_wofid)

				if err != nil {
					http.Error(rsp, "400 Bad Gateway", http.StatusBadRequest)
					return
				}

				w_lat, w_lon, err := lookup.GetLatLonById(wof_id)

				if err != nil {
					http.Error(rsp, err.Error(), http.StatusInternalServerError)
					return
				}

				lat = w_lat
				lon = w_lon

			} else {

				if str_lat == "" {
					http.Error(rsp, "Missing latitude parameter", http.StatusBadRequest)
					return
				}

				if str_lon == "" {
					http.Error(rsp, "Missing longitude parameter", http.StatusBadRequest)
					return
				}

				var err error

				lat, err = strconv.ParseFloat(str_lat, 64)

				if err != nil {
					http.Error(rsp, "Invalid latitude parameter", http.StatusBadRequest)
					return
				}

				lon, err = strconv.ParseFloat(str_lon, 64)

				if err != nil {
					http.Error(rsp, "Invalid longitude parameter", http.StatusBadRequest)
					return
				}
			}

			if lat > 90.0 || lat < -90.0 {
				http.Error(rsp, "E_IMPOSSIBLE_LATITUDE", http.StatusBadRequest)
				return
			}

			if lon > 180.0 || lon < -180.0 {
				http.Error(rsp, "E_IMPOSSIBLE_LONGITUDE", http.StatusBadRequest)
				return
			}

			t_string := ""
			t_format := "2006-01-02T15:04:05Z-0700"

			if str_ymd != "" {

				if !re_ymd.MatchString(str_ymd) {
					http.Error(rsp, "E_IMPOSSIBLE_YEARMONTHDAY", http.StatusBadRequest)
				}

				m := re_ymd.FindStringSubmatch(str_ymd)

				t_string = fmt.Sprintf("%s-%s-%sT00:00:00Z-0000", m[1], m[2], m[3])
			}

			if t_string == "" {
				t = time.Now()
			} else {

				tp, err := time.Parse(t_format, t_string)

				if err != nil {
					http.Error(rsp, "E_IMPOSSIBLE_TIMESTRING", http.StatusBadRequest)
					return
				}

				t = tp
			}

			times := suncalc.SunTimes(t, lat, lon)

			js, err := json.Marshal(times)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			rsp.Header().Set("Content-Type", "application/json")
			rsp.Write(js)
		}

		return http.HandlerFunc(fn)
	}

	handler := cors.EnsureCORSHandler(request_handler(), *cors_enable, *cors_allow)

	if *tls_enable {

		var cert string
		var key string

		if *tls_cert == "" && *tls_key == "" {

			root, err := tls.EnsureTLSRoot()

			if err != nil {
				panic(err)
			}

			cert, key, err = tls.GenerateTLSCert(*host, root)

			if err != nil {
				panic(err)
			}

		} else {
			cert = *tls_cert
			key = *tls_key
		}

		fmt.Printf("start and listen for requests at https://%s\n", endpoint)
		err = http.ListenAndServeTLS(endpoint, cert, key, handler)

	} else {

		fmt.Printf("start and listen for requests at http://%s\n", endpoint)
		err = http.ListenAndServe(endpoint, handler)
	}

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
