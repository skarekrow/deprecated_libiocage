package templateDatasets

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// templateDatasets returns all template datasets in a slice of a slice.
func Args(Pool *string) []string {
	var temp []string
	match, _ := regexp.Compile(*Pool + "/ioc/templates/[[:alnum:]]*")

	out, _ := exec.Command("/sbin/zfs", "list", "-Hrd", "1",
		*Pool+"/ioc/templates").Output() // TODO: Use C Lib

	outs := strings.Split(string(out), "\n")

	for _, i := range outs {
		om := match.FindAllString(i, -1)
		if om != nil {
			// temp looks like POOL/ioc/templates/NAME
			temp = append(temp, fmt.Sprintf("%s", strings.TrimSpace(om[0])))
		}
	}
	return temp
}
