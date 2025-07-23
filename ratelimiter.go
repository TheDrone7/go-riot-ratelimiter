package ratelimiter

import (
	"net/http"
	"strconv"
	"time"
)

type RateLimits struct {
	Limit      int
	Counts     int
	Duration   time.Duration
	RetryAfter time.Duration
	LastAt     time.Time
}

// RateLimiter represents the rate limiting functionality
type RateLimiter struct {
	cache Store
}

func NewRateLimiter(store Store) *RateLimiter {
	return &RateLimiter{
		cache: store,
	}
}

// Reserve creates a reservation for a URL and method, incrementing the reservation count
func (rl *RateLimiter) Reserve(url string, method string) error {
	details, err := urlHelper(url, method)
	if err != nil {
		return err
	}

	platform := details.PlatformName
	methodKey := platform + ":" + details.ServiceName + ":" + details.MethodName

	// Platform reserve key
	platformReserveKey := platform + ":reserve"
	if reserveCountRaw, exists := rl.cache.Get(platformReserveKey); exists {
		if count, ok := reserveCountRaw.(int); ok {
			rl.cache.Set(platformReserveKey, count+1)
		}
	} else {
		rl.cache.Set(platformReserveKey, 1)
	}

	// Method reserve key
	methodReserveKey := methodKey + ":reserve"
	if reserveCountRaw, exists := rl.cache.Get(methodReserveKey); exists {
		if count, ok := reserveCountRaw.(int); ok {
			rl.cache.Set(methodReserveKey, count+1)
		}
	} else {
		rl.cache.Set(methodReserveKey, 1)
	}

	return nil
}

// RemoveReservationN reduces reservations for a URL and method by n (but not lower than 0)
func (rl *RateLimiter) RemoveReservationN(url string, method string, n int) error {
	details, err := urlHelper(url, method)
	if err != nil {
		return err
	}

	platform := details.PlatformName
	methodKey := platform + ":" + details.ServiceName + ":" + details.MethodName

	// Reduce platform reservation by n
	platformReserveKey := platform + ":reserve"
	if reserveCountRaw, exists := rl.cache.Get(platformReserveKey); exists {
		if count, ok := reserveCountRaw.(int); ok {
			newCount := count - n
			if newCount < 0 {
				newCount = 0
			}
			rl.cache.Set(platformReserveKey, newCount)
		}
	}

	// Reduce method reservation by n
	methodReserveKey := methodKey + ":reserve"
	if reserveCountRaw, exists := rl.cache.Get(methodReserveKey); exists {
		if count, ok := reserveCountRaw.(int); ok {
			newCount := count - n
			if newCount < 0 {
				newCount = 0
			}
			rl.cache.Set(methodReserveKey, newCount)
		}
	}

	return nil
}

// Extracts platform, service, and method names from the URL and method
// Then updates the ratelimits in the cache
// Returns an error if the URL or method is invalid
func (rl *RateLimiter) UpdateRateLimits(url string, method string, limitType LimitType, limits []RateLimits) error {
	details, err := urlHelper(url, method)
	if err != nil {
		return err
	}

	platform := details.PlatformName
	methodKey := platform + ":" + details.ServiceName + ":" + details.MethodName

	switch limitType {
	case LIMIT_TYPE_METHOD:
		rl.cache.Set(methodKey, limits)
	case LIMIT_TYPE_APPLICATION:
		rl.cache.Set(platform, limits)
	}

	return nil
}

