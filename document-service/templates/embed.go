package templates

import "embed"

//go:embed *.typ
var FS embed.FS

func GetTemplate(filename string) (string, error) {
	bytes, err := FS.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
