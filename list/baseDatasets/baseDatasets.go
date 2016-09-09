package baseDatasets

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// bseDatasets returns all base datasets in a slice of a slice.
func Args(Pool *string) []string {
	var base []string
	match, _ := regexp.Compile(*Pool + "/ioc/bases/[[:digit:]]{1,2}." +
		"[[:digit:]]-RELEASE|[[:digit:]]{1,2}.[[:digit:]]-[[:alnum:]]{0,9}")

	out, _ := exec.Command("/sbin/zfs", "list", "-Hrd", "1",
		*Pool+"/ioc/bases").Output() // TODO: Use C Lib

	outs := strings.Split(string(out), "\n")

	for _, i := range outs {
		om := match.FindAllString(i, -1)
		if om != nil {
			fmt.Println(om)
			// base looks like POOL/ioc/bases/10.3-RELEASE
			base = append(base, fmt.Sprintf("%s", strings.TrimSpace(om[0])))
		}
	}
	return base
}
