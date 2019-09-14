package gapi

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestCreatePanelExportURL(t *testing.T) {
	var export GrafanaPanelExport

	export.Panel.ID = 1

	timeRange := GrafanaPanelExportTimeRange{}

	str := "2019-09-05T00:00:00.000Z"
	t1, _ := time.Parse(time.RFC3339, str)

	str = "2019-09-05T23:59:59.999Z"
	t2, _ := time.Parse(time.RFC3339, str)

	timeRange.New(t1, t2)

	export.ExportRange = timeRange

	export.Dashboard.Title = "testTitle"
	export.Dashboard.UID = "abc"

	export.ExportSize.Height = 500
	export.ExportSize.Width = 1000
	export.Tz = "Europe/Berlin"
	url := export.AsRenderPartOfUrl()
	expectedURL := "/render/d-solo/abc/testTitle?orgId=0&panelId=1&from=1567641600000&to=1567727999000&width=1000&height=500&tz=Europe%2FBerlin"
	assert.Equal(t, url, expectedURL, "The two words should be the same.")

}

func TestConvertDateToGrafanaUrlQueryFormat(t *testing.T) {
	//layout := "2006-01-02T15:04:05.000Z"
	layout := "02.01.2006 15:04 MST"
	//str := "2019-09-05T00:00:00.000Z"
	str := "05.09.2019 00:00 CEST"
	convertTime, _ := time.Parse(layout, str)
	grafanaFormattedDate := timeToGrafanaString(convertTime)
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
