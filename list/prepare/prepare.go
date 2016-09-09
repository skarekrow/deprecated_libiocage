package prepare

import (
	"fmt"
	"github.com/iocage/libiocage/get/stateAndJid"
	"github.com/iocage/libiocage/get/uclProp"
	"github.com/iocage/libiocage/get/uuidPathTag"
	// "github.com/fatih/color"
	"os/exec"
	"strconv"
	"strings"
)

// red := color.New(color.FgHiRed).PrintfFunc()
// yellow := color.New(color.FgHiYellow).PrintfFunc()

// prepare accepts a Pool, a slice of datasets, and a boolean that governs
// whether to return a header + table or a string separated by whitespace.
func Args(Pool *string, ds []string, ltype string, header bool) [][]string {
	var d [][]string
	var i int = 1

	for _, da := range ds {
		istr := strconv.Itoa(i)

		switch ltype {
		case "all":
			uuid, path, jtag, _ := uuidPathTag.Args(Pool, []string{da}, da)
			jboot, _ := uclProp.Args(path, "boot", false)
			jstate, jid := stateAndJid.Args(uuid)
			jip4addr := "-"

			if jstate == "up" {
				cmd := fmt.Sprintf("jls -j %s -h ip4.addr|grep -v -x ip4.addr",
					jid)
				out, _ := exec.Command("sh", "-c", cmd).Output()
				jip4addr = strings.TrimSpace(string(out))
			}

			jbase, _ := uclProp.Args(path, "base", false)

			if !header {
				d = append(d, []string{jid, uuid, jboot, jstate, jtag,
					jip4addr, jbase, istr})
			} else {
				d = append(d, []string{fmt.Sprintf("%s %s %s %s %s %s %s",
					jid, strings.TrimSpace(uuid), jboot, jstate, jtag, jip4addr,
					jbase, istr)})
			}
			i++
		case "base":
			da = strings.TrimSpace(strings.Split(da, "/")[3])

			if !header {
				d = append(d, []string{da})
			} else {
				d = append(d, []string{fmt.Sprintf("%s", da)})
			}
		case "template":
			da = strings.TrimSpace(strings.Split(da, "/")[3])
			_, path, jtag, _ := uuidPathTag.Args(Pool, ds, da)

			uuid, _ := uclProp.Args(path, "host_hostuuid", false)
			// jtag, _ := uclProp.Args(path, "tag", false)
			jboot, _ := uclProp.Args(path, "boot", false)
			jstate, jid := stateAndJid.Args(uuid)
			jip4addr := "-"

			if jstate == "up" {
				cmd := fmt.Sprintf("jls -j %s -h ip4.addr|grep -v -x ip4.addr",
					jid)
				out, _ := exec.Command("sh", "-c", cmd).Output()
				jip4addr = strings.TrimSpace(string(out))
			}

			jbase, _ := uclProp.Args(path, "base", false)

			if !header {
				d = append(d, []string{jid, uuid, jboot, jstate, jtag,
					jip4addr, jbase, istr})
			} else {
				d = append(d, []string{fmt.Sprintf("%s %s %s %s %s %s %s",
					jid, strings.TrimSpace(uuid), jboot, jstate, jtag, jip4addr,
					jbase, istr)})
			}
		}
	}
	return d
}
