package stateAndJid

import (
	"os/exec"
	"strings"
)

// stateAndJid accepts a UUID as input and then returns the state as either up
// or down, depending on if the jid is non-empty. It then also returns either
// a "-" if the jid is empty or the actual jid of the jail.
func Args(uuid string) (string, string) {
	iocuuid := "ioc-" + uuid
	state := "down"
	jid, _ := exec.Command("/usr/sbin/jls", "-j", iocuuid, "jid").Output() // TODO: Use C lib
	jidstr := strings.TrimSpace(string(jid))
	if jidstr != "" {
		state = "up"
		return state, jidstr
	}
	return state, "-"
}
