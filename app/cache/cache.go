package cache

import (
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/types"
)

type RecipeCache struct {
	recipes []types.Recipe
	logger  *logging.Logger
}

func NewRecipeCache(logger *logging.Logger) RecipeCache {
	return RecipeCache{logger: logger}
}

func (r *RecipeCache) Get(isAdmin bool) *[]types.Recipe {
	var recipes *[]types.Recipe
	if isAdmin {
		recipes = r.getWithPending()
	} else {
		recipes = r.getWithoutPending()
	}
	r.logger.Debugf("read %d recipes from cache", len(*recipes))
	return recipes
}

func (r *RecipeCache) getWithPending() *[]types.Recipe {
	if r.recipes == nil {
		return &[]types.Recipe{}
	}
	return &r.recipes
}

func (r *RecipeCache) getWithoutPending() *[]types.Recipe {
	var recipes []types.Recipe
	if r.recipes == nil {
		return &recipes
	}
	for _, recipe := range r.recipes {
		if !recipe.Pending {
			recipes = append(recipes, recipe)
		}
	}
	return &recipes
}

func (r *RecipeCache) Set(recipes []types.Recipe) {
	r.recipes = recipes
	r.logger.Debugf("put %d recipes into cache", len(recipes))
}

func (r *RecipeCache) Invalidate() {
	r.recipes = nil
	r.logger.Debug("invalidated cache")
}

func (r *RecipeCache) IsValid() bool {
	return r.recipes != nil
}
