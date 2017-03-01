// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

/*
Package cache contains the interfaces, structs, and logic that form the basis of protonds cache plugin subsystem.

Protond currently implements the following cache plugins:
  - Noop
    - A no operation cache which just returns a hard coded event and discards any new event added to it, this is used for testing protond.
  - Lru
    - An in memory least recently used cache that allows for look backs of arbitrary size, only limited by memory available to the protond application.
*/
package cache
