package domain

// Project represents an interface for managing Java classes within a project.
type Project interface {
	// Classes retrieves all Java classes within the project.
	Classes() ([]Class, error)
}
