package copyFile

import (
	"io/ioutil"
	"log"
)

// Reads the source file and then writes it into the destination. 644 perms.
func Args(src string, dst string) {
	f, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dst, f, 0644)

	if err != nil {
		log.Fatal(err)
	}
}
