package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Redd MCP Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	getAlerts := mcp.NewTool(
		"get_alerts",
		mcp.WithDescription("Get weather alerts for a US state."),
		mcp.WithString(
			"state",
			mcp.Required(),
			mcp.Description("The US state to get alerts for.")),
	)
	s.AddTool(getAlerts, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		state := request.Params.Arguments["state"].(string)
		url := fmt.Sprintf("https://api.weather.gov/alerts/active/area/%s", state)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http error: %s", resp.Status)
		}
		var payload WeatherAlert
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}
		if len(payload.Features) == 0 {
			return mcp.NewToolResultText("No active alerts for this state."), nil
		}
		var alerts []string
		for _, feature := range payload.Features {
			alerts = append(alerts, fmt.Sprintf("Event: %s\nArea: %s\nSeverity: %s\nDescription: %s\nInstruction: %s\n",
				feature.Properties.Event,
				feature.Properties.AreaDesc,
				feature.Properties.Severity,
				feature.Properties.Description,
				feature.Properties.Instruction))
		}
		return mcp.NewToolResultText(strings.Join(alerts, "\n--\n")), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

type WeatherAlert struct {
	Features []WeatherAlertFeature `json:"features"`
}

type WeatherAlertFeature struct {
	Properties WeatherAlertFeatureProperty `json:"properties"`
}

type WeatherAlertFeatureProperty struct {
	Event       string `json:"event"`
	AreaDesc    string `json:"areaDesc"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`
}
