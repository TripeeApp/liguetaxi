package integration

import (
	"bytes"
	"flag"
	"log"
	"math/rand"
	"net/url"
	"net/http"
	"io/ioutil"
	"os"
	"time"
	"unsafe"

	"bitbucket.org/mobilitee/liguetaxi"
)

// Transport used to log the requests.
type transportLogger struct {
	base http.RoundTripper
}

func (t *transportLogger) RoundTrip(r *http.Request) (*http.Response, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	res, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	log.Printf("Request /%s %s %s --> Response %s %s",
		r.Method, r.URL.String(), string(reqBody), res.Status, string(resBody))

	res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))

	return res, nil
}

const (
	envKeyLiguetaxiToken  = "LIGUETAXI_TOKEN"
	envKeyLiguetaxiHost = "LIGUETAXI_HOST"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax = 63 / letterIdxBits
)

// Ligue Taxi Client
var ligtaxi *liguetaxi.Client

var src = rand.NewSource(time.Now().UnixNano())

var logging = flag.Bool("log", false, "Define if tests should log the requests")

func init() {
	flag.Parse()

	token := os.Getenv(envKeyLiguetaxiToken)
	if token == "" {
		panic("No token defined!!")
	}

	host, _ := url.Parse(os.Getenv(envKeyLiguetaxiHost))

	var hc *http.Client

	if *logging {
		hc = &http.Client{
			Transport: &transportLogger{http.DefaultTransport},
		}
	}

	ligtaxi = liguetaxi.New(host, token, hc)
}

func randString(max int, rangeBytes string) string {
	b := make([]byte, max)

	for i, cache, remain := max-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMax); idx < len(rangeBytes) {
			b[i] = rangeBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
