package uclProp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/iocage/libiocage/get/stateAndJid"
	"os"
	"os/exec"
	"strings"
)

// uclProp accepts a path where the config file lives and the property.
// A default boolean governs if it's a default property or a jail property.
// It then attempts to return a native jail/ioc property, and if that fails
// attempts to find a ZFS property.
func Args(path string, prop string, def bool) (string, error) {
	var str string

	if def {
		out, _ := exec.Command("/usr/local/bin/uclcmd", "get", "-qf",
			"/ioc/.default", prop).Output() // TODO: Use C Lib, support different mountpoints for /ioc

		str = strings.TrimSpace(string(out))
	} else {
		switch prop {
		case "all":
			config := path + "/config"
			f, err := os.Open(config)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				r := strings.Replace(scanner.Text(), " = ", ":", 1)
				r = strings.Replace(r, ";", "", 1)
				r = strings.Replace(r, "\"", "", -1)
				fmt.Println(r)
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case "state":
			uuid := strings.TrimSpace(strings.Split(path, "/")[3])
			state, _ := stateAndJid.Args(uuid)
			str = state
		default:
			config := path + "/config"
			out, _ := exec.Command("/usr/local/bin/uclcmd", "get", "-qf",
				config, prop).Output() // TODO: Use C Lib

			str = strings.TrimSpace(string(out))

			if str == "null" {
				out, err := exec.Command("/sbin/zfs", "get", "-H", "-o", "value",
					prop, path).Output() // TODO: Use C Lib
				if err != nil {
					serr := fmt.Sprintf("Invalid property: %s!", prop)
					return "", errors.New(serr)
				}

				str = strings.TrimSpace(string(out))
			}
		}
	}

	return str, nil
}
