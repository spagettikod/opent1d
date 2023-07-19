package librelinkup

import "fmt"

type Endpoint struct {
	Hostname string
	Region   string
}

var (
	EndpointAE  = Endpoint{Hostname: "api-ae.libreview.io", Region: "ae"}
	EndpointAP  = Endpoint{Hostname: "api-ap.libreview.io", Region: "ap"}
	EndpointAU  = Endpoint{Hostname: "api-au.libreview.io", Region: "au"}
	EndpointCA  = Endpoint{Hostname: "api-ca.libreview.io", Region: "ca"}
	EndpointDE  = Endpoint{Hostname: "api-de.libreview.io", Region: "de"}
	EndpointEU  = Endpoint{Hostname: "api-eu.libreview.io", Region: "eu"}
	EndpointEU2 = Endpoint{Hostname: "api-eu2.libreview.io", Region: "eu2"}
	EndpointFR  = Endpoint{Hostname: "api-fr.libreview.io", Region: "fr"}
	EndpointJP  = Endpoint{Hostname: "api-jp.libreview.io", Region: "jp"}
	EndpointUS  = Endpoint{Hostname: "api-us.libreview.io", Region: "us"}

	EndpointDefault Endpoint = EndpointEU

	Endpoints = []Endpoint{
		EndpointAE,
		EndpointAP,
		EndpointAU,
		EndpointCA,
		EndpointDE,
		EndpointEU,
		EndpointEU2,
		EndpointFR,
		EndpointJP,
		EndpointUS,
	}
)

func EndpointByRegion(region string) (Endpoint, bool) {
	for _, e := range Endpoints {
		if e.Region == region {
			return e, true
		}
	}
	return Endpoint{}, false
}

func (e Endpoint) LoginURL() string {
	return fmt.Sprintf("https://%s/llu/auth/login", e.Hostname)
}

func (e Endpoint) ConnectionsURL() string {
	return fmt.Sprintf("https://%s/llu/connections/", e.Hostname)
}

func (e Endpoint) GraphURL(patientID string) string {
	return fmt.Sprintf("https://%s/llu/connections/%s/graph", e.Hostname, patientID)
}
