package openstack

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud/internal/cli/lib"
)

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default Gophercloud RoundTripper to allow for logging.
type LogRoundTripper struct {
	rt                http.RoundTripper
	numReauthAttempts int
}

// newHTTPClient return a custom HTTP client that allows for logging relevant
// information before and after the HTTP request.
func newHTTPClient() http.Client {
	lrt := new(LogRoundTripper)
	lrt.rt = http.DefaultTransport
	lrt.rt.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return http.Client{Transport: lrt}
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	defer func() {
		if request.Body != nil {
			request.Body.Close()
		}
	}()

	var err error

	if request.Body != nil {
		request.Body, err = lrt.logRequestBody(request.Body, request.Header)
		if err != nil {
			return nil, err
		}
	}

	info, err := json.MarshalIndent(request.Header, "", "  ")
	if err != nil {
		lib.Log.Debugf(fmt.Sprintf("Error logging request headers: %s\n", err))
	}
	lib.Log.Debugf("Request Headers: %+v\n", string(info))

	lib.Log.Debugf("Request URL: %s\n", request.URL)

	response, err := lrt.rt.RoundTrip(request)
	if response == nil {
		return nil, err
	}

	var err2 error
	response.Body, err2 = lrt.logResponseBody(response.Body, response.Header)
	if err2 != nil {
		lib.Log.Debugf("Unable to log response body: %s", err2)
	}

	if response.StatusCode == http.StatusUnauthorized {
		if lrt.numReauthAttempts == 3 {
			return response, fmt.Errorf("Tried to re-authenticate 3 times with no success.")
		}
		lrt.numReauthAttempts++
	}

	lib.Log.Debugf("Response Status: %s\n", response.Status)

	var err3 error
	info, err3 = json.MarshalIndent(response.Header, "", "  ")
	if err != nil {
		lib.Log.Debugf(fmt.Sprintf("Error logging response headers: %s\n", err3))
	}
	lib.Log.Debugf("Response Headers: %+v\n", string(info))

	return response, err
}

func (lrt *LogRoundTripper) logResponseBody(original io.ReadCloser, headers http.Header) (io.ReadCloser, error) {
	contentType := headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		defer original.Close()

		var bs bytes.Buffer
		_, err := io.Copy(&bs, original)
		if err != nil {
			return nil, err
		}

		debugInfo := lrt.formatJSON(bs.Bytes())
		lib.Log.Debugf("Response: %s\n", debugInfo)
		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}
	return original, nil
}

func (lrt *LogRoundTripper) logRequestBody(original io.ReadCloser, headers http.Header) (io.ReadCloser, error) {
	contentType := headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		lib.Log.Debugln("Logging request body...")
		defer original.Close()
		var bs bytes.Buffer
		_, err := io.Copy(&bs, original)
		if err != nil {
			return nil, err
		}
		debugInfo := lrt.formatJSON(bs.Bytes())
		lib.Log.Debugf("Request Options: %s\n", debugInfo)
		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}

	return original, nil
}

func (lrt *LogRoundTripper) formatJSON(raw []byte) string {
	var data interface{}

	var m map[string]interface{}
	err := json.Unmarshal(raw, &m)
	switch err {
	case nil:
		data = m
	default:
		var slice []map[string]interface{}
		err = json.Unmarshal(raw, &slice)
		if err != nil {
			lib.Log.Debugf("Unable to parse JSON: %s\n\n", err)
			return string(raw)
		}
		data = slice
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		lib.Log.Debugf("Unable to re-marshal JSON: %s\n", err)
		return string(raw)
	}

	return string(pretty)
}
