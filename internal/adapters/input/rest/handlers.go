package rest

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../oapi-codegen-types.yaml ../../../../openapi.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../oapi-codegen.yaml ../../../../openapi.yaml

import (
	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest/gen"
	"github.com/charlires/go-base-backend-service/internal/core/services"
)

// Handlers implements gen.StrictServerInterface and wires all service dependencies.
type Handlers struct {
	UserService     services.UserService
	TrackService    services.TrackService
	PlaylistService services.PlaylistService
}

// Compile-time check that Handlers satisfies the StrictServerInterface.
var _ gen.StrictServerInterface = (*Handlers)(nil)

// NewHandlers creates a new Handlers instance.
func NewHandlers(
	userService services.UserService,
	trackService services.TrackService,
	playlistService services.PlaylistService,
) *Handlers {
	return &Handlers{
		UserService:     userService,
		TrackService:    trackService,
		PlaylistService: playlistService,
	}
}
