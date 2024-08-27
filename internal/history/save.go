package history

import (
	"encoding/json"
	"fmt"
	"github.com/byawitz/ggh/internal/config"
	"os"
	"slices"
	"strings"
	"time"
)

func AddHistoryFromArgs(args []string) {
	if len(args) == 1 && !strings.Contains(args[0], "@") {
		localConfig, err := config.GetConfig(args[0])
		if err != nil || localConfig.Name == "" {
			fmt.Printf("couldn't fetch %s from config file, error:%v.\n", args[0], err)
			return
		}

		AddHistory(localConfig)
		return
	}

	generatedConfig := config.SSHConfig{}

	skipNext := false
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		switch {
		case strings.HasPrefix(arg, "-p"):
			if arg == "-p" {
				generatedConfig.Port = args[i+1]
				skipNext = true
			} else {
				generatedConfig.Port = args[i][2:]
			}
		case arg == "-i":
			generatedConfig.Key = args[i+1]
			skipNext = true
		case strings.Contains(arg, "@"):
			values := strings.Split(arg, "@")
			generatedConfig.User = values[0]
			generatedConfig.Host = values[1]
		}
	}
	AddHistory(generatedConfig)
}

func AddHistory(c config.SSHConfig) {
	if c.Host == "" {
		return
	}

	list, err := Fetch(getFile())

	if err != nil {
		fmt.Println("error getting ggh file")
		return
	}

	err = saveFile(SSHHistory{Connection: c, Date: time.Now()}, list)
	if err != nil {
		fmt.Println("error saving ggh file")
		return
	}
}

func saveFile(n SSHHistory, l []SSHHistory) error {
	file := getFileLocation()
	fileContent := stringify(n, l)

	err := os.WriteFile(file, []byte(fileContent), 0644)

	return err
}

func stringify(n SSHHistory, l []SSHHistory) string {
	history := make([]SSHHistory, 0)

	for i, sshHistory := range l {
		if sshHistory.Connection.Host == n.Connection.Host &&
			sshHistory.Connection.Name == n.Connection.Name {
			l = slices.Delete(l, i, i+1)
		}
	}

	history = append(history, n)
	history = append(history, l...)
	content, err := json.Marshal(history)

	if err != nil {
		return ""
	}

	return string(content)
}
