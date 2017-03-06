// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

/*
Package alert contains the interfaces, structs, and logic that form the basis of protonds alert plugin subsystem.

Protond currently implements the following alert plugins:
  - Noop
    - A no alert which just noops the event emission.
  - Http
    - An http alert that allows for emitting events to the specified endpoint for alerting.
*/
package alert
