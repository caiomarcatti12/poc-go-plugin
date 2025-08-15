package pluginapi

const ABI = 1

type Greeter interface {
	Greet(name string) string
}

type Info struct {
	Name        string
	Version     string
	Description string
	ABI         int
}
