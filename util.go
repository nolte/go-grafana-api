package gapi

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// MissingQueryParameterError the error Type for Missing Query Parameters
var MissingQueryParameterError error

type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (t TimeRange) AsPartOfUrl() (string, error) {
	if t.From == "" {
		return "", MissingQueryParameterError
	}
	if t.To == "" {
		return "", MissingQueryParameterError
	}
	template := "from=%s&to=%s"

	return fmt.Sprintf(template, t.From, t.To), nil
}

func TimeToGrafanaString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10) + "000"
}

func dasboardVarsToQueryString(vars map[string][]string) string {
	var queryString string
	queryString = ""
	template := "var-%s=%v"
	currentElement := 0
	for key, value := range vars {
		currentElement++

		for index, elemet := range value {
			queryString += fmt.Sprintf(template, key, elemet)
			if index+1 < len(value) {
				queryString += "&"
			}
		}

		if currentElement < len(vars) {
			queryString += "&"
		}

	}

	return queryString
}

func buildPathAndQuery(path string, data map[string]string) string {
	pathAndQuery := fmt.Sprintf("%s?", path)
	params := url.Values{}

	for k, v := range data {
		params.Add(k, v)
	}

	return pathAndQuery + params.Encode()
}
