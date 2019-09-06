package gapi

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestCreatePanelExportURL(t *testing.T) {
	var export GrafanaPanelExport

	export.Panel.ID = 1

	export.ExportRange.From = "1"
	export.ExportRange.End = "2"

	export.Dashboard.Title = "testTitle"
	export.Dashboard.UID = "abc"

	export.ExportSize.Height = 1000
	export.ExportSize.Width = 500
	export.Tz = "Europe/Berlin"
	url, _ := export.asPartOfUrl()
	expectedURL := "/render/d-solo/abc/testTitle?orgId=0&panelId=1&from=1&end=2&width=1000&height=500&tz=Europe%2FBerlin"
	assert.Equal(t, url, expectedURL, "The two words should be the same.")

}
