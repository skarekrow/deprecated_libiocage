package jailDatasets

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// jailDatasets returns all jail datasets in a slice of a slice.
func Args(Pool *string) []string {
	var jail []string
	match, _ := regexp.Compile(*Pool + "/ioc/jails/[[:alnum:]]{8}-" +
		"(([[:alnum:]]{4}-){3})[[:alnum:]]{12}[[:space:]]")

	out, _ := exec.Command("/sbin/zfs", "list", "-Hrd", "1",
		*Pool+"/ioc/jails").Output() // TODO: Use C Lib

	outs := strings.Split(string(out), "\n")

	for _, i := range outs {
		om := match.FindAllString(i, -1)
		if om != nil {
			// jail looks like POOL/ioc/jails/UUID
			jail = append(jail, fmt.Sprintf("%s", strings.TrimSpace(om[0])))
		}
	}
	return jail
}
