package configfile

import "gitee.com/openeuler/PilotGo/app/server/service/internal/dao"

type ConfigFile = dao.ConfigFile

func AddConfigFile(cf ConfigFile) error {
	return dao.AddConfigFile(cf)
}
