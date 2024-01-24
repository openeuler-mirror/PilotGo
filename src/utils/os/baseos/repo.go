package baseos

import (
	"fmt"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils"
)

// 获取repo文件
func (b *BaseOS) GetRepoConfig() (string, error) {
	result := `[`
	exitc, filename, stde, err := utils.RunCommand("find " + b.RepoPath + " -type f -name \"*.repo\"")
	if exitc == 0 && len(filename) > 0 && stde == "" && err == nil {
		for _, v := range strings.Split(filename, "\n") {
			name := strings.Split(v, "/")[len(strings.Split(v, "/"))-1]
			context, err := utils.FileReadString(v)
			if err != nil {
				return "", err
			}
			result = result + `{"path": ` + b.RepoPath + `, "name":` + name + `, "file": ` + context + `}`
		}
	} else {
		logger.Error("failed to get rpoe file context: %d, %s, %s, %v", exitc, filename, stde, err)
		return result, fmt.Errorf("failed to get rpoe file context: %d, %s, %s, %v", exitc, filename, stde, err)
	}
	result = result + `]`
	return result, nil
}
