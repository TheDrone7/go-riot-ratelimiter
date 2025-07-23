# RIOT Rate Limiter

Some go code to help with rate limiting for the RIOT API.

> **IMPORTANT**
>
> This is not designed to be a public package and maintenance is not guaranteed.
> Most use cases will probably require modifications.

---

## How to Use

```go
store := NewStore() // Create a new store instance
rateLimiter := NewRateLimiter(store) // Create a new rate limiter instance

waitDuration := rateLimiter.CheckFor("https://na1.api.riotgames.com/lol/summoner/v4/summoners/me", "get", "spread")
// url, HTTP method, and strategy (spread/burst)

// Optionally, you can reserve limits if you are planning to make multiple calls asynchronously.
err := rateLimiter.Reserve("https://na1.api.riotgames.com/lol/summoner/v4/summoners/me", "get")

// Wait for the returned duration
time.Sleep(waitDuration)

// ... Perform your API call here ...
resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/me")
// ... handle the response ...

// After the API call, update the rate limiter if there is a response.
// This will also remove the reservation made earlier.
err = rateLimiter.UpdateFromHeaders("https://na1.api.riotgames.com/lol/summoner/v4/summoners/me", "get", resp.Header)

// If there are no headers, you need to manually remove the reservation.
err = rateLimiter.RemoveReservationN("https://na1.api.riotgames.com/lol/summoner/v4/summoners/me", "get", 1)
```

You can also set the rate limits manually if needed:

```go
err = rateLimiter.UpdateRateLimits("https://na1.api.riotgames.com/lol/summoner/v4/summoners/me", "get", myLimits)

// where myLimits is a slice of RateLimits struct defined in ratelimiter.go
```

---

## Files (and modifications)

- helpers.go (Contains helper functions for rate limiting)
- constants.go (Defines constants for rate limiting)
  - Update if necessary to add/remove API methods or platforms (supports all as of 22 July 2025)
- ratelimiter.go (Implements the rate limiting logic)
  - Update if necessary to change rate limiting strategies or logic
- store.go (Implements the storage layer for rate limits)
  - Update if necessary to change storage backend or logic
  - Maybe if you want to add thread safety or redis support

---
