// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

/*
Package output contains the interfaces, structs, and logic that form the basis of protonds output plugin subsystem.

Protond currently implements the following output plugins:
  - Stdout
    - This plugin writes to stdout and is used for testing filters and other pieces of functionality of protond.
  - TCP
    - This plugin allows connecting to an arbitrary tcp server, and pushes events over the connection.
*/
package output
