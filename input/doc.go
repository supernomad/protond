// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

/*
Package input contains the interfaces, structs, and logic that form the basis of protonds input plugin subsystem.

Protond currently implements the following input plugins:
  - Stdin
    - An input plugin that reads from stdin and is used for testing filters and other pieces of functionality of protond.
  - TCP
    - This plugin allows listening on an arbitrary tcp socket, and reads new line terminated strings from the connected clients.
*/
package input
