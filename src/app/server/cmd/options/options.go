package options

type ServerOptions struct {
	Config string
}

func NewServerOptions() *ServerOptions {
	s := &ServerOptions{
		Config: "./config_server.yaml",
	}
	return s
}
