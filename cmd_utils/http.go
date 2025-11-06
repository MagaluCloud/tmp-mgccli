package cmdutils

import "net/http"

type Transport struct {
	Headers map[string]string
	Base    http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.Headers {
		req.Header.Add(k, v)
	}
	httpHeaders := http.Header{}
	for reqKey, reqValue := range req.Header {
		if reqKey == "X-Api-Key" && reqValue[0] == "" {
			continue
		}
		httpHeaders.Add(reqKey, reqValue[0])
	}
	req.Header = httpHeaders

	Base := t.Base
	if Base == nil {
		Base = http.DefaultTransport
	}
	return Base.RoundTrip(req)
}
