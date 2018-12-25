package base

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostProfile(ctx context.Context, p Profile) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostProfile", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostProfile(ctx, p)
}

func (mw loggingMiddleware) GetProfile(ctx context.Context, id string) (p Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetProfile", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetProfile(ctx, id)
}

func (mw loggingMiddleware) PutProfile(ctx context.Context, id string, p Profile) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutProfile", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutProfile(ctx, id, p)
}

func (mw loggingMiddleware) DeleteProfile(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteProfile", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteProfile(ctx, id)
}

func (mw loggingMiddleware) GetCards() (p PokemonResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetCards", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetCards()
}

func (mw loggingMiddleware) FetchData() (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "FetchData", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.FetchData()
}

func InstrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	countResult metrics.Histogram,
) Middleware {
	return func(next Service) Service {
		return instrmw{requestCount, requestLatency, countResult, next}
	}
}

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	Service
}

func (mw instrmw) PostProfile(ctx context.Context, p Profile) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "postprofile", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Service.PostProfile(ctx, p)
	return
}

func (mw instrmw) GetProfile(ctx context.Context, id string) (p Profile, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getprofile", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	p, err = mw.Service.GetProfile(ctx, id)
	return
}

func (mw instrmw) PutProfile(ctx context.Context, id string, p Profile) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "putprofile", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Service.PutProfile(ctx, id, p)
	return
}

func (mw instrmw) DeleteProfile(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "deleteprofile", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Service.DeleteProfile(ctx, id)
	return
}

func (mw instrmw) GetCards() (p PokemonResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetCards", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.GetCards()
}
