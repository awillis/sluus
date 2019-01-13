package core

import "os"

var (
	VERSION string
	HOMEDIR = os.Getenv("SLUUS_HOMEDIR")
	CONFDIR = os.Getenv("SLUUS_CONFDIR")
	DATADIR = os.Getenv("SLUUS_DATADIR")
	PLUGDIR = os.Getenv("SLUUS_HOMEDIR") + string(os.PathSeparator) + "plugin"
	LOGDIR  = os.Getenv("SLUUS_HOMEDIR") + string(os.PathSeparator) + "log"
)
