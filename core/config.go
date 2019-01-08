package core

var VERSION = "0.0.1"
var HOMEDIR string
var CONFDIR string
var DATADIR string
var PLUGDIR string
var LOGDIR string

type PluginType uint8

const (
	CONDUIT PluginType = iota
	SOURCE
	SINK
)
