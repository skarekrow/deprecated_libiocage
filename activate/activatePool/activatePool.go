package activatePool

import (
	"errors"
	"fmt"
	"os/exec"
)

func Args(pool string) (string, error) {
	_, err := exec.Command("/sbin/zfs", "set",
		"org.freebsd.ioc:active=yes", pool).Output()
	str := fmt.Sprintf("%s successfully activated.\n", pool)

	if err != nil {
		errstr := fmt.Sprintf("Pool %s can not be found!", pool)
		return "", errors.New(errstr)
	}

	return str, nil
}
