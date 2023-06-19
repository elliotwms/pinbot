package build

var version = ""

func Version() string {
	if version != "" {
		return version
	}

	return "develop"
}
