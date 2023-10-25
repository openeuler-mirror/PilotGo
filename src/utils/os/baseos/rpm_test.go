package baseos

import (
	"strings"
	"testing"

	"gitee.com/openeuler/PilotGo/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAllRpm(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetAllRpm()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}

func TestGetRpmSource(t *testing.T) {
	var osobj BaseOS
	rpm := "time"
	tmp, err := osobj.GetRpmSource(rpm)
	assert.Nil(t, err)
	assert.Equal(t, rpm, strings.Split(tmp[0].Name, "-")[0])
}

func TestGetRpmInfo(t *testing.T) {
	var osobj BaseOS
	rpm := "time"
	tmp, err := osobj.GetRpmInfo(rpm)
	assert.Nil(t, err)
	assert.Equal(t, rpm, tmp.Name)
}

func TestInstallAndRemoveRpm(t *testing.T) {
	var osobj BaseOS
	rpm := "bind"

	exitc, stdo, stde, err := utils.RunCommand("rpm -qi " + rpm)
	if exitc == 0 && len(stdo) > 0 && stde == "" && err == nil {
		t.Run("test remove rpm", func(t *testing.T) {
			assert.Nil(t, osobj.RemoveRpm(rpm))
			_, err := osobj.GetRpmInfo(rpm)
			assert.NotNil(t, err)
		})

		t.Run("test install rpm", func(t *testing.T) {
			assert.Nil(t, osobj.InstallRpm(rpm))
			tmp, err := osobj.GetRpmInfo(rpm)
			assert.Nil(t, err)
			assert.Equal(t, rpm, tmp.Name)
		})
	} else if exitc == 1 && strings.Replace(stdo, "\n", "", -1) == "package bind is not installed" && stde == "" && err == nil {
		t.Run("test install rpm", func(t *testing.T) {
			assert.Nil(t, osobj.InstallRpm(rpm))
			tmp, err := osobj.GetRpmInfo(rpm)
			assert.Nil(t, err)
			assert.Equal(t, rpm, tmp.Name)
		})

		t.Run("test remove rpm", func(t *testing.T) {
			assert.Nil(t, osobj.RemoveRpm(rpm))
			_, err := osobj.GetRpmInfo(rpm)
			assert.NotNil(t, err)
		})
	} else if exitc == 127 && stdo == "" && strings.Contains(stde, "command not found") == true && err == nil {
		t.Errorf("[TestInstallAndRemoveRpm]rpm command not found: %d, %s, %s, %v\n", exitc, stdo, stde, err)
	} else {
		t.Errorf("[TestInstallAndRemoveRpm]other error: %d, %s, %s, %v\n", exitc, stdo, stde, err)
	}
}

func TestYumsource(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetRepoSource()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}

func TestParseRepoContent(t *testing.T) {
	s1 := `[OS]
	name=OS
	baseurl=http://repo.openeuler.org/openEuler-22.03-LTS-SP2/OS/$basearch/
	metalink=https://mirrors.openeuler.org/metalink?repo=$releasever/OS&arch=$basearch
	metadata_expire=1h
	enabled=1
	gpgcheck=1
	gpgkey=http://repo.openeuler.org/openEuler-22.03-LTS-SP2/OS/$basearch/RPM-GPG-KEY-openEuler

	[baseos]
	name=CentOS Stream $releasever - BaseOS
	mirrorlist=http://mirrorlist.centos.org/?release=$stream&arch=$basearch&repo=BaseOS&infra=$infra
	#baseurl=http://mirror.centos.org/$contentdir/$stream/BaseOS/$basearch/os/
	gpgcheck=1
	enabled=1
	gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-centosofficial`

	repos, err := parseRepoContent(s1)
	assert.NoError(t, err)
	assert.Len(t, repos, 2)

	assert.NotEmpty(t, repos[0].Name)
	assert.NotEmpty(t, repos[0].BaseURL)
	assert.NotEmpty(t, repos[0].MetaLink)
	assert.NotEmpty(t, repos[0].MetadataExpire)
	assert.NotEmpty(t, repos[0].Enabled)
	assert.NotEmpty(t, repos[0].GPGCheck)
	assert.NotEmpty(t, repos[0].GPGKey)
	assert.Empty(t, repos[0].MirrorList)

	assert.NotEmpty(t, repos[1].Name)
	assert.NotEmpty(t, repos[1].MirrorList)
	assert.NotEmpty(t, repos[1].GPGCheck)
	assert.NotEmpty(t, repos[1].Enabled)
	assert.NotEmpty(t, repos[1].GPGKey)
	assert.Empty(t, repos[1].MetadataExpire)
	assert.Empty(t, repos[1].MetaLink)
	assert.Empty(t, repos[1].BaseURL)
}
