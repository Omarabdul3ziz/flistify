package builder

import (
	"fmt"
	"os"
	"strings"
)

func Build(params []string) error {
	var err error
	filePath := params[0]

	// TODO: implement reader
	// TODO: allow multiple lines for the same command
	zerofileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("can't read Zerofile: %v", err.Error())
	}
	lines := strings.Split(string(zerofileBytes), "\n")

	// TODO: pass flist name from flag or something
	err = startRootfs(path)
	if err != nil {
		return fmt.Errorf("can't start the flist: %v", err.Error())
	}

	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}

		content := strings.Split(strings.TrimSpace(line), " ")
		switch content[0] {
		case "FROM":
			err = handleFrom(content[1:])
		case "KERNEL":
			err = handleKernel(content[1:])
		case "RUN":
			err = handleRun(content[1:])
		case "ENV":
			err = handleEnv(content[1:])
		case "ENTRYPOINT":
			err = handleEntrypoint(content[1:])
		default:
			err = fmt.Errorf("%+v keyword not supported", content[0])
		}
	}

	return err
}
