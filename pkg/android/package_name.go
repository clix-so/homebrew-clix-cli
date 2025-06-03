package android

import (
	"io/ioutil"
	"regexp"
)

func GetPackageName(projectRoot string) string {
	appBuildGradlePath := GetAppBuildGradlePath(projectRoot)
	if appBuildGradlePath == "" {
		return ""
	}

	data, err := ioutil.ReadFile(appBuildGradlePath)
	if err != nil {
		return ""
	}
	content := string(data)

	// namespace = "com.example.app"
	namespaceRegex := regexp.MustCompile(`(?m)^\s*namespace\s*=\s*"(.*?)"`)
	// applicationId = "com.example.app"
	appIdRegex := regexp.MustCompile(`(?m)^\s*applicationId\s*=\s*"(.*?)"`)

	// Priority: namespace > applicationId
	if matches := namespaceRegex.FindStringSubmatch(content); len(matches) > 1 {
		return matches[1]
	}
	if matches := appIdRegex.FindStringSubmatch(content); len(matches) > 1 {
		return matches[1]
	}
	return "" // not found
}
