package main

type Config struct {
	Main    ConfMain
	Session ConfSession
}

// ConfMain - basic configuration
type ConfMain struct {
	Port string
	Name string
}

// ConfSession parameters
type ConfSession struct {
	Secure string
	Name   string
}
