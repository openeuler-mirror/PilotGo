package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewVersionCommand() *cobra.Command {

	versionCmd := &cobra.Command{
		Example: `
		# Print  the pilotgo-server version
		pilotgo-server version
		`,
		Use:   "version",
		Short: "Print the version of pilotgo server",
		Run: func(cmd *cobra.Command, args []string) {
			klog.InfoS("version info", "Version", GetVersion())
		},
	}
	return versionCmd
}

var (
	version        = "99.99.99"
	buildDate      = "1970-01-01T00:00:00Z"
	gitCommit      = ""
	gitTag         = ""
	gitTreeState   = ""
	extraBuildInfo = ""
)

type Version struct {
	Version        string
	BuildDate      string
	GitCommit      string
	GitTag         string
	GitTreeState   string
	GoVersion      string
	Compiler       string
	Platform       string
	ExtraBuildInfo string
}

func GetVersion() Version {
	var versionStr string

	if gitCommit != "" && gitTag != "" && gitTreeState == "clean" {
		versionStr = gitTag
	} else {
		versionStr = version
		if len(gitCommit) >= 7 {
			versionStr += "+" + gitCommit[0:7]
			if gitTreeState != "clean" {
				versionStr += ".dirty"
			}
		} else {
			versionStr += "+unknown"
		}
	}

	return Version{
		Version:        versionStr,
		BuildDate:      buildDate,
		GitCommit:      gitCommit,
		GitTag:         gitTag,
		GitTreeState:   gitTreeState,
		GoVersion:      runtime.Version(),
		Compiler:       runtime.Compiler,
		Platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		ExtraBuildInfo: extraBuildInfo,
	}
}
