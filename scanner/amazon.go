package scanner

import (
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/future-architect/vuls/config"
	"github.com/future-architect/vuls/logging"
	"github.com/future-architect/vuls/models"
)

// inherit OsTypeInterface
type amazon struct {
	redhatBase
}

// NewAmazon is constructor
func newAmazon(c config.ServerInfo) *amazon {
	r := &amazon{
		redhatBase{
			base: base{
				osPackages: osPackages{
					Packages:  models.Packages{},
					VulnInfos: models.VulnInfos{},
				},
			},
			sudo: rootPrivAmazon{},
		},
	}
	r.log = logging.NewNormalLogger()
	r.setServerInfo(c)
	return r
}

func (o *amazon) checkScanMode() error {
	return nil
}

func (o *amazon) checkDeps() error {
	if o.getServerInfo().Mode.IsFast() {
		return o.execCheckDeps(o.depsFast())
	}
	if o.getServerInfo().Mode.IsFastRoot() {
		return o.execCheckDeps(o.depsFastRoot())
	}
	if o.getServerInfo().Mode.IsDeep() {
		return o.execCheckDeps(o.depsDeep())
	}
	return xerrors.New("Unknown scan mode")
}

func (o *amazon) depsFast() []string {
	if o.getServerInfo().Mode.IsOffline() {
		return []string{}
	}
	// repoquery
	switch s := strings.Fields(o.getDistro().Release)[0]; s {
	case "1", "2":
		return []string{"yum-utils"}
	default:
		if _, err := time.Parse("2006.01", s); err == nil {
			return []string{"yum-utils"}
		}
		return []string{"dnf-utils"}
	}
}

func (o *amazon) depsFastRoot() []string {
	switch s := strings.Fields(o.getDistro().Release)[0]; s {
	case "1", "2":
		return []string{"yum-utils"}
	default:
		if _, err := time.Parse("2006.01", s); err == nil {
			return []string{"yum-utils"}
		}
		return []string{"dnf-utils"}
	}
}

func (o *amazon) depsDeep() []string {
	return o.depsFastRoot()
}

func (o *amazon) checkIfSudoNoPasswd() error {
	if o.getServerInfo().Mode.IsFast() {
		return o.execCheckIfSudoNoPasswd(o.sudoNoPasswdCmdsFast())
	}
	if o.getServerInfo().Mode.IsFastRoot() {
		return o.execCheckIfSudoNoPasswd(o.sudoNoPasswdCmdsFastRoot())
	}
	return o.execCheckIfSudoNoPasswd(o.sudoNoPasswdCmdsDeep())
}

func (o *amazon) sudoNoPasswdCmdsFast() []cmd {
	return []cmd{}
}

func (o *amazon) sudoNoPasswdCmdsFastRoot() []cmd {
	return []cmd{
		{"needs-restarting", exitStatusZero},
		{"which which", exitStatusZero},
		{"stat /proc/1/exe", exitStatusZero},
		{"ls -l /proc/1/exe", exitStatusZero},
		{"cat /proc/1/maps", exitStatusZero},
		{"lsof -i -P -n", exitStatusZero},
	}
}

func (o *amazon) sudoNoPasswdCmdsDeep() []cmd {
	return o.sudoNoPasswdCmdsFastRoot()
}

type rootPrivAmazon struct{}

func (o rootPrivAmazon) repoquery() bool {
	return false
}

func (o rootPrivAmazon) yumMakeCache() bool {
	return false
}

func (o rootPrivAmazon) yumPS() bool {
	return false
}
