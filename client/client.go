// Package client provides a profilesvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"arood/base"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func DecodePokemonResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response base.PokemonResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func EncodePokemonRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}
