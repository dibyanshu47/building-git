package gitrepository

import "gopkg.in/ini.v1"

func RepoDefaultConfig() *ini.File {
	config := ini.Empty()
	config.Section("core").Key("repositoryformatversion").SetValue("0")
	config.Section("core").Key("filemode").SetValue("false")
	config.Section("core").Key("bare").SetValue("false")
	return config
}
