package hub

import (
	"fmt"
	"os"
	"path"

	"github.com/omarabdul3ziz/flistify/pkg/types"
	"github.com/omarabdul3ziz/flistify/pkg/utils"
	"github.com/pkg/errors"
)

const (
	HUB_PUSH_ENDPOINT = "https://hub.grid.tf/api/flist/me/upload"
)

func Push(source string) error {
	fileInfo, err := os.Stat(source)
	if err != nil {
		return errors.Wrapf(err, "file not found: %v", source)
	}

	if fileInfo.IsDir() {
		if utils.IsRootFS(source) {
			name := path.Base(path.Clean(source))
			output := fmt.Sprintf("/tmp/%v.tar.gz", name)
			utils.ExecuteCommand(types.Command{
				Name: "tar",
				Args: []string{"-czf", output, "-C", source, "."},
			})

			source = output
		}
	}

	token := os.Getenv("HUB_JWT")
	if token == "" {
		return errors.New("token is required. export HUB_JWT")
	}

	if err := uploadFile(source, token); err != nil {
		return errors.Wrapf(err, "failed uploading the file: %v", source)
	}

	return nil
}

func uploadFile(path string, token string) error {

	authHeader := fmt.Sprintf("'Authorization: Bearer %s'", token)
	fileForm := fmt.Sprintf("'file=@%s'", path)
	curlCommand := fmt.Sprintf("curl -X Post -H %s -F %s %s", authHeader, fileForm, HUB_PUSH_ENDPOINT)

	if err := utils.ExecuteCommand(types.Command{
		Name: "bash",
		Args: []string{"-c", curlCommand},
	}); err != nil {
		return errors.Wrapf(err, "failed upload flist: %v", path)
	}
	return nil
}
