package enricher

import (
	"context"
	"shorter/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

type IpGeoEnricher struct {
	geoClient    *GeoClient
	deviceParser *DeviceParser
}

func NewIpGeoEnricher(geoClient *GeoClient, deviceParser *DeviceParser) *IpGeoEnricher {
	return &IpGeoEnricher{
		geoClient:    geoClient,
		deviceParser: deviceParser,
	}
}

func (e *IpGeoEnricher) Enrich(ctx context.Context, task *ClickTask) (*EnrichedClick, error) {
	timer := prometheus.NewTimer(metrics.EnrichDuration)
	defer timer.ObserveDuration()

	var country, city *string

	if task.IP != "" {
		geoResp, err := e.geoClient.GetLocation(ctx, task.IP)
		if err == nil {
			country = &geoResp.CountryName
			city = &geoResp.City
		}
	}

	device, os, browser := e.deviceParser.Parse(task.UserAgent)

	referer := task.Referer
	var refererPtr *string
	if referer != "" {
		refererPtr = &referer
	}

	return &EnrichedClick{
		Alias:     task.Alias,
		IP:        task.IP,
		Country:   country,
		City:      city,
		Device:    device,
		OS:        os,
		Browser:   browser,
		Referer:   refererPtr,
		Timestamp: task.Timestamp,
	}, nil
}
