// Package client provides a profilesvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"arood/base"
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

// New returns a service that's load-balanced over instances of profilesvc found
// in the provided Consul server.
func New(consulAddr string, logger log.Logger) (base.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	var (
		consulService = "profilesvc"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints base.Endpoints
	)
	{
		factory := factoryFor(base.MakePostProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfileEndpoint = retry
	}
	{
		factory := factoryFor(base.MakeGetProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetProfileEndpoint = retry
	}
	{
		factory := factoryFor(base.MakePutProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PutProfileEndpoint = retry
	}
	{
		factory := factoryFor(base.MakeDeleteProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.DeleteProfileEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(base.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := base.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
