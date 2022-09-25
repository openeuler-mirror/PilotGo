package model

type PluginLists struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
	Url    string `json:"url"`
}

type LoadPlugin struct {
	Name string `json:"Name"`
}

type UnLoadPlugin struct {
	Name string `json:"Name"`
}
