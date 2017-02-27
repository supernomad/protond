// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

/*
Package worker contains the structs, and logic that form the basis of protonds worker subsystem.

Protond currently implements a single worker type, that is responsible for ingesting events from an arbitrary set of user defined input plugins, processing those events with an arbitrary set of filter plugins, and pushing those filtered events to an arbitrary set of output plugins.
*/
package worker
