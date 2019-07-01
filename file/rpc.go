package file

import "os"

// RPC ..
type RPC struct{}

// Ls ...
func (r *RPC) Ls(path string) []*FileInfo {
	if s, err := os.Stat(path); err == nil {
		if s.IsDir() {
			if items, err := ListingDirectory(path); err == nil {
				return items
			} else {
				panic(err)
			}
		}
	} else {
		panic(err)
	}
	return nil
}

// Cat ...
func (r *RPC) Cat(path string) []byte {
	return nil
}

// Mv ...
func (r *RPC) Mv(path string, dst string) string {
	return "OK"
}

// Cp ...
func (r *RPC) Cp(path string, dst string) string {
	return "OK"
}

// Rm ...
func (r *RPC) Rm(path ...string) string {
	return "OK"
}
