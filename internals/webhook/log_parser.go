package webhook

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)


func parseExtractedFiles(folder string) (map[string]string, error) {
	values := map[string]string{
		"game":      "",
		"run_id":    "",
		"user_id":   "",
		"server_ip": "",
	}

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(strings.ToLower(info.Name()), "output") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", path, err)
			}

			content := string(data)
			extracted := extractValues(content)

			for key, value := range extracted {
				if value != "" {
					values[key] = value
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking through folder %s: %v", folder, err)
	}

	if values["game"] == "" && values["run_id"] == "" && values["user_id"] == "" && values["server_ip"] == "" {
		return nil, nil
	}

	return values, nil
}

func extractValues(content string) map[string]string {
	values := map[string]string{
		"game":      "",
		"run_id":    "",
		"user_id":   "",
		"server_ip": "",
	}

	gamePattern := regexp.MustCompile(`"game":\s*"(.+?)"`)
	runIDPattern := regexp.MustCompile(`"run_id":\s*"(.+?)"`)
	userIDPattern := regexp.MustCompile(`"user_id":\s*"(.+?)"`)
	serverIPPattern := regexp.MustCompile(`server_ip\s*[=:]\s*"(.+?)"`)

	gameMatch := gamePattern.FindStringSubmatch(content)
	runIDMatch := runIDPattern.FindStringSubmatch(content)
	userIDMatch := userIDPattern.FindStringSubmatch(content)
	serverIPMatch := serverIPPattern.FindStringSubmatch(content)

	if len(gameMatch) > 1 {
		values["game"] = gameMatch[1]
	}
	if len(runIDMatch) > 1 {
		values["run_id"] = runIDMatch[1]
	}
	if len(userIDMatch) > 1 {
		values["user_id"] = userIDMatch[1]
	}
	if len(serverIPMatch) > 1 {
		values["server_ip"] = serverIPMatch[1]
	}

	return values
}


