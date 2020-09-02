package test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	testUtil "github.com/x893675/gocron/test/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

var (
	serviceAddress        = ""
	servicePort           = "8080"
	serviceEndpointPrefix = ""
)

func init() {
	flag.StringVar(&serviceAddress, "addr", "localhost", "service address")
	flag.Parse()
	if len(flag.Args()) > 0 {
		for _, iFlag := range flag.Args() {
			keyVal := strings.Split(iFlag, "=")
			if keyVal[0] == "addr" {
				serviceAddress = keyVal[1]
			}
		}
	}
	serviceEndpointPrefix = fmt.Sprintf("http://%s:%s", serviceAddress, servicePort)
}

func requestTest(t *testing.T, requestURL, httpMethod string, postBody json.RawMessage, header map[string]string, skipTLSCheck, disableKeepAlive bool) ([]byte, int) {
	body, statusCode, err := testUtil.CommonRequest(requestURL,
		httpMethod, postBody, header,
		skipTLSCheck, disableKeepAlive, 100*time.Second)
	if err != nil {
		t.Errorf("Failed To Request url %s with data %s", requestURL, postBody)
		return nil, -1
	}
	return body, statusCode
}

func requestReturnWithHeader(requestUrl, httpMethod string, postBody json.RawMessage, header map[string]string, skipTlsCheck, disableKeepAlive bool, timeout time.Duration) ([]byte, int, http.Header, error) {
	var req *http.Request
	var reqErr error

	req, reqErr = http.NewRequest(httpMethod, requestUrl, bytes.NewReader(postBody))

	if reqErr != nil {
		return []byte{}, http.StatusInternalServerError, nil, reqErr
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	for key, val := range header {
		req.Header.Add(key, val)
	}
	client := &http.Client{}
	client.Timeout = timeout
	if skipTlsCheck {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: disableKeepAlive}
	} else {
		client.Transport = &http.Transport{DisableKeepAlives: disableKeepAlive}
	}
	resp, respErr := client.Do(req)
	if respErr != nil {
		return []byte{}, http.StatusInternalServerError, nil, respErr
	}
	defer resp.Body.Close()
	body, readBodyErr := ioutil.ReadAll(resp.Body)
	if readBodyErr != nil {
		return []byte{}, http.StatusInternalServerError, nil, readBodyErr
	}
	respHeader := resp.Header
	return body, resp.StatusCode, respHeader, nil
}

func getRequestURL(path string) string {
	return fmt.Sprintf("%s%s", serviceEndpointPrefix, path)
}
