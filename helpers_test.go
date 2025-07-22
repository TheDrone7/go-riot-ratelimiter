package ratelimiter

import (
	"reflect"
	"testing"
)

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []RateLimitPair
		hasError bool
	}{
		{
			name:     "Valid single pair",
			input:    "100:120",
			expected: []RateLimitPair{{Limit: 100, Duration: 120}},
			hasError: false,
		},
		{
			name:     "Valid multiple pairs",
			input:    "100:120,20:1",
			expected: []RateLimitPair{{Limit: 100, Duration: 120}, {Limit: 20, Duration: 1}},
			hasError: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []RateLimitPair{},
			hasError: false,
		},
		{
			name:     "Invalid format",
			input:    "invalid",
			expected: nil,
			hasError: true,
		},
		{
			name:     "Mixed valid and invalid",
			input:    "100:120,invalid,20:1",
			expected: []RateLimitPair{{Limit: 100, Duration: 120}, {Limit: 20, Duration: 1}},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHeader(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchesPath(t *testing.T) {
	tests := []struct {
		name       string
		urlPath    string
		methodPath string
		expected   bool
	}{
		{
			name:       "Exact match",
			urlPath:    "/lol/summoner/v4/summoners",
			methodPath: "/lol/summoner/v4/summoners",
			expected:   true,
		},
		{
			name:       "With parameter",
			urlPath:    "/lol/summoner/v4/summoners/abc123",
			methodPath: "/lol/summoner/v4/summoners/:puuid",
			expected:   true,
		},
		{
			name:       "Multiple parameters",
			urlPath:    "/riot/account/v1/accounts/by-riot-id/player/tag",
			methodPath: "/riot/account/v1/accounts/by-riot-id/:gameName/:tagLine",
			expected:   true,
		},
		{
			name:       "Different paths",
			urlPath:    "/lol/summoner/v4/summoners",
			methodPath: "/lol/match/v5/matches",
			expected:   false,
		},
		{
			name:       "Different segment count",
			urlPath:    "/lol/summoner/v4",
			methodPath: "/lol/summoner/v4/summoners/:puuid",
			expected:   false,
		},
		{
			name:       "Empty paths",
			urlPath:    "",
			methodPath: "",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesPath(tt.urlPath, tt.methodPath)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for urlPath: %s, methodPath: %s",
					tt.expected, result, tt.urlPath, tt.methodPath)
			}
		})
	}
}

func TestUrlHelper(t *testing.T) {
	tests := []struct {
		name             string
		inputUrl         string
		httpMethod       string
		expectedPlatform string
		expectedService  string
		expectedMethod   string
		hasError         bool
	}{
		{
			name:             "Valid summoner GET request",
			inputUrl:         "https://na1.api.riotgames.com/lol/summoner/v4/summoners/me",
			httpMethod:       "get",
			expectedPlatform: "NA1",
			expectedService:  "SUMMONER",
			expectedMethod:   "GET_BY_ACCESS_TOKEN",
			hasError:         false,
		},
		{
			name:             "Valid account GET request with parameters",
			inputUrl:         "https://europe.api.riotgames.com/riot/account/v1/accounts/by-puuid/some-puuid",
			httpMethod:       "GET",
			expectedPlatform: "EUROPE",
			expectedService:  "ACCOUNT",
			expectedMethod:   "GET_BY_PUUID",
			hasError:         false,
		},
		{
			name:             "Invalid URL",
			inputUrl:         "://invalid-url",
			httpMethod:       "GET",
			expectedPlatform: "",
			expectedService:  "",
			expectedMethod:   "",
			hasError:         true,
		},
		{
			name:             "Unknown endpoint",
			inputUrl:         "https://na1.api.riotgames.com/unknown/endpoint",
			httpMethod:       "GET",
			expectedPlatform: "",
			expectedService:  "",
			expectedMethod:   "",
			hasError:         true,
		},
		{
			name:             "Wrong HTTP method",
			inputUrl:         "https://na1.api.riotgames.com/lol/summoner/v4/summoners/me",
			httpMethod:       "POST",
			expectedPlatform: "",
			expectedService:  "",
			expectedMethod:   "",
			hasError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := urlHelper(tt.inputUrl, tt.httpMethod)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.PlatformName != tt.expectedPlatform {
				t.Errorf("Expected platform %s, got %s", tt.expectedPlatform, result.PlatformName)
			}

			if result.ServiceName != tt.expectedService {
				t.Errorf("Expected service %s, got %s", tt.expectedService, result.ServiceName)
			}

			if result.MethodName != tt.expectedMethod {
				t.Errorf("Expected method %s, got %s", tt.expectedMethod, result.MethodName)
			}
		})
	}
}
