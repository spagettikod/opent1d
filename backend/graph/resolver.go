package graph

import "github.com/spagettikod/opent1d/datastore"

//go:generate go run github.com/99designs/gqlgen generate --verbose

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Store datastore.Store
}
