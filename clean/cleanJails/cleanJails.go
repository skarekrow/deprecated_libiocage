package cleanJails

import (
	"fmt"
	"github.com/iocage/iocage-cli/src/cmd/destroy"
	"github.com/iocage/libiocage/askUser"
	"os"
	"strings"
)

func Args(Pool *string, force bool, datasets []string) {
	var q bool

	if !force {
		qstr := fmt.Sprintf(strings.Repeat("-", 9) + "\nWARNING: \n" +
			strings.Repeat("-", 9) + "\n\n" + "This will destroy all jails!\n" +
			"Are you sure? y[N]: ")
		q, _ = askUser.Args(qstr, true)
	}

	if !q && !force {
		fmt.Printf("Command not confirmed.  No action taken.\n")
		os.Exit(0)
	} else {
		destroy.Args(Pool, true, true, datasets)
	}
}
