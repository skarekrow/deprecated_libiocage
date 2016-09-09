package rootCmds

func Args(cmd string) bool {
	switch cmd {
	case
		"activate",
		"create",
		"destroy",
		"deactivate",
		"set":
		return true
	}
	return false
}
