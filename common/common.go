// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package common

import (
	"os"
)

// PathExists determines whether or not the specified path exists on disk.
func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
