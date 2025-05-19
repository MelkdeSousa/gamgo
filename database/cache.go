package database

import "strings"

const CACHE_SEARCH_GAME_KEY_PREFIX = "search:game"

func GetCacheKey(key ...string) string {
	return strings.Join(key, ":")
}
