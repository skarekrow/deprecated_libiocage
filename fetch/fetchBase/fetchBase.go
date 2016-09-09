package fetchBase

// I now need to extract these to /download.
import (
	"fmt"
	"github.com/iocage/libiocage/askUser"
	"github.com/iocage/libiocage/get/uclProp"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/secsy/goftp"
)

func Args(Pool, Iocroot *string, props []string) {
	var r bool
	var base string
	var bases []string

	u, _ := exec.Command("/usr/bin/uname", "-r").Output()
	uname := strings.TrimSpace(string(u))

	// COLORS!
	cyan := color.New(color.FgHiCyan)
	cyanul := cyan.Add(color.Underline)
	green := color.New(color.FgHiGreen)
	greenul := green.Add(color.Underline)
	red := color.New(color.FgHiRed)
	redsp := red.SprintfFunc()
	white := color.New(color.FgHiWhite)
	whitesp := white.SprintfFunc()
	yellow := color.New(color.FgHiYellow)
	yellowsp := yellow.SprintfFunc()

	if len(props) < 0 {
		if !strings.Contains(props[0], "=") {
			fmt.Printf("Property %s has invalid syntax!\n", props[0])
			os.Exit(1)
		}
	}

	// A quick check to see if they supplied us with any properties we care
	// about..
	m := make(map[string]string)
	for _, v := range props {
		p := strings.Split(v, "=")
		prop, val := p[0], p[1]
		m[prop] = val
	}

	// Check if any of these properties exist, as the users choice takes
	// preference.
	if m["base"] != "" {
		r = true
		base = m["base"]
	}

	// We are going to assume the root of their ftp host, they can supply
	// any different root.
	if m["ftphost"] != "" && m["ftpdir"] == "" {
		m["ftpdir"] = "/"
	}

	if m["ftphost"] == "" {
		m["ftphost"], _ = uclProp.Args(*Iocroot, "ftphost", true)
	}

	if m["ftpdir"] == "" {
		m["ftpdir"], _ = uclProp.Args(*Iocroot, "ftpdir", true)
	}

	if m["ftpfiles"] == "" {
		m["ftpfiles"], _ = uclProp.Args(*Iocroot, "ftpfiles", true)
	}

	// Create client object with default config
	client, err := goftp.Dial(m["ftphost"])
	defer client.Close()
	if err != nil {
		panic(err)
	}

	if !r {
		files, _ := client.ReadDir(m["ftpdir"])
		match, _ := regexp.Compile("[[:alnum:]]{1,2}.[[:alnum:]]-RELEASE")
		cyanul.Printf("Bases available from ")
		greenul.Printf(m["ftphost"])
		cyanul.Printf(":\n")
		for _, f := range files {
			// fstr := fmt.Sprintf("%v", f.Name())
			findmatch := match.FindAllString(f.Name(), -1)
			if findmatch != nil {
				bases = append(bases, findmatch[0])
				yellow.Printf("%s\n", findmatch[0])
			}
		}
		qstr := whitesp("Please select a base [%s]", yellowsp(uname)) +
			whitesp(" or type ") + redsp("EXIT") + whitesp(": ")
		_, base = askUser.Args(qstr, false)

		if base == "" {
			base = uname
		}

		if base == "EXIT" {
			os.Exit(0)
		}

		exists := baseExists(base, bases)

		if !exists {
			for {
				qstr := whitesp("Please select a valid base or type ") +
					redsp("EXIT") + whitesp(": ")
				_, base = askUser.Args(qstr, false)

				if base == "" {
					base = uname
				}

				if base == "EXIT" {
					os.Exit(0)
				}

				exists = baseExists(base, bases)
				if exists {
					break
				}
			}
		}
	}

	dlpath := *Iocroot + "/download/" + base
	basepath := *Iocroot + "/bases/" + base + "/root"
	fmt.Println(m["ftphost"], m["ftpfiles"], base)
	os.Exit(0)
	compression := "compression=lz4"
	exec.Command("/sbin/zfs", "create", "-o", compression, "-p",
		*Pool+basepath).Run()
	exec.Command("/sbin/zfs", "create", "-o", compression, "-p",
		*Pool+dlpath).Run()

	txz, _ := uclProp.Args(*Iocroot, "ftpfiles", true)
	txzsp := strings.Split(txz, " ")

	for _, t := range txzsp {
		ftxz, err := os.Create(dlpath + "/" + t)
		if err != nil {
			panic(err)
		}
		tpath := m["ftpdir"] + "/" + base + "/" + t
		white.Printf("Fetching %s...", t)
		err = client.Retrieve(tpath, ftxz)
		if err != nil {
			panic(err)
		} else {
			white.Printf(" done\n")
		}
	}

}

func baseExists(base string, bases []string) bool {
	for _, b := range bases {
		if b == base {
			return true
		}
	}
	return false
}
