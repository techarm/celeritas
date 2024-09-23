package filesystems

import "time"

// FS is the interface for file systems. In order to satify the interface,
// all of its functions must exist
type FS interface {
	Put(fileName, filder string) error
	Get(destination string, items ...string) error
	List(prefix string) ([]Listing, error)
	Delete(itemsToDelete []string) bool
}

// Listing describes one file on a remote file system
type Listing struct {
	Etag         string
	LastModified time.Time
	Key          string
	Size         float64
	IsDir        bool
}
