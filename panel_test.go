package gapi

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestCreatePanelRenderURL(t *testing.T) {

	timeRange := TimeRange{}
	timeRange.From = "1567641600000"
	timeRange.To = "1567727999000"

	exportSize := GrafanaPanelExportSize{}
	exportSize.New(1000, 800)
	url, err := buildRenderURL("V3TD6Z5Wk", 1, 2, timeRange, exportSize, nil, "Europe/Berlin")
	expectedURL := "/render/d-solo/V3TD6Z5Wk/neverlandUnicorns?orgId=1&panelId=2&from=1567641600000&to=1567727999000&width=1000&height=800&tz=Europe%2FBerlin"
	assert.Equal(t, err, nil, "No error Expected")
	assert.Equal(t, url, expectedURL, "The two words should be the same.")
}

func TestConvertDateToGrafanaUrlQueryFormat(t *testing.T) {
	//layout := "2006-01-02T15:04:05.000Z"
	layout := "02.01.2006 15:04 MST"
	//str := "2019-09-05T00:00:00.000Z"
	str := "05.09.2019 00:00 CEST"
	convertTime, _ := time.Parse(layout, str)
	grafanaFormattedDate := TimeToGrafanaString(convertTime)
	assert.Equal(t, grafanaFormattedDate, "1567634400000", "The two words should be the same.")
}
func TestConvertDashbardVarsToQueryString(t *testing.T) {
	vars := make(map[string][]string)
	vars["firstVar"] = []string{"test"}
	grafanaFormattedDate := dasboardVarsToQueryString(vars)
	assert.Equal(t, grafanaFormattedDate, "var-firstVar=test", "The two words should be the same.")
}
func TestConvertDashbardMultiplyVarsToQueryString(t *testing.T) {
	vars := make(map[string][]string)
	vars["firstVar"] = []string{"test"}
	vars["scondVar"] = []string{"3"}
	grafanaFormattedDate := dasboardVarsToQueryString(vars)
	assert.Equal(t, grafanaFormattedDate, "var-firstVar=test&var-scondVar=3", "The two words should be the same.")
}
func TestConvertDashbardListVarsToQueryString(t *testing.T) {
	vars := make(map[string][]string)
	vars["listVar"] = []string{"10", "20"}

	grafanaFormattedDate := dasboardVarsToQueryString(vars)
	assert.Equal(t, grafanaFormattedDate, "var-listVar=10&var-listVar=20", "The two words should be the same.")
}
