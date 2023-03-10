package round

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"net/url"
	"sync/atomic"
)

// ErrServersNotExists is the error that servers dose not exists
var ErrServersNotExists = errors.New("servers dose not exist")

type URLProxy struct {
	url     *url.URL
	Service *api.CatalogService
}

// RoundRobin is an interface for representing round-robin balancing.
type RoundRobin interface {
	Next() *URLProxy
	GetUrls() []*URLProxy
}

type roundrobin struct {
	urls []*URLProxy
	next uint32
}

// New returns RoundRobin implementation(*roundrobin).
func New(urls []*URLProxy) (RoundRobin, error) {
	if len(urls) == 0 {
		return nil, ErrServersNotExists
	}

	return &roundrobin{
		urls: urls,
	}, nil
}

func BuildRound(service []*api.CatalogService) (RoundRobin, error) {
	arr := make([]*URLProxy, len(service))
	for i := range service {
		arr[i] = &URLProxy{url: &url.URL{Host: service[i].ServiceAddress}, Service: service[i]}
	}
	return New(arr)
}

// Next returns next address
func (r *roundrobin) Next() *URLProxy {
	n := atomic.AddUint32(&r.next, 1)
	return r.urls[(int(n)-1)%len(r.urls)]
}

func (r *roundrobin) GetUrls() []*URLProxy {
	return r.urls
}
