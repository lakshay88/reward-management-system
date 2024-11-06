package routers

import (
	"github.com/go-chi/chi/v5"
)

type Routers struct{}

func NewRouter() *Routers {
	return &Routers{}
}

func (r *Routers) RegisterRoutes(router *chi.Mux) {

	// router.Mount("/api", api.OrderRoutes{}.Routes())

}
