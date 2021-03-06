package main

import (
	"arood/base"
	"arood/client"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kithttp "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	var (
		httpAddr   = flag.String("http.addr", ":8080", "HTTP listen address")
		pokemonURL = flag.String("pokemon.url", "https://api.pokemontcg.io/v1/cards", "Outbound Pokemon address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	errs := make(chan error)

	clientURL, err := url.Parse(*pokemonURL)
	if err != nil {
		errs <- err
	}

	options := []kithttp.ClientOption{}
	pokemonEndpoint := kithttp.NewClient(http.MethodGet, clientURL, client.EncodePokemonRequest, client.DecodePokemonResponse, options...).Endpoint()

	var s base.Service
	{
		s = base.NewInmemService(pokemonEndpoint)
		s = base.LoggingMiddleware(logger)(s)
		s = base.InstrumentingMiddleware(requestCount, requestLatency, countResult)(s)
	}

	var h http.Handler
	{
		h = base.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		err := s.FetchData()
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
