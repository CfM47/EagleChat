package repositories

import "eaglechat/apps/client/internal/middleware/domain/entities"

// IDManagerRepository defines the contract for storing and retrieving data
// about discovered ID Manager services. Implementations of this interface are
// expected to be thread-safe.
type IDManagerRepository interface {
	// Add adds a new ID Manager's data to the repository or updates the
	// existing entry if an entry with the same ID already exists. This method
	// should also update the internal 'last seen' timestamp for the entry.
	Add(data entities.IDManagerData)

	// GetAll retrieves the data for all active ID Managers currently held
	// in the repository.
	GetAll() []entities.IDManagerData
}
