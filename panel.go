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

const (
	// GrafanaAPIPanelRender Base for the Grafana Render API
	GrafanaAPIPanelRender = "/render/d-solo"

	// GrafanaAPIPanelRenderWidth Default Panel Export Width
	GrafanaAPIPanelRenderWidth = 1000

	// GrafanaAPIPanelRenderHeight Default Panel Export Height
	GrafanaAPIPanelRenderHeight = 500
)

type DashboardPanel struct {
	Bars         bool   `json:"bars"`
	DashLength   int    `json:"dashLength"`
	Dashes       bool   `json:"dashes"`
	Description  string `json:"description"`
	Fill         int    `json:"fill"`
	FillGradient int    `json:"fillGradient"`
	GridPos      struct {
		H int `json:"h"`
		W int `json:"w"`
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"gridPos"`
	ID     int64 `json:"id"`
	Legend struct {
		Avg     bool `json:"avg"`
		Current bool `json:"current"`
		Max     bool `json:"max"`
		Min     bool `json:"min"`
		Show    bool `json:"show"`
		Total   bool `json:"total"`
		Values  bool `json:"values"`
	} `json:"legend"`
	Lines         bool   `json:"lines"`
	Linewidth     int    `json:"linewidth"`
	Links         []Link `json:"links"`
	NullPointMode string `json:"nullPointMode"`
	Options       struct {
		DataLinks []interface{} `json:"dataLinks"`
	} `json:"options"`
	Percentage      bool          `json:"percentage"`
	Pointradius     int           `json:"pointradius"`
	Points          bool          `json:"points"`
	Renderer        string        `json:"renderer"`
	SeriesOverrides []interface{} `json:"seriesOverrides"`
	SpaceLength     int           `json:"spaceLength"`
	Stack           bool          `json:"stack"`
	SteppedLine     bool          `json:"steppedLine"`
	Thresholds      []interface{} `json:"thresholds"`
	TimeFrom        interface{}   `json:"timeFrom"`
	TimeRegions     []interface{} `json:"timeRegions"`
	TimeShift       interface{}   `json:"timeShift"`
	Title           string        `json:"title"`
	Tooltip         struct {
		Shared    bool   `json:"shared"`
		Sort      int    `json:"sort"`
		ValueType string `json:"value_type"`
	} `json:"tooltip"`
	Type  string `json:"type"`
	Xaxis struct {
		Buckets interface{}   `json:"buckets"`
		Mode    string        `json:"mode"`
		Name    interface{}   `json:"name"`
		Show    bool          `json:"show"`
		Values  []interface{} `json:"values"`
	} `json:"xaxis"`
	Yaxes []struct {
		Format  string      `json:"format"`
		Label   interface{} `json:"label"`
		LogBase int         `json:"logBase"`
		Max     interface{} `json:"max"`
		Min     interface{} `json:"min"`
		Show    bool        `json:"show"`
	} `json:"yaxes"`
	Yaxis struct {
		Align      bool        `json:"align"`
		AlignLevel interface{} `json:"alignLevel"`
	} `json:"yaxis"`
}

func (p DashboardPanel) AsPartOfUrl() string {
	return fmt.Sprintf("panelId=%v&fullscreen", p.ID)
}

type GrafanaPanelExportSize struct {
	Width  int
	Height int
}

func (size *GrafanaPanelExportSize) New(width int, height int) {
	if width == 0 {
		size.Width = GrafanaAPIPanelRenderWidth
	} else {
		size.Width = width
	}
	if height == 0 {
		size.Height = GrafanaAPIPanelRenderHeight
	} else {
		size.Height = height
	}
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

func buildRenderURL(dashboardID string,
	orgID int64,
	panelID int64,
	timeRange TimeRange,
	exportSize GrafanaPanelExportSize,
	dashboardVars map[string][]string,
	timeZone string,
) (string, error) {

	// set the base from the url
	exportURL := GrafanaAPIPanelRender

	//append dashboard url part (the dashboardslug is not important for the Request)
	exportURL += "/" + dashboardID + "/neverlandUnicorns?"

	// append the OrgID to the url
	exportURL += fmt.Sprintf("orgId=%v", orgID)

	// append the panelId to the url
	exportURL += fmt.Sprintf("&panelId=%v", panelID)

	// append the TimeRange to the url
	timeRangeQueryParm, err := timeRange.AsPartOfUrl()
	if err != nil {
		return "", err
	}
	exportURL += "&" + timeRangeQueryParm

	dashboardVarsString := dasboardVarsToQueryString(dashboardVars)
	if dashboardVarsString != "" {
		exportURL += "&" + dashboardVarsString
	}

	exportSizeQueryPart, err := exportSize.AsPartOfUrl()
	if err != nil {
		return "", err
	}
	exportURL += "&" + exportSizeQueryPart

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return "", err
	}
	exportURL += fmt.Sprintf("&tz=%v", url.QueryEscape(loc.String()))

	return exportURL, nil
}

// ExportPanelAsImage a Grafana Panel to the Local Filesystem
// Save the Panel as PNG
func (c *Client) ExportPanelAsImage(
	dashboardID string,
	orgID int64,
	panelID int64,
	timeRange TimeRange,
	exportSize GrafanaPanelExportSize,
	dashboardVars map[string][]string,
	timeZone string,
	output string) error {

	renderURL, err := buildRenderURL(dashboardID, orgID, panelID, timeRange, exportSize, dashboardVars, timeZone)
	if err != nil {
		return err
	}
	req, err := c.newRequest("GET", renderURL, nil, nil)
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
	//myImage := image.NewRGBA(image.Rect(0, 0, 100, 200))
	image, err := png.Decode(resp.Body)
	if err != nil {
		log.Panic(err)
		return err
	}

	f, err := os.Create(output)
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