// Updates the rate limits based on URL, HTTP method and response headers
func (rl *RateLimiter) UpdateFromHeaders(url string, method string, headers http.Header) error {
	now := time.Now()

	// Extract rate limit headers with default values
	appRateLimit := headers.Get("X-App-Rate-Limit")
	if appRateLimit == "" {
		appRateLimit = "100:120,20:1"
	}

	appRateLimitCount := headers.Get("X-App-Rate-Limit-Count")
	if appRateLimitCount == "" {
		appRateLimitCount = "1:120,1:1"
	}

	methodRateLimit := headers.Get("X-Method-Rate-Limit")
	if methodRateLimit == "" {
		methodRateLimit = "100:120,20:1"
	}

	methodRateLimitCount := headers.Get("X-Method-Rate-Limit-Count")
	if methodRateLimitCount == "" {
		methodRateLimitCount = "1:120,1:1"
	}

	retryAfterStr := headers.Get("Retry-After")
	if retryAfterStr == "" {
		retryAfterStr = "0" // Default to 0 seconds if not provided
	}

	retryAfter, err := strconv.ParseFloat(retryAfterStr, 64)
	if err != nil {
		return err
	}

	rl.RemoveReservationN(url, method, 1)

	appLimitPairs, err := parseHeader(appRateLimit)
	if err != nil {
		return err
	}

	appCountPairs, err := parseHeader(appRateLimitCount)
	if err != nil {
		return err
	}

	methodLimitPairs, err := parseHeader(methodRateLimit)
	if err != nil {
		return err
	}

	methodCountPairs, err := parseHeader(methodRateLimitCount)
	if err != nil {
		return err
	}

	var appRateLimits []RateLimits
	for i, limitPair := range appLimitPairs {
		rateLimits := RateLimits{
			Limit:      limitPair.Limit,
			Duration:   time.Duration(limitPair.Duration) * time.Second,
			RetryAfter: time.Duration(retryAfter) * time.Second,
			LastAt:     now,
		}

		if i < len(appCountPairs) {
			rateLimits.Counts = appCountPairs[i].Limit
		} else {
			rateLimits.Counts = 0
		}

		appRateLimits = append(appRateLimits, rateLimits)
	}

	var methodRateLimits []RateLimits
	for i, limitPair := range methodLimitPairs {
		rateLimits := RateLimits{
			Limit:      limitPair.Limit,
			Duration:   time.Duration(limitPair.Duration) * time.Second,
			RetryAfter: time.Duration(retryAfter) * time.Second,
			LastAt:     now,
		}

		if i < len(methodCountPairs) {
			rateLimits.Counts = methodCountPairs[i].Limit
		} else {
			rateLimits.Counts = 0
		}

		methodRateLimits = append(methodRateLimits, rateLimits)
	}

	rl.UpdateRateLimits(url, method, LIMIT_TYPE_APPLICATION, appRateLimits)
	rl.UpdateRateLimits(url, method, LIMIT_TYPE_METHOD, methodRateLimits)

	return nil
}

// GetWaitFor calculates the wait time for a given URL, HTTP method, and limit strategy
func (rl *RateLimiter) GetWaitFor(url string, httpMethod string, strategy LimitStrategy) (time.Duration, error) {
	// Parse URL and method to get platform, service and method details
	details, err := urlHelper(url, httpMethod)
	if err != nil {
		return 0, err
	}

	now := time.Now()

	// Generate keys for cache lookup
	platform := details.PlatformName
	methodKey := platform + ":" + details.ServiceName + ":" + details.MethodName

	// Get application limits from cache
	appLimits := []RateLimits{}
	if appLimitsRaw, exists := rl.cache.Get(platform); exists {
		if limits, ok := appLimitsRaw.([]RateLimits); ok {
			appLimits = limits
		}
	}

	// Get method limits from cache
	methodLimits := []RateLimits{}
	if methodLimitsRaw, exists := rl.cache.Get(methodKey); exists {
		if limits, ok := methodLimitsRaw.([]RateLimits); ok {
			methodLimits = limits
		}
	}

	platformReserveCount := 0
	if reserveCountRaw, exists := rl.cache.Get(platform + ":reserve"); exists {
		if count, ok := reserveCountRaw.(int); ok {
			platformReserveCount = count
		}
	}

	methodReserveCount := 0
	if reserveCountRaw, exists := rl.cache.Get(methodKey + ":reserve"); exists {
		if count, ok := reserveCountRaw.(int); ok {
			methodReserveCount = count
		}
	}

	allLimits := append(appLimits, methodLimits...)
	waitTime := time.Duration(0)

	if strategy == LIMIT_STRATEGY_BURST {
		for i, limit := range allLimits {
			// Add reservation count to effective counts
			reserveCount := platformReserveCount
			if i >= len(appLimits) {
				// This is a method limit, use method reserve count
				reserveCount = methodReserveCount
			}

			// Check if current count + reservations >= limit
			if limit.Counts+reserveCount >= limit.Limit {
				timeElapsed := now.Sub(limit.LastAt)
				if timeElapsed < limit.Duration {
					tempWait := limit.Duration - timeElapsed
					if tempWait > waitTime {
						waitTime = tempWait
					}
				}
			}
		}
	} else {
		for i, limit := range allLimits {
			// Add reservation count to effective counts
			reserveCount := platformReserveCount
			if i >= len(appLimits) {
				// This is a method limit, use method reserve count
				reserveCount = methodReserveCount
			}

			effectiveCounts := limit.Counts + reserveCount

			if effectiveCounts >= limit.Limit {
				timeElapsed := now.Sub(limit.LastAt)
				if timeElapsed < limit.Duration {
					tempWait := limit.Duration - timeElapsed
					if tempWait > waitTime {
						waitTime = tempWait
					}
				}
			} else {
				timeElapsed := now.Sub(limit.LastAt)
				if timeElapsed < limit.Duration {
					remainingRequests := limit.Limit - effectiveCounts
					remainingTime := limit.Duration - timeElapsed

					if remainingRequests > 0 && remainingTime > 0 {
						averageWaitTime := remainingTime / time.Duration(remainingRequests)
						if averageWaitTime > waitTime {
							waitTime = averageWaitTime
						}
					}
				}
			}
		}
	}

	// Return the calculated wait time
	return waitTime - time.Since(now), nil
}
