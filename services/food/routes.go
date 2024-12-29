package food

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

type FoodStore interface {
}

type Handler struct {
	store     FoodStore
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(store FoodStore, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{store: store, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	// routeGrp := router.Group(link)

}
