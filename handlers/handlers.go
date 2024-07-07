package handlers

import "github.com/Micah-Shallom/stage-two/config"

type Handlers struct {
	App *config.Application
}

func NewHandlers(app *config.Application) *Handlers {
	return &Handlers{
		App: app,
	}
}
