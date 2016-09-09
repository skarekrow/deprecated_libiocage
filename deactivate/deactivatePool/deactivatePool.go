package deactivatePool

import (
	"errors"
	"fmt"
	"os/exec"
)

func Args(pool string) (string, error) {
	_, err := exec.Command("/sbin/zfs", "set",
		"org.freebsd.ioc:active=no", pool).Output()
	str := fmt.Sprintf("%s successfully deactivated.\n", pool)

	if err != nil {
		errstr := fmt.Sprintf("Pool %s can not be found!", pool)
		return "", errors.New(errstr)
	}

	return str, nil
}
