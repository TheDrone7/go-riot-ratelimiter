package ratelimiter

// METHODS defines the Riot API endpoints organized by service
var METHODS = map[string]map[string]string{
	"ACCOUNT": {
		"GET_BY_PUUID":                 "/riot/account/v1/accounts/by-puuid/:puuid",
		"GET_BY_RIOT_ID":               "/riot/account/v1/accounts/by-riot-id/:gameName/:tagLine",
		"GET_BY_ACCESS_TOKEN":          "/riot/account/v1/accounts/me",
		"GET_ACTIVE_SHARD_FOR_PLAYER":  "/riot/account/v1/active-shards/by-game/:game/by-puuid/:puuid",
		"GET_ACTIVE_REGION_FOR_PLAYER": "/riot/account/v1/region/by-game/:game/by-puuid/:puuid",
	},
	"CHAMPION_MASTERY": {
		"GET_ALL_CHAMPIONS":          "/lol/champion-mastery/v4/champion-masteries/by-puuid/:encryptedPUUID",
		"GET_CHAMPION_MASTERY":       "/lol/champion-mastery/v4/champion-masteries/by-puuid/:encryptedPUUID/by-champion/:championId",
		"GET_TOP_CHAMPIONS":          "/lol/champion-mastery/v4/champion-masteries/by-puuid/:encryptedPUUID/top",
		"GET_CHAMPION_MASTERY_SCORE": "/lol/champion-mastery/v4/scores/by-puuid/:encryptedPUUID",
	},
	"CHAMPION": {
		"GET_CHAMPION_ROTATIONS": "/lol/platform/v3/champion-rotations",
	},
	"CLASH": {
		"GET_PLAYERS_BY_PUUID": "/lol/clash/v1/players/by-puuid/:puuid",
		"GET_TEAM":             "/lol/clash/v1/teams/:teamId",
		"GET_TOURNAMENTS":      "/lol/clash/v1/tournaments",
		"GET_TOURNAMENT":       "/lol/clash/v1/tournaments/:tournamentId",
		"GET_TOURNAMENT_TEAM":  "/lol/clash/v1/tournaments/by-team/:teamId",
	},
	"LEAGUE_EXP": {
		"GET_LEAGUE_ENTRIES": "/lol/league-exp/v4/entries/:queue/:tier/:division",
	},
	"LEAGUE": {
		"GET_CHALLENGER_BY_QUEUE":  "/lol/league/v4/challengerleagues/by-queue/:queue",
		"GET_ENTRIES_BY_PUUID":     "/lol/league/v4/entries/by-puuid/:puuid",
		"GET_ALL_ENTRIES":          "/lol/league/v4/entries/:queue/:tier/:division",
		"GET_GRANDMASTER_BY_QUEUE": "/lol/league/v4/grandmasterleagues/by-queue/:queue",
		"GET_LEAGUE_BY_ID":         "/lol/league/v4/leagues/:leagueId",
		"GET_MASTER_BY_QUEUE":      "/lol/league/v4/masterleagues/by-queue/:queue",
	},
	"LOL_CHALLENGES": {
		"GET_CONFIG":               "/lol/challenges/v1/challenges/config",
		"GET_PERCENTILES":          "/lol/challenges/v1/challenges/percentiles",
		"GET_CONFIG_BY_ID":         "/lol/challenges/v1/challenges/:challengeId/config",
		"GET_LEADERBOARD_BY_ID":    "/lol/challenges/v1/challenges/:challengeId/leaderboards/by-level/:level",
		"GET_PERCENTILES_BY_ID":    "/lol/challenges/v1/challenges/:challengeId/percentiles",
		"GET_PLAYER_DATA_BY_PUUID": "/lol/challenges/v1/player-data/:puuid",
	},
	"LOL_RSO_MATCH": {
		"GET_MATCH_IDS_BY_ACCESS_TOKEN": "/lol/rso-match/v1/matches/ids",
		"GET_MATCH_BY_ID":               "/lol/rso-match/v1/matches/:matchId",
		"GET_MATCH_TIMELINE_BY_ID":      "/lol/rso-match/v1/matches/:matchId/timeline",
	},
	"LOL_STATUS": {
		"GET_PLATFORM_DATA": "/lol/status/v4/platform-data",
	},
	"LOR_DECK": {
		"GET_DECKS_FOR_PLAYER":        "/lor/deck/v1/decks/me",
		"POST_CREATE_DECK_FOR_PLAYER": "/lor/deck/v1/decks/me",
	},
	"LOR_INVENTORY": {
		"GET_CARDS_OWNED_BY_PLAYER": "/lor/inventory/v1/cards/me",
	},
	"LOR_MATCH": {
		"GET_MATCH_IDS_BY_PUUID": "/lor/match/v1/matches/by-puuid/:puuid/ids",
		"GET_MATCH_BY_ID":        "/lor/match/v1/matches/:matchId",
	},
	"LOR_RANKED": {
		"GET_MASTER_TIER": "/lor/ranked/v1/leaderboards",
	},
	"LOR_STATUS_V1": {
		"GET_PLATFORM_DATA": "/lor/status/v1/platform-data",
	},
	"MATCH_V5": {
		"GET_IDS_BY_PUUID":         "/lol/match/v5/matches/by-puuid/:puuid/ids",
		"GET_MATCH_BY_ID":          "/lol/match/v5/matches/:matchId",
		"GET_MATCH_TIMELINE_BY_ID": "/lol/match/v5/matches/:matchId/timeline",
	},
	"RIFTBOUND_CONTENT": {
		"GET_RIFTBOUND_CONTENT": "/riftbound-content/v1/contents",
	},
	"SPECTATOR_TFT_V5": {
		"GET_GAME_BY_PUUID":  "/lol/spectator/tft/v5/active-games/by-puuid/:puuid",
		"GET_FEATURED_GAMES": "/lol/spectator/tft/v5/featured-games",
	},
	"SPECTATOR": {
		"GET_GAME_BY_PUUID":  "/lol/spectator/v5/active-games/by-summoner/:puuid",
		"GET_FEATURED_GAMES": "/lol/spectator/v5/featured-games",
	},
	"SUMMONER": {
		"GET_BY_ACCESS_TOKEN": "/lol/summoner/v4/summoners/me",
		"GET_BY_PUUID":        "/lol/summoner/v4/summoners/by-puuid/:puuid",
	},
	"TFT_LEAGUE": {
		"GET_BY_PUUID":                  "/tft/league/v1/by-puuid/:puuid",
		"GET_CHALLENGER":                "/tft/league/v1/challenger",
		"GET_ALL_ENTRIES":               "/tft/league/v1/entries/:tier/:division",
		"GET_GRANDMASTER":               "/tft/league/v1/grandmaster",
		"GET_MASTER":                    "/tft/league/v1/master",
		"GET_TOP_RATED_LADDER_BY_QUEUE": "/tft/league/v1/rated-ladders/:queue/top",
		"GET_LEAGUE_BY_ID":              "/tft/league/v1/leagues/:leagueId",
	},
	"TFT_MATCH": {
		"GET_MATCH_IDS_BY_PUUID": "/tft/match/v1/matches/by-puuid/:puuid/ids",
		"GET_MATCH_BY_ID":        "/tft/match/v1/matches/:matchId",
	},
	"TFT_STATUS_V1": {
		"GET_PLATFORM_DATA": "/tft/status/v1/platform-data",
	},
	"TFT_SUMMONER": {
		"GET_BY_PUUID":        "/tft/summoner/v1/summoners/by-puuid/:puuid",
		"GET_BY_ACCESS_TOKEN": "/tft/summoner/v1/summoners/me",
	},
	"TOURNAMENT_STUB_V5": {
		"POST_CREATE_CODES":                   "/lol/tournament-stub/v5/codes",
		"GET_TOURNAMENT_BY_CODE":              "/lol/tournament-stub/v5/codes/:tournamentCode",
		"GET_LOBBY_EVENTS_BY_TOURNAMENT_CODE": "/lol/tournament-stub/v5/lobby-events/by-code/:tournamentCode",
		"POST_CREATE_PROVIDER":                "/lol/tournament-stub/v5/providers",
		"POST_CREATE_TOURNAMENT":              "/lol/tournament-stub/v5/tournaments",
	},
	"TOURNAMENT_V5": {
		"POST_CREATE_CODES":                   "/lol/tournament/v5/codes",
		"GET_TOURNAMENT_BY_CODE":              "/lol/tournament/v5/codes/:tournamentCode",
		"PUT_TOURNAMENT_CODE":                 "/lol/tournament/v5/codes/:tournamentCode",
		"GET_TOURNAMENT_GAME_DETAILS":         "/lol/tournament/v5/games/by-code/:tournamentCode",
		"GET_LOBBY_EVENTS_BY_TOURNAMENT_CODE": "/lol/tournament/v5/lobby-events/by-code/:tournamentCode",
		"POST_CREATE_PROVIDER":                "/lol/tournament/v5/providers",
		"POST_CREATE_TOURNAMENT":              "/lol/tournament/v5/tournaments",
	},
	"VAL_CONSOLE_MATCH": {
		"GET_MATCH_BY_ID":             "/val/match/console/v1/matches/:matchId",
		"GET_MATCHLIST_BY_PUUID":      "/val/match/console/v1/matchlists/by-puuid/:puuid",
		"GET_RECENT_MATCHES_BY_QUEUE": "/val/match/console/v1/recent-matches/by-queue/:queue",
	},
	"VAL_CONSOLE_RANKED": {
		"GET_LEADERBOARD_BY_QUEUE": "/val/console/ranked/v1/leaderboards/by-act/:actId",
	},
	"VAL_CONTENT": {
		"GET_CONTENT": "/val/content/v1/contents",
	},
	"VAL_MATCH": {
		"GET_MATCH_BY_ID":             "/val/match/v1/matches/:matchId",
		"GET_MATCHLIST_BY_PUUID":      "/val/match/v1/matchlists/by-puuid/:puuid",
		"GET_RECENT_MATCHES_BY_QUEUE": "/val/match/v1/recent-matches/by-queue/:queue",
	},
	"VAL_RANKED": {
		"GET_LEADERBOARD_BY_QUEUE": "/val/ranked/v1/leaderboards/by-act/:actId",
	},
	"VAL_STATUS_V1": {
		"GET_PLATFORM_DATA": "/val/status/v1/platform-data",
	},
}

type LimitType string

const (
	LIMIT_TYPE_APPLICATION LimitType = "application"
	LIMIT_TYPE_METHOD      LimitType = "method"
)

type LimitStrategy string

const (
	LIMIT_STRATEGY_SPREAD LimitStrategy = "spread"
	LIMIT_STRATEGY_BURST  LimitStrategy = "burst"
)
