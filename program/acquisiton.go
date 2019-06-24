package program

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// SoftwareAcquisition ...
type SoftwareAcquisition interface {
	GetSoftwarePath() (string, error)
}

// LocalAcquisition ...
type LocalAcquisition struct{}

// GetSoftwarePath ...
func (l *LocalAcquisition) GetSoftwarePath(software string) (string, error) {
	location, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	location = filepath.Join(location, ".dm", "bin", software)
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return "", err
	}
	return location, nil
}
