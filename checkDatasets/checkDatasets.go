package checkDatasets

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"os/user"
)

func Args(Pool, Iocroot *string) {
	checkds(Pool, Iocroot, "ioc")
	checkds(Pool, Iocroot, "ioc/download")
	checkds(Pool, Iocroot, "ioc/jails")
	checkds(Pool, Iocroot, "ioc/bases")
	checkds(Pool, Iocroot, "ioc/templates")
}

func checkds(Pool, Iocroot *string, dataset string) {
	u, _ := user.Current()
	red := color.New(color.FgHiRed)
	yellow := color.New(color.FgHiYellow)

	out, _ := exec.Command("/sbin/zfs", "get", "-H", "creation",
		*Pool+"/"+dataset).Output() // TODO: Use C Lib

	str := string(out)
	if str == "" {
		if u.Uid != "0" {
			fmt.Println("Please run as root to create missing datasets")
			os.Exit(1)
		}

		exec.Command("/sbin/zfs", "create", "-p",
			*Pool+"/"+dataset).Run() // TODO: Use C Lib
		if dataset == "ioc" {
			mountpoint := "mountpoint=" + *Iocroot
			exec.Command("/sbin/zfs", "set", mountpoint,
				*Pool+"/"+dataset).Run()
		}

		red.Printf("Dataset: ")
		yellow.Printf(dataset)
		red.Printf(" is missing, creating it now.\n")
	}
}
