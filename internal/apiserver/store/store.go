package store

// client is a package-level variable that holds the instance of Factory.
var client Factory

// Factory is an interface that abstracts the creation of different stores.
// It provides methods to access different data stores and to close them.
type Factory interface {
	Users() UserStore // Users returns an instance of UserStore for user-related data operations.
	Close() error     // Close is responsible for closing any resources used by the factory, e.g., database connections.
}

// Client is a function that returns the current instance of Factory.
func Client() Factory {
	return client // Returning the package-level client variable.
}

// SetClient is a function for setting the package-level client variable.
// It is used to initialize the client with a concrete implementation of Factory.
func SetClient(factory Factory) {
	client = factory // Assigning the provided factory to the package-level client variable.
}
