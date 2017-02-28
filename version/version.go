// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

// Package version exposes the current version of protond in symantic format.
package version

const (
	// VERSION is the current version of the protond application
	VERSION = "0.1.0"
)

// GetVersion returns the current version of the protond application
func GetVersion() string {
	return VERSION
}
