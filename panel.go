package gapi

import (
	"errors"
	"fmt"
	"image/png"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

var MissingQueryParameterError error

type GrafanaPanelExport struct {
	Output      string
	Panel       GrafanaPanel
	ExportRange GrafanaPanelExportTimeRange
	ExportSize  GrafanaPanelExportSize

	Dashboard DashboardMeta
	Org       Org
	Tz        string

	DashboardVars map[string][]string
}

func (e *GrafanaPanelExport) AsQueryParameters() string {
	template := "orgId=%v%s&%s&%s"
	panelQuery, _ := e.Panel.AsPartOfUrl()
	rangeQuery, _ := e.ExportRange.AsPartOfUrl()

	dashboardVarsQuery := ""
	if len(e.DashboardVars) > 0 {
		dashboardVarsQuery = "&" + dasboardVarsToQueryString(e.DashboardVars)
	}
	url := fmt.Sprintf(template, e.Org.Id, dashboardVarsQuery, panelQuery, rangeQuery)
	return url
}
func (e *GrafanaPanelExport) AsRenderPartOfUrl() string {

	template := "/render/d-solo/%s/%s?%s&%s&tz=%s"
	sizeQuery, _ := e.ExportSize.AsPartOfUrl()
	loc, _ := time.LoadLocation(e.Tz)
	url := fmt.Sprintf(template, e.Dashboard.UID, e.Dashboard.Title, e.AsQueryParameters(), sizeQuery, url.QueryEscape(loc.String()))
	return url
}

type GrafanaPanel struct {
	ID int
}

func (e *GrafanaPanel) AsPartOfUrl() (string, error) {
	return fmt.Sprintf("panelId=%v", e.ID), nil
}

type GrafanaExportQueryParameter interface {
	AsPartOfUrl() (string, error)
}
type GrafanaPanelExportTimeRange struct {
	From time.Time
	End  time.Time
}

func (r *GrafanaPanelExportTimeRange) New(from time.Time, end time.Time) {
	r.From = from
	r.End = end
}

type GrafanaPanelExportSize struct {
	Width  int
	Height int
}

func (size GrafanaPanelExportSize) AsPartOfUrl() (string, error) {
	if size.Width == 0 {
		return "", MissingQueryParameterError
	}
	if size.Height == 0 {
		return "", MissingQueryParameterError
	}
	template := "width=%v&height=%v"

	return fmt.Sprintf(template, size.Width, size.Height), nil
}

func (timeRange GrafanaPanelExportTimeRange) AsPartOfUrl() (string, error) {
	if timeRange.From.IsZero() {
		return "", MissingQueryParameterError
	}
	if timeRange.End.IsZero() {
		return "", MissingQueryParameterError
	}
	template := "from=%s&to=%s"

	return fmt.Sprintf(template, timeToGrafanaString(timeRange.From), timeToGrafanaString(timeRange.End)), nil
}

func (c *Client) ExportPanelAsImage(export GrafanaPanelExport) error {
	path := export.AsRenderPartOfUrl()
	req, err := c.newRequest("GET", path, nil, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	image, err := png.Decode(resp.Body)
	if err != nil {
		log.Panic(err)
		return err
	}

	f, err := os.Create(export.Output)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// todo define the size
	err = png.Encode(f, image)

	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

func timeToGrafanaString(t time.Time) string {
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
