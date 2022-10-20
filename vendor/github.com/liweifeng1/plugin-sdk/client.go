package plugin_sdk

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type PluginManifest struct {
	Id string
	Name string
	Author string
	Version string
}

type ExtensionType int

const (
	CSS = iota+1
	JavaScript
	HTML
)

type WebExtension struct {
	Type ExtensionType
	PathMatchRegex string
	Source string
}

const (
	UrlValue = iota
	PortValue
	ProtocolValue
)

type PluginConfig struct {
	Title string
	Description string
	Key string
	Type int
	Values string
}

type PluginInterface interface {
	OnLoad() error
	OnClose() error
	GetManifest() PluginManifest
	GetConfiguration() []PluginConfig
	GetWebExtension() []WebExtension
}

type PilotGoRPC struct {
	client *rpc.Client
}

type OnLoadReply struct {
	Err error
}


func(p *PilotGoRPC) OnLoad() error {
	rep:=OnLoadReply{}
	err:=p.client.Call("Plugin.OnLoad",new(interface{}),&rep)
	if err!=nil {
		panic(err)
	}
	return rep.Err
}

type GetManifestReply struct {
	Manifest PluginManifest
}

func(p *PilotGoRPC) GetManifest() PluginManifest {
	rep:=GetManifestReply{}
	err:=p.client.Call("Plugin.GetManifest",new(interface{}),&rep)
	if err!=nil {
		panic(err)
	}
	return rep.Manifest
}


type GetConfigurationReply struct {
	Configuration []PluginConfig
}

func (p *PilotGoRPC) GetConfiguration() []PluginConfig {
	rep := &GetConfigurationReply{}
	err := p.client.Call("Plugin.GetConfiguration", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Configuration
}

type GetWebExtensionReply struct {
	Extensions []WebExtension
}

func (p *PilotGoRPC) GetWebExtension() []WebExtension {
	rep := &GetWebExtensionReply{}
	err := p.client.Call("Plugin.GetWebExtension", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Extensions
}


type OnCloseReply struct {
	Err error
}

func(p *PilotGoRPC) OnClose() error {
	rep:=OnCloseReply{}
	err:=p.client.Call("Plugin.OnClose",new(interface{}),&rep)
	if err!=nil {
		panic(err)
	}
	return rep.Err
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion: 1,
	MagicCookieKey: "PILOTGO_PLUGIN",
	MagicCookieValue: "Mz1K0OGpIRs",
}

type PilotGoPlugin struct {
	Impl PluginInterface
}

func (p *PilotGoPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PilotGoServer{Impl: p.Impl}, nil
}

func (PilotGoPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PilotGoRPC{client: c}, nil
}