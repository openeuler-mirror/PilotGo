package plugin_sdk

type PilotGoServer struct {
	Impl PluginInterface
}

func (s *PilotGoServer) OnLoad(args interface{}, resp *OnLoadReply) error {
	resp.Err = s.Impl.OnLoad()
	return nil
}

func (s *PilotGoServer) GetManifest(args interface{}, resp *GetManifestReply) error {
	resp.Manifest = s.Impl.GetManifest()
	return nil
}

func (s *PilotGoServer) GetConfiguration(args interface{}, resp *GetConfigurationReply) error {
	resp.Configuration = s.Impl.GetConfiguration()
	return nil
}

func (s *PilotGoServer) GetWebExtension(args interface{},resp *GetWebExtensionReply) error {
	resp.Extensions = s.Impl.GetWebExtension()
	return nil
}

func (s *PilotGoServer) OnClose(args interface{}, resp *OnCloseReply) error {
	resp.Err = s.Impl.OnClose()
	return nil
}

