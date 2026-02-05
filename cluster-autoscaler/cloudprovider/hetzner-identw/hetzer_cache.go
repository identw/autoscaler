package hetznerIdentw

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/utils/clock"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/hetzner/hcloud-go/hcloud"
)

const (
	defaultCacheTTL   = time.Minute * 3
	jitterCacheMinTTL = 10
	jitterCacheMaxTTL = 60
	serversCacheKey   = "servers-cache"
)

type serverCache struct {
	name    string
	servers []*hcloud.Server
}

type serversCache struct {
	cache.Store
	client *hcloud.Client
}

type serversCacheClock struct {
	clock.Clock

	jitter bool
	sync.RWMutex
}

func newServersCache(hc *hcloud.Client) *serversCache {
	clock := &serversCacheClock{}

	store := cache.NewExpirationStore(func(obj interface{}) (string, error) {
		return obj.(serverCache).name, nil
	}, &cache.TTLPolicy{
		TTL:   defaultCacheTTL,
		Clock: clock,
	})
	return &serversCache{
		Store:  store,
		client: hc,
	}
}

func randomRange(a int, b int) int {
	return rand.Intn(b-a+1) + a
}

func (c *serversCacheClock) Since(t time.Time) time.Duration {
	jitter := time.Duration(randomRange(jitterCacheMinTTL, jitterCacheMaxTTL)) * time.Second
	return time.Since(t.Add(jitter))
}

func (sc *serversCache) getServers() ([]*hcloud.Server, error) {
	obj, exists, err := sc.Store.GetByKey(serversCacheKey)
	if err != nil {
		return nil, err
	}

	if exists {
		return obj.(serverCache).servers, nil
	}

	if !exists {
		listOpts := hcloud.ListOpts{Page: 1, PerPage: 0}
		var serverStatus []hcloud.ServerStatus
		serverStatus = append(serverStatus, hcloud.ServerStatusRunning)
		serverStatus = append(serverStatus, hcloud.ServerStatusInitializing)
		serverStatus = append(serverStatus, hcloud.ServerStatusStarting)
		serverStatus = append(serverStatus, hcloud.ServerStatusDeleting)
		serverListOpts := hcloud.ServerListOpts{ListOpts: listOpts, Name: "", Status: serverStatus}

		servers, err := sc.client.Server.AllWithOpts(context.Background(), serverListOpts)
		if err != nil {
			klog.Errorf("getServers() error get servers. Hetzner API (ServerClient.AllWithOpts: https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud#ServerClient.AllWithOpts), error: %v\n", err)
			return nil, err
		}
		sc.Add(serverCache{
			name:    serversCacheKey,
			servers: servers,
		})
		return servers, nil
	}

	return nil, nil

}
