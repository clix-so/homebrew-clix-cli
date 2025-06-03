package android

import (
	"encoding/xml"
	"os"
)

type Manifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	Package     string      `xml:"package,attr"`
	Application Application `xml:"application"`
}

type Application struct {
	Name string `xml:"name,attr"` // `android:name` => use a custom struct to capture full namespace
}

func extractApplicationClassName(manifestPath string) (string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	type AppWithAnyAttr struct {
		Attrs []xml.Attr `xml:",any,attr"`
	}

	type ManifestWrapper struct {
		XMLName     xml.Name       `xml:"manifest"`
		Package     string         `xml:"package,attr"`
		Application AppWithAnyAttr `xml:"application"`
	}

	var m ManifestWrapper
	if err := xml.Unmarshal(data, &m); err != nil {
		return "", err
	}

	var appName string
	for _, attr := range m.Application.Attrs {
		if attr.Name.Local == "name" {
			appName = attr.Value
			break
		}
	}

	return appName, nil
}
