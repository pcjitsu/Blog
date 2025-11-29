module Blog // The name of your module

go 1.25.4 // The Go version this module is intended for

require (
	// External libraries your code uses
	github.com/google/uuid v1.6.0 // Used for generating unique user IDs
	github.com/lib/pq v1.10.9 // The PostgreSQL driver for Go
)
