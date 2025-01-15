package script

import "gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"

type DangerousCommands dao.DangerousCommands

type CommandsWithKey struct {
	ID      uint   `json:"id"`
	Key     uint   `json:"key"`
	Command string `json:"command"`
	Active  bool   `json:"active"`
}

func CreateDangerousCommands() error {
	return dao.CreateDangerousCommands()
}

func UpdateCommandsBlackList(_whitelist []uint) error {
	return dao.UpdateCommandsBlackList(_whitelist)
}

func GetDangerousCommandsList() ([]*CommandsWithKey, error) {
	commands, err := dao.GetDangerousCommandsList()
	if err != nil {
		return nil, err
	}

	_commands := []*CommandsWithKey{}
	for _, c := range commands {
		_c := &CommandsWithKey{
			ID:      c.ID,
			Key:     c.ID,
			Command: c.Command,
			Active:  c.Active,
		}
		_commands = append(_commands, _c)
	}
	return _commands, nil
}

func GetDangerousCommandsInBlackList() ([]string, error) {
	commands, err := dao.GetDangerousCommandsList()
	if err != nil {
		return nil, err
	}

	_commands := make([]string, 0)
	for _, c := range commands {
		if c.Active {
			_commands = append(_commands, c.Command)
		}
	}
	return _commands, nil
}
