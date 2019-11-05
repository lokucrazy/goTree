package goTree

import "errors"

var parseError error = errors.New("cannot parse empty text")
var matchError error = errors.New("Could not parse/find commit fields")
