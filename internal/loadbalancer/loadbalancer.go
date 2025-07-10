package loadbalancer

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// LoadBalancer defines the interface for load balancing strategies
type LoadBalancer interface {
	NextEndpoint() string
	AddEndpoint(endpoint string)
	RemoveEndpoint(endpoint string)
	GetHealthyEndpoints() []string
}

// Strategy types
const (
	RoundRobin     = "round-robin"
	Random         = "random"
	WeightedRandom = "weighted-random"
	LeastConn      = "least-conn"
)

// New creates a new load balancer with the specified strategy
func New(endpoints []string, strategy string) LoadBalancer {
	switch strategy {
	case Random:
		return NewRandomLB(endpoints)
	case WeightedRandom:
		return NewWeightedRandomLB(endpoints)
	case LeastConn:
		return NewLeastConnLB(endpoints)
	default:
		return NewRoundRobinLB(endpoints)
	}
}

// Round Robin Load Balancer
type roundRobinLB struct {
	endpoints []string
	current   uint64
	mu        sync.RWMutex
}

func NewRoundRobinLB(endpoints []string) LoadBalancer {
	return &roundRobinLB{
		endpoints: endpoints,
	}
}

func (rr *roundRobinLB) NextEndpoint() string {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	if len(rr.endpoints) == 0 {
		return ""
	}
	
	index := atomic.AddUint64(&rr.current, 1) % uint64(len(rr.endpoints))
	return rr.endpoints[index]
}

func (rr *roundRobinLB) AddEndpoint(endpoint string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.endpoints = append(rr.endpoints, endpoint)
}

func (rr *roundRobinLB) RemoveEndpoint(endpoint string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	for i, ep := range rr.endpoints {
		if ep == endpoint {
			rr.endpoints = append(rr.endpoints[:i], rr.endpoints[i+1:]...)
			break
		}
	}
}

func (rr *roundRobinLB) GetHealthyEndpoints() []string {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	result := make([]string, len(rr.endpoints))
	copy(result, rr.endpoints)
	return result
}

// Random Load Balancer
type randomLB struct {
	endpoints []string
	rand      *rand.Rand
	mu        sync.RWMutex
}

func NewRandomLB(endpoints []string) LoadBalancer {
	return &randomLB{
		endpoints: endpoints,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *randomLB) NextEndpoint() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.endpoints) == 0 {
		return ""
	}
	
	index := r.rand.Intn(len(r.endpoints))
	return r.endpoints[index]
}

func (r *randomLB) AddEndpoint(endpoint string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.endpoints = append(r.endpoints, endpoint)
}

func (r *randomLB) RemoveEndpoint(endpoint string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	for i, ep := range r.endpoints {
		if ep == endpoint {
			r.endpoints = append(r.endpoints[:i], r.endpoints[i+1:]...)
			break
		}
	}
}

func (r *randomLB) GetHealthyEndpoints() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]string, len(r.endpoints))
	copy(result, r.endpoints)
	return result
}

// Weighted Random Load Balancer
type weightedRandomLB struct {
	endpoints []WeightedEndpoint
	rand      *rand.Rand
	mu        sync.RWMutex
}

type WeightedEndpoint struct {
	URL    string
	Weight int
}

func NewWeightedRandomLB(endpoints []string) LoadBalancer {
	weighted := make([]WeightedEndpoint, len(endpoints))
	for i, ep := range endpoints {
		weighted[i] = WeightedEndpoint{URL: ep, Weight: 1}
	}
	
	return &weightedRandomLB{
		endpoints: weighted,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (wr *weightedRandomLB) NextEndpoint() string {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	
	if len(wr.endpoints) == 0 {
		return ""
	}
	
	totalWeight := 0
	for _, ep := range wr.endpoints {
		totalWeight += ep.Weight
	}
	
	if totalWeight == 0 {
		return ""
	}
	
	target := wr.rand.Intn(totalWeight)
	current := 0
	
	for _, ep := range wr.endpoints {
		current += ep.Weight
		if current > target {
			return ep.URL
		}
	}
	
	return wr.endpoints[0].URL
}

func (wr *weightedRandomLB) AddEndpoint(endpoint string) {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	wr.endpoints = append(wr.endpoints, WeightedEndpoint{URL: endpoint, Weight: 1})
}

func (wr *weightedRandomLB) RemoveEndpoint(endpoint string) {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	
	for i, ep := range wr.endpoints {
		if ep.URL == endpoint {
			wr.endpoints = append(wr.endpoints[:i], wr.endpoints[i+1:]...)
			break
		}
	}
}

func (wr *weightedRandomLB) GetHealthyEndpoints() []string {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	
	result := make([]string, len(wr.endpoints))
	for i, ep := range wr.endpoints {
		result[i] = ep.URL
	}
	return result
}

// Least Connection Load Balancer
type leastConnLB struct {
	endpoints   []ConnEndpoint
	mu          sync.RWMutex
}

type ConnEndpoint struct {
	URL         string
	Connections int64
}

func NewLeastConnLB(endpoints []string) LoadBalancer {
	connEndpoints := make([]ConnEndpoint, len(endpoints))
	for i, ep := range endpoints {
		connEndpoints[i] = ConnEndpoint{URL: ep, Connections: 0}
	}
	
	return &leastConnLB{
		endpoints: connEndpoints,
	}
}

func (lc *leastConnLB) NextEndpoint() string {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	
	if len(lc.endpoints) == 0 {
		return ""
	}
	
	minConn := int64(-1)
	selectedIndex := 0
	
	for i, ep := range lc.endpoints {
		if minConn == -1 || ep.Connections < minConn {
			minConn = ep.Connections
			selectedIndex = i
		}
	}
	
	lc.endpoints[selectedIndex].Connections++
	return lc.endpoints[selectedIndex].URL
}

func (lc *leastConnLB) AddEndpoint(endpoint string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.endpoints = append(lc.endpoints, ConnEndpoint{URL: endpoint, Connections: 0})
}

func (lc *leastConnLB) RemoveEndpoint(endpoint string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	
	for i, ep := range lc.endpoints {
		if ep.URL == endpoint {
			lc.endpoints = append(lc.endpoints[:i], lc.endpoints[i+1:]...)
			break
		}
	}
}

func (lc *leastConnLB) GetHealthyEndpoints() []string {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	
	result := make([]string, len(lc.endpoints))
	for i, ep := range lc.endpoints {
		result[i] = ep.URL
	}
	return result
}