package ai

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

// SmartRetry uses AI to determine optimal retry strategies
type SmartRetry struct {
	history    []RetryAttempt
	mu         sync.RWMutex
	model      *RetryModel
}

type RetryAttempt struct {
	URL        string
	Method     string
	StatusCode int
	Duration   time.Duration
	Success    bool
	Timestamp  time.Time
}

type RetryModel struct {
	weights map[string]float64
	bias    float64
}

func NewSmartRetry() *SmartRetry {
	return &SmartRetry{
		history: make([]RetryAttempt, 0),
		model: &RetryModel{
			weights: map[string]float64{
				"status_code": -0.1,
				"duration":    -0.05,
				"hour":        0.02,
				"method":      0.1,
			},
			bias: 0.5,
		},
	}
}

func (sr *SmartRetry) ShouldRetry(req *http.Request, resp *http.Response, attempt int) bool {
	if attempt >= 5 {
		return false
	}

	// Use AI model to predict success probability
	probability := sr.predictSuccessProbability(req, resp)
	
	// Adaptive threshold based on attempt number
	threshold := 0.3 + (float64(attempt) * 0.1)
	
	return probability > threshold
}

func (sr *SmartRetry) predictSuccessProbability(req *http.Request, resp *http.Response) float64 {
	features := sr.extractFeatures(req, resp)
	
	score := sr.model.bias
	for feature, value := range features {
		if weight, exists := sr.model.weights[feature]; exists {
			score += weight * value
		}
	}
	
	// Sigmoid activation
	return 1.0 / (1.0 + math.Exp(-score))
}

func (sr *SmartRetry) extractFeatures(req *http.Request, resp *http.Response) map[string]float64 {
	features := make(map[string]float64)
	
	if resp != nil {
		features["status_code"] = float64(resp.StatusCode)
	}
	
	features["hour"] = float64(time.Now().Hour())
	
	switch req.Method {
	case "GET":
		features["method"] = 1.0
	case "POST":
		features["method"] = 0.8
	case "PUT":
		features["method"] = 0.6
	default:
		features["method"] = 0.5
	}
	
	return features
}

func (sr *SmartRetry) RecordAttempt(req *http.Request, resp *http.Response, duration time.Duration, success bool) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	
	attempt := RetryAttempt{
		URL:       req.URL.String(),
		Method:    req.Method,
		Duration:  duration,
		Success:   success,
		Timestamp: time.Now(),
	}
	
	if resp != nil {
		attempt.StatusCode = resp.StatusCode
	}
	
	sr.history = append(sr.history, attempt)
	
	// Keep only recent history
	if len(sr.history) > 1000 {
		sr.history = sr.history[len(sr.history)-1000:]
	}
	
	// Update model periodically
	if len(sr.history)%100 == 0 {
		sr.updateModel()
	}
}

func (sr *SmartRetry) updateModel() {
	// Simple online learning update
	learningRate := 0.01
	
	for _, attempt := range sr.history[len(sr.history)-10:] {
		features := map[string]float64{
			"status_code": float64(attempt.StatusCode),
			"duration":    attempt.Duration.Seconds(),
			"hour":        float64(attempt.Timestamp.Hour()),
		}
		
		predicted := sr.model.bias
		for feature, value := range features {
			if weight, exists := sr.model.weights[feature]; exists {
				predicted += weight * value
			}
		}
		
		actual := 0.0
		if attempt.Success {
			actual = 1.0
		}
		
		error := actual - predicted
		
		// Update weights
		for feature, value := range features {
			if _, exists := sr.model.weights[feature]; exists {
				sr.model.weights[feature] += learningRate * error * value
			}
		}
		sr.model.bias += learningRate * error
	}
}

// SmartCache uses AI to optimize caching decisions
type SmartCache struct {
	accessPatterns map[string]*AccessPattern
	mu             sync.RWMutex
}

type AccessPattern struct {
	URL           string
	AccessTimes   []time.Time
	HitRate       float64
	LastAccess    time.Time
	PredictedNext time.Time
}

func NewSmartCache() *SmartCache {
	return &SmartCache{
		accessPatterns: make(map[string]*AccessPattern),
	}
}

func (sc *SmartCache) ShouldCache(url string, size int64) bool {
	sc.mu.RLock()
	pattern, exists := sc.accessPatterns[url]
	sc.mu.RUnlock()
	
	if !exists {
		// Cache new URLs by default
		return size < 1024*1024 // 1MB limit for new items
	}
	
	// Use access frequency and recency to decide
	frequency := sc.calculateFrequency(pattern)
	recency := time.Since(pattern.LastAccess).Hours()
	
	score := frequency * math.Exp(-recency/24) // Decay over 24 hours
	
	return score > 0.1 && size < 10*1024*1024 // 10MB limit for frequent items
}

func (sc *SmartCache) calculateFrequency(pattern *AccessPattern) float64 {
	if len(pattern.AccessTimes) < 2 {
		return 0.1
	}
	
	// Calculate access frequency over the last 24 hours
	now := time.Now()
	recentAccesses := 0
	
	for _, accessTime := range pattern.AccessTimes {
		if now.Sub(accessTime).Hours() <= 24 {
			recentAccesses++
		}
	}
	
	return float64(recentAccesses) / 24.0 // Accesses per hour
}

func (sc *SmartCache) RecordAccess(url string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	pattern, exists := sc.accessPatterns[url]
	if !exists {
		pattern = &AccessPattern{
			URL:         url,
			AccessTimes: make([]time.Time, 0),
		}
		sc.accessPatterns[url] = pattern
	}
	
	now := time.Now()
	pattern.AccessTimes = append(pattern.AccessTimes, now)
	pattern.LastAccess = now
	
	// Keep only recent access times
	if len(pattern.AccessTimes) > 100 {
		pattern.AccessTimes = pattern.AccessTimes[len(pattern.AccessTimes)-100:]
	}
	
	// Predict next access time
	pattern.PredictedNext = sc.predictNextAccess(pattern)
}

func (sc *SmartCache) predictNextAccess(pattern *AccessPattern) time.Time {
	if len(pattern.AccessTimes) < 3 {
		return time.Now().Add(time.Hour)
	}
	
	// Calculate average interval between accesses
	intervals := make([]time.Duration, 0)
	for i := 1; i < len(pattern.AccessTimes); i++ {
		interval := pattern.AccessTimes[i].Sub(pattern.AccessTimes[i-1])
		intervals = append(intervals, interval)
	}
	
	// Calculate weighted average (more recent intervals have higher weight)
	totalWeight := 0.0
	weightedSum := time.Duration(0)
	
	for i, interval := range intervals {
		weight := math.Exp(float64(i-len(intervals)) * 0.1) // Exponential decay
		totalWeight += weight
		weightedSum += time.Duration(float64(interval) * weight)
	}
	
	avgInterval := time.Duration(float64(weightedSum) / totalWeight)
	
	return pattern.LastAccess.Add(avgInterval)
}

// AdaptiveTimeout adjusts timeouts based on historical performance
type AdaptiveTimeout struct {
	endpointStats map[string]*EndpointStats
	mu            sync.RWMutex
}

type EndpointStats struct {
	URL              string
	ResponseTimes    []time.Duration
	SuccessRate      float64
	RecommendedTimeout time.Duration
	LastUpdate       time.Time
}

func NewAdaptiveTimeout() *AdaptiveTimeout {
	return &AdaptiveTimeout{
		endpointStats: make(map[string]*EndpointStats),
	}
}

func (at *AdaptiveTimeout) GetTimeout(url string, defaultTimeout time.Duration) time.Duration {
	at.mu.RLock()
	stats, exists := at.endpointStats[url]
	at.mu.RUnlock()
	
	if !exists || time.Since(stats.LastUpdate) > time.Hour {
		return defaultTimeout
	}
	
	return stats.RecommendedTimeout
}

func (at *AdaptiveTimeout) RecordResponse(url string, duration time.Duration, success bool) {
	at.mu.Lock()
	defer at.mu.Unlock()
	
	stats, exists := at.endpointStats[url]
	if !exists {
		stats = &EndpointStats{
			URL:           url,
			ResponseTimes: make([]time.Duration, 0),
			SuccessRate:   1.0,
		}
		at.endpointStats[url] = stats
	}
	
	stats.ResponseTimes = append(stats.ResponseTimes, duration)
	stats.LastUpdate = time.Now()
	
	// Keep only recent response times
	if len(stats.ResponseTimes) > 100 {
		stats.ResponseTimes = stats.ResponseTimes[len(stats.ResponseTimes)-100:]
	}
	
	// Update success rate
	successCount := 0
	for i := len(stats.ResponseTimes) - 10; i < len(stats.ResponseTimes) && i >= 0; i++ {
		if i == len(stats.ResponseTimes)-1 && success {
			successCount++
		} else if i < len(stats.ResponseTimes)-1 {
			successCount++ // Assume previous responses were successful if recorded
		}
	}
	
	if len(stats.ResponseTimes) >= 10 {
		stats.SuccessRate = float64(successCount) / 10.0
	}
	
	// Calculate recommended timeout
	stats.RecommendedTimeout = at.calculateOptimalTimeout(stats)
}

func (at *AdaptiveTimeout) calculateOptimalTimeout(stats *EndpointStats) time.Duration {
	if len(stats.ResponseTimes) < 5 {
		return 30 * time.Second
	}
	
	// Calculate 95th percentile response time
	times := make([]time.Duration, len(stats.ResponseTimes))
	copy(times, stats.ResponseTimes)
	
	// Simple sort for percentile calculation
	for i := 0; i < len(times); i++ {
		for j := i + 1; j < len(times); j++ {
			if times[i] > times[j] {
				times[i], times[j] = times[j], times[i]
			}
		}
	}
	
	p95Index := int(float64(len(times)) * 0.95)
	if p95Index >= len(times) {
		p95Index = len(times) - 1
	}
	
	p95Time := times[p95Index]
	
	// Add buffer based on success rate
	buffer := time.Duration(float64(p95Time) * (2.0 - stats.SuccessRate))
	
	recommendedTimeout := p95Time + buffer
	
	// Ensure reasonable bounds
	if recommendedTimeout < 5*time.Second {
		recommendedTimeout = 5 * time.Second
	}
	if recommendedTimeout > 300*time.Second {
		recommendedTimeout = 300 * time.Second
	}
	
	return recommendedTimeout
}

// PredictivePreloader anticipates future requests
type PredictivePreloader struct {
	requestPatterns map[string]*RequestPattern
	mu              sync.RWMutex
	preloadFunc     func(url string)
}

type RequestPattern struct {
	BaseURL       string
	FollowupURLs  map[string]float64 // URL -> probability
	LastRequests  []string
	PatternLength int
}

func NewPredictivePreloader(preloadFunc func(url string)) *PredictivePreloader {
	return &PredictivePreloader{
		requestPatterns: make(map[string]*RequestPattern),
		preloadFunc:     preloadFunc,
	}
}

func (pp *PredictivePreloader) RecordRequest(url string) {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	
	// Update patterns for all recent requests
	for baseURL, pattern := range pp.requestPatterns {
		if len(pattern.LastRequests) > 0 {
			lastURL := pattern.LastRequests[len(pattern.LastRequests)-1]
			if lastURL == baseURL {
				// This URL followed the base URL
				if pattern.FollowupURLs == nil {
					pattern.FollowupURLs = make(map[string]float64)
				}
				pattern.FollowupURLs[url] += 0.1
			}
		}
	}
	
	// Create or update pattern for this URL
	pattern, exists := pp.requestPatterns[url]
	if !exists {
		pattern = &RequestPattern{
			BaseURL:       url,
			FollowupURLs:  make(map[string]float64),
			LastRequests:  make([]string, 0),
			PatternLength: 3,
		}
		pp.requestPatterns[url] = pattern
	}
	
	pattern.LastRequests = append(pattern.LastRequests, url)
	if len(pattern.LastRequests) > pattern.PatternLength {
		pattern.LastRequests = pattern.LastRequests[1:]
	}
	
	// Trigger preloading for likely next requests
	pp.triggerPreloading(url)
}

func (pp *PredictivePreloader) triggerPreloading(currentURL string) {
	if pp.preloadFunc == nil {
		return
	}
	
	pattern, exists := pp.requestPatterns[currentURL]
	if !exists {
		return
	}
	
	// Preload URLs with high probability
	for url, probability := range pattern.FollowupURLs {
		if probability > 0.5 {
			go pp.preloadFunc(url)
		}
	}
}

// AIManager coordinates all AI features
type AIManager struct {
	smartRetry          *SmartRetry
	smartCache          *SmartCache
	adaptiveTimeout     *AdaptiveTimeout
	predictivePreloader *PredictivePreloader
	enabled             bool
}

func NewAIManager() *AIManager {
	return &AIManager{
		smartRetry:      NewSmartRetry(),
		smartCache:      NewSmartCache(),
		adaptiveTimeout: NewAdaptiveTimeout(),
		enabled:         true,
	}
}

func (ai *AIManager) SetPreloadFunction(fn func(url string)) {
	ai.predictivePreloader = NewPredictivePreloader(fn)
}

func (ai *AIManager) ShouldRetry(req *http.Request, resp *http.Response, attempt int) bool {
	if !ai.enabled {
		return attempt < 3 // Fallback to simple retry
	}
	return ai.smartRetry.ShouldRetry(req, resp, attempt)
}

func (ai *AIManager) ShouldCache(url string, size int64) bool {
	if !ai.enabled {
		return size < 1024*1024 // Fallback to simple size check
	}
	return ai.smartCache.ShouldCache(url, size)
}

func (ai *AIManager) GetAdaptiveTimeout(url string, defaultTimeout time.Duration) time.Duration {
	if !ai.enabled {
		return defaultTimeout
	}
	return ai.adaptiveTimeout.GetTimeout(url, defaultTimeout)
}

func (ai *AIManager) RecordRequest(req *http.Request, resp *http.Response, duration time.Duration, success bool) {
	if !ai.enabled {
		return
	}
	
	url := req.URL.String()
	
	ai.smartRetry.RecordAttempt(req, resp, duration, success)
	ai.smartCache.RecordAccess(url)
	ai.adaptiveTimeout.RecordResponse(url, duration, success)
	
	if ai.predictivePreloader != nil {
		ai.predictivePreloader.RecordRequest(url)
	}
}

func (ai *AIManager) Enable() {
	ai.enabled = true
}

func (ai *AIManager) Disable() {
	ai.enabled = false
}