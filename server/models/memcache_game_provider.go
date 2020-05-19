package models

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// MemcacheGameProvider provides game storage through in-memory caching
type MemcacheGameProvider struct {
	internalCache *cache.Cache
}

var (
	providerInstance = &MemcacheGameProvider{
		internalCache: cache.New(20*time.Minute, 5*time.Minute),
	}
)

func getMemcacheGameProviderInstance() *MemcacheGameProvider {
	return providerInstance
}

func (provider *MemcacheGameProvider) LoadGame(groupName string) *Game {
	state, found := provider.internalCache.Get(groupName)
	if found {
		return state.(*Game)
	}
	return nil
}

func (provider *MemcacheGameProvider) SaveGame(game *Game) error {
	provider.internalCache.Set(game.GroupName, game, cache.DefaultExpiration)
	return nil
}
