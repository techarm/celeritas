package sftpfs

import "github.com/techarm/celeritas/filesystems"

type SFTP struct {
	Host string
	User string
	Pass string
	Port string
}

func (s *SFTP) Put(fileName, folder string) error {
	return nil
}

func (s *SFTP) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing
	return listing, nil
}

func (s *SFTP) Delete(itemsToDelete []string) bool {
	return true
}

func (s *SFTP) Get(destination string, item ...string) error {
	return nil
}
