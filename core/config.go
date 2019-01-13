package core

import "os"

var (
	VERSION string
	HOMEDIR = getDefaultEnv("SLUUS_HOMEDIR", DEFAULT_HOME)
	CONFDIR = getDefaultEnv("SLUUS_CONFDIR", DEFAULT_CONF)
	DATADIR = getDefaultEnv("SLUUS_DATADIR", DEFAULT_DATA)
	PLUGDIR = HOMEDIR + string(os.PathSeparator) + "plugin"
	LOGDIR  = HOMEDIR + string(os.PathSeparator) + "log"
)

func getDefaultEnv(key string, builtin string) (value string) {
	if value, ok := os.LookupEnv(key); !ok {
		return builtin
	} else {
		return value
	}
}
