package main

import (
	"context"
	"io"
	"net/http"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = registerer("krakend-no-redirects")

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(_ context.Context, extra map[string]interface{}) (http.Handler, error) {
	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		client := http.DefaultClient
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy headers, status codes, and body from the backend to the response writer
		for k, hs := range resp.Header {
			for _, h := range hs {
				w.Header().Add(k, h)
			}
		}
		w.WriteHeader(resp.StatusCode)
		if resp.Body == nil {
			return
		}
		io.Copy(w, resp.Body)
		resp.Body.Close()

	}), nil
}

func main() {}
