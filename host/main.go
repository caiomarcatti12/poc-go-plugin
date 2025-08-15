package main

import (
	"flag"
	"fmt"
	"log"
	"plugin"

	"github.com/caiomarcatti12/poc-go-plugin/pluginapi"
)

func main() {
	var soPath string
	flag.StringVar(&soPath, "plugin", "../bin/greeter_v1.so", "caminho do arquivo .so do plugin")
	flag.Parse()

	p, err := plugin.Open(soPath)
	if err != nil {
		log.Fatalf("falha ao abrir plugin %q: %v", soPath, err)
	}

	// 1) Checa ABI
	abiSym, err := p.Lookup("ABI")
	if err != nil {
		log.Fatalf("símbolo ABI não encontrado no plugin: %v", err)
	}
	abi, ok := abiSym.(*int)
	if !ok || *abi != pluginapi.ABI {
		log.Fatalf("ABI incompatível: plugin=%v host=%v", valueOr(abiSym), pluginapi.ABI)
	}

	// 2) (Opcional) Lê Info
	if infoSym, err := p.Lookup("Info"); err == nil {
		if info, ok := infoSym.(*pluginapi.Info); ok {
			log.Printf("Plugin carregado: %s %s (ABI=%d) — %s", info.Name, info.Version, info.ABI, info.Description)
		}
	}

	// 3) Obtém a implementação
	sym, err := p.Lookup("Plugin")
	if err != nil {
		log.Fatalf("símbolo Plugin não encontrado: %v", err)
	}

	// O plugin exporta uma variável; Lookup retorna um ponteiro para ela.
	var greeter pluginapi.Greeter
	switch v := sym.(type) {
	case *pluginapi.Greeter:
		if v == nil {
			log.Fatalf("símbolo Plugin é *pluginapi.Greeter, mas está nil")
		}
		greeter = *v
	case pluginapi.Greeter:
		greeter = v
	default:
		log.Fatalf("símbolo Plugin não implementa pluginapi.Greeter (tipo real: %T)", sym)
	}

	fmt.Println(greeter.Greet("Mundo"))
}

func valueOr(v any) any {
	switch x := v.(type) {
	case *int:
		if x == nil {
			return nil
		}
		return *x
	default:
		return v
	}
}
