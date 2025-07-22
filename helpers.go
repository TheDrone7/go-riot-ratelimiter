package ratelimiter

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// RateLimitPair represents a pair of numbers
type RateLimitPair struct {
	Limit    int
	Duration int
}

type RateLimitDetails struct {
	PlatformName string
	ServiceName  string
	MethodName   string
}

// Parses the ratelimit header string
// Expected format: "100:120,20:1"
func parseHeader(input string) ([]RateLimitPair, error) {
	if input == "" {
		return []RateLimitPair{}, nil
	}

	var pairs []RateLimitPair

	// Use regex to find all number pairs in format "number:number"
	re := regexp.MustCompile(`(\d+):(\d+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		return nil, errors.New("no valid number pairs found in input")
	}

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		first, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, errors.New("invalid limit/count number: " + match[1])
		}

		second, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, errors.New("invalid duration number: " + match[2])
		}

		pairs = append(pairs, RateLimitPair{Limit: first, Duration: second})
	}

	return pairs, nil
}

// Checks if a URL path matches a method path template (e.g., "/lol/summoner/v4/summoners/:puuid")
func matchesPath(urlPath string, methodPath string) bool {
	urlSegments := strings.Split(strings.Trim(urlPath, "/"), "/")
	methodSegments := strings.Split(strings.Trim(methodPath, "/"), "/")

	if len(urlSegments) != len(methodSegments) {
		return false
	}

	for i, methodSegment := range methodSegments {
		// Skip parameter segments (those starting with :)
		if strings.HasPrefix(methodSegment, ":") {
			continue
		}

		// Must match exactly for non-parameter segments
		if urlSegments[i] != methodSegment {
			return false
		}
	}

	return true
}

// Validates a URL with HTTP method and returns a RateLimitDetails object with extracted platform and path
func urlHelper(inputUrl string, httpMethod string) (*RateLimitDetails, error) {
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		return nil, errors.New("invalid URL format: " + err.Error())
	}

	host := parsedUrl.Host
	hostParts := strings.Split(host, ".")
	platform := strings.ToUpper(hostParts[0])
	path := parsedUrl.Path

	// Find matching service and method by iterating through METHODS
	serviceName := ""
	methodName := ""

	for service, methods := range METHODS {
		for method, methodPath := range methods {
			if matchesPath(path, methodPath) && strings.HasPrefix(method, strings.ToUpper(httpMethod)) {
				serviceName = service
				methodName = method
				break
			}
		}
		if serviceName != "" {
			break
		}
	}

	// Return error if no matching endpoint is found
	if serviceName == "" || methodName == "" {
		return nil, errors.New("unknown endpoint: " + httpMethod + " " + path)
	}

	return &RateLimitDetails{
		PlatformName: platform,
		ServiceName:  serviceName,
		MethodName:   methodName,
	}, nil
}
