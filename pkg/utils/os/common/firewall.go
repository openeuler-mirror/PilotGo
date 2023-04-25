package common

type FireWalldConfig struct {
	Status      string
	Nic         string
	DefaultZone string
	Zones       []string
	Services    []string
}

func (firewalldConfig *FireWalldConfig) Set() {
	firewalldConfig.Status = ""
	firewalldConfig.Nic = ""
	firewalldConfig.DefaultZone = ""
	firewalldConfig.Zones = nil
	firewalldConfig.Services = nil
}
