package main

import "github.com/caiomarcatti12/poc-go-plugin/pluginapi"

// impl concreta do contrato
type impl struct{}

func (impl) Greet(name string) string {
	return "Olá, " + name + " 👋"
}

// Símbolos exportados (NOMES DEVEM BATER com o Lookup do host)
var (
	Plugin pluginapi.Greeter = impl{}
	ABI                      = pluginapi.ABI
	Info                     = pluginapi.Info{
		Name:        "greeter",
		Version:     "v1.0.0",
		Description: "Plugin de exemplo que dá oi",
		ABI:         pluginapi.ABI,
	}
)
