package model

var (
	CmdVersion = []string{"version", "--version", "-v"}
	CmdHelp    = []string{"help", "--help", "-h"}
	CmdPrefix  = []string{"prefix", "--prefix", "-p"}
)

// Contains checks if a string is in a slice
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}
