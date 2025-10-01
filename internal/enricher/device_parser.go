package enricher

import "github.com/mssola/user_agent"

type DeviceParser struct{}

func NewDeviceParser() *DeviceParser {
	return &DeviceParser{}
}

func (d *DeviceParser) Parse(userAgent string) (device, os, browser string) {
	parser := user_agent.New(userAgent)

	if parser.Mobile() {
		device = "mobile"
	} else {
		device = "desktop"
	}

	os = parser.OSInfo().FullName

	browser, _ = parser.Browser()

	return device, os, browser
}