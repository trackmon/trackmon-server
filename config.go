package main

type Configuration struct {
	ListeningAddress      string
	DatabaseAddress       string
	DatabasePassword      string
	AutoUpdateChecker     bool
	LogFileLocationPrefix string
}
