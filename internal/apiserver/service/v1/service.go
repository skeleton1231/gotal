package service

import "github.com/skeleton1231/gotal/internal/apiserver/store"

// Service is the interface that abstracts the functionalities of your services.
type Service interface {
	Users() UserSrv // Users returns an instance of UserSrv which handles user-related operations.
}

// service is a struct that implements the Service interface.
// It holds a reference to the store factory to access different stores.
type service struct {
	store store.Factory
}

// NewService is a constructor function for creating a new instance of service.
// It takes a store.Factory as an argument and returns a Service.
func NewService(store store.Factory) Service {
	return &service{
		store: store, // Initializing the store field with the provided store factory.
	}
}

// Users is a method on service struct that returns a new instance of UserSrv.
func (s *service) Users() UserSrv {
	return newUsers(s) // Creating a new UserSrv using the current service instance.
}
