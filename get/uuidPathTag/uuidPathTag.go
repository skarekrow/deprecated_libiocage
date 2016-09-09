package uuidPathTag

import (
	"errors"
	"fmt"
	"github.com/iocage/libiocage/get/uclProp"
	"strings"
)

// uuidPathTag accepts the zpool, a slice containing all the datasets,
// and what we are trying to match. It then returns the UUID of the match,
// the path to the jail, and the tag. Along with an error if the jail cannot
// be found.
func Args(Pool *string, datasets []string, dataset string) (string,
	string, string, error) {
	for _, d := range datasets {
		path := strings.TrimSpace(strings.Split(d, *Pool)[1])
		uuid := strings.TrimSpace(strings.Split(d, "/")[3])
		tag, _ := uclProp.Args(path, "tag", false)

		if tag == dataset {
			return uuid, path, tag, nil
		} else if len(datasets) == 1 {
			return uuid, path, tag, nil
		}
	}

	err := fmt.Sprintf("%s was not found!", dataset)
	return "", "", "", errors.New(err)
}
