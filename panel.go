package gapi

import (
	"errors"
	"fmt"
	"image/png"
	"log"
	"net/url"
	"os"
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
}

func (e *GrafanaPanelExport) asPartOfUrl() (string, error) {

	template := "/render/d-solo/%s/%s?orgId=%v&%s&%s&%s&tz=%s"
	panelQuery, _ := e.Panel.asPartOfUrl()
	rangeQuery, _ := e.ExportRange.asPartOfUrl()
	sizeQuery, _ := e.ExportSize.asPartOfUrl()
	log.Printf(e.Tz)
	loc, _ := time.LoadLocation(e.Tz)
	url := fmt.Sprintf(template, e.Dashboard.UID, e.Dashboard.Title, e.Org.Id, panelQuery, rangeQuery, sizeQuery, url.QueryEscape(loc.String()))
	return url, nil
}

type GrafanaPanel struct {
	ID int
}

func (e *GrafanaPanel) asPartOfUrl() (string, error) {
	return fmt.Sprintf("panelId=%v", e.ID), nil
}

type GrafanaExportQueryParameter interface {
	asPartOfUrl() (string, error)
}
type GrafanaPanelExportTimeRange struct {
	From string
	End  string
}

type GrafanaPanelExportSize struct {
	Width  int
	Height int
}

func (size GrafanaPanelExportSize) asPartOfUrl() (string, error) {
	if size.Width == 0 {
		return "", MissingQueryParameterError
	}
	if size.Height == 0 {
		return "", MissingQueryParameterError
	}
	template := "width=%v&height=%v"

	return fmt.Sprintf(template, size.Width, size.Height), nil
}

func (timeRange GrafanaPanelExportTimeRange) asPartOfUrl() (string, error) {
	if timeRange.From == "" {
		return "", MissingQueryParameterError
	}
	if timeRange.End == "" {
		return "", MissingQueryParameterError
	}
	template := "from=%s&end=%s"

	return fmt.Sprintf(template, timeRange.From, timeRange.End), nil
}

func (c *Client) ExportPanelAsImage(export GrafanaPanelExport) error {
	path, _ := export.asPartOfUrl()
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
