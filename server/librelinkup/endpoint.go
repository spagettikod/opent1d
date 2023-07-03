package librelinkup

import "fmt"

type Endpoint string

const (
	EndpointAE  Endpoint = "api-ae.libreview.io"
	EndpointAP  Endpoint = "api-ap.libreview.io"
	EndpointAU  Endpoint = "api-au.libreview.io"
	EndpointCA  Endpoint = "api-ca.libreview.io"
	EndpointDE  Endpoint = "api-de.libreview.io"
	EndpointEU  Endpoint = "api-eu.libreview.io"
	EndpointEU2 Endpoint = "api-eu2.libreview.io"
	EndpointFR  Endpoint = "api-fr.libreview.io"
	EndpointJP  Endpoint = "api-jp.libreview.io"
	EndpointUS  Endpoint = "api-us.libreview.io"
)

var (
	Endpoints = map[string]Endpoint{
		"ae":  EndpointAE,
		"ap":  EndpointAP,
		"au":  EndpointAU,
		"ca":  EndpointCA,
		"de":  EndpointDE,
		"eu":  EndpointEU,
		"eu2": EndpointEU2,
		"fr":  EndpointFR,
		"jp":  EndpointJP,
		"us":  EndpointUS,
	}
)

func (e Endpoint) LoginURL() string {
	return fmt.Sprintf("https://%s/llu/auth/login", e)
}

func (e Endpoint) ConnectionsURL() string {
	return fmt.Sprintf("https://%s/llu/connections/", e)
}

func (e Endpoint) GraphURL(patientID string) string {
	return fmt.Sprintf("https://%s/llu/connections/%s/graph", e, patientID)
}
