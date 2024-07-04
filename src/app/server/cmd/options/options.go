package options

type ServerOptions struct {
	ConfigFile string
}

func NewServerOptions() *ServerOptions {
	s := &ServerOptions{
		ConfigFile: "./config_server.yaml",
	}
	return s
}
