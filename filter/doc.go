// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

/*
Package filter contains the interfaces, structs, and logic that form the basis of protonds filter plugin subsystem.

Protond currently implements the following filter plugins:
  - Noop
    - A no operation filter which just returns the event unchanged, this is used for pass through protond relays and testing protond.
  - Javascript
    - This plugin allows for arbitrary javascript scripts that can modify and call certain functions on all events.
*/
package filter
