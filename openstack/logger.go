package openstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
)

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default Gophercloud RoundTripper to allow for logging.
type LogRoundTripper struct {
	Logger            *logrus.Logger
	rt                http.RoundTripper
	numReauthAttempts int
}

// newHTTPClient return a custom HTTP client that allows for logging relevant
// information before and after the HTTP request.
func newHTTPClient(l *logrus.Logger) http.Client {
	return http.Client{
		Transport: &LogRoundTripper{
			Logger: l,
			rt:     http.DefaultTransport,
		},
	}
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	defer func() {
		if request.Body != nil {
			request.Body.Close()
		}
	}()

	var err error

	if lrt.Logger.Level == logrus.DebugLevel && request.Body != nil {
		lrt.Logger.Debugln("Logging request body...")
		request.Body, err = lrt.logRequestBody(request.Body, request.Header)
		if err != nil {
			return nil, err
		}
	}

	info, err := json.MarshalIndent(request.Header, "", "  ")
	if err != nil {
		lrt.Logger.Debugf(fmt.Sprintf("Error logging request headers: %s\n", err))
	}
	lrt.Logger.Debugf("Request Headers: %+v\n", string(info))

	lrt.Logger.Infof("Request URL: %s\n", request.URL)

	response, err := lrt.rt.RoundTrip(request)
	if response == nil {
		return nil, err
	}
	response.Body, err = lrt.logResponseBody(response.Body, response.Header)
	if err != nil {
		lrt.Logger.Debugf("Unable to log response body: %s", err)
	}

	if response.StatusCode == http.StatusUnauthorized {
		if lrt.numReauthAttempts == 3 {
			return response, fmt.Errorf("Tried to re-authenticate 3 times with no success.")
		}
		lrt.numReauthAttempts++
	}

	lrt.Logger.Debugf("Response Status: %s\n", response.Status)

	info, err = json.MarshalIndent(response.Header, "", "  ")
	if err != nil {
		lrt.Logger.Debugf(fmt.Sprintf("Error logging response headers: %s\n", err))
	}
	lrt.Logger.Debugf("Response Headers: %+v\n", string(info))

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
		lrt.Logger.Debugf("Response: %s\n", debugInfo)
		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}
	return original, nil
}

func (lrt *LogRoundTripper) logRequestBody(original io.ReadCloser, headers http.Header) (io.ReadCloser, error) {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	contentType := headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		debugInfo := lrt.formatJSON(bs.Bytes())
		lrt.Logger.Debugf("Request Options: %s\n", debugInfo)
	} else {
		lrt.Logger.Debugf("Request Options: %s\n", bs.String())
	}

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil
}

func (lrt *LogRoundTripper) formatJSON(raw []byte) string {
	var data map[string]interface{}

	err := json.Unmarshal(raw, &data)
	if err != nil {
		lrt.Logger.Debugf("Unable to parse JSON: %s\n\n", err)
		return string(raw)
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		lrt.Logger.Debugf("Unable to re-marshal JSON: %s\n", err)
		return string(raw)
	}

	return string(pretty)
}
