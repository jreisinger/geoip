package check

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/jreisinger/checkip"
)

type censys struct {
	Result result `json:"result"`
}

type result struct {
	Data censysData      `json:"services"`
	OS   operatingSystem `json:"operating_system"`
}

type operatingSystem struct {
	Product                   string `json:"product"`
	Vendor                    string `json:"vendor"`
	Version                   string `json:"version"`
	Edition                   string `json:"edition"`
}

type censysData []struct {
	Port                int    `json:"port"`
	Transport           string `json:"transport_protocol"` // tcp, udp
	ServiceName         string `json:"service_name"`
	ExtendedServiceName string `json:"extended_service_name"`
}

var censysUrl = "https://search.censys.io/api/v2"

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Censys gets generic information from search.censys.io.
func Censys(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "censys.io",
		Type: checkip.TypeInfoSec,
	}

	apiKey, err := getConfigValue("CENSYS_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		return result, nil
	}

	apiSec, err := getConfigValue("CENSYS_SEC")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiSec == "" {
		return result, nil
	}

	headers := map[string]string{
		"Authorization": "Basic " + basicAuth(apiKey, apiSec),
		"Accept":        "application/json",
		"Content-Type":  "application/x-www-form-urlencoded;charset=UTF-8",
	}

	var censys censys
	apiURL := fmt.Sprintf("%s/hosts/%s", censysUrl, ipaddr)
	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &censys); err != nil {
		return result, newCheckError(err)
	}
	result.Info = censys

	for _, d := range censys.Result.Data {
		port := d.Port
		if port != 80 && port != 443 && port != 53 { // undecidable ports
			result.Malicious = true
		}
	}

	return result, nil
}

type byPortC censysData

func (x byPortC) Len() int           { return len(x) }
func (x byPortC) Less(i, j int) bool { return x[i].Port < x[j].Port }
func (x byPortC) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Info returns interesting information from the check.
func (c censys) Summary() string {
	var portInfo []string
	var os string
	if c.Result.OS.Product != "" {
		os = c.Result.OS.Product
	}
	var osvendor string
	if c.Result.OS.Vendor != "" {
		osvendor = c.Result.OS.Vendor
	}

	service := make(map[string]int)
	sort.Sort(byPortC(c.Result.Data))
	for _, d := range c.Result.Data {
		port := d.Port

		sport := fmt.Sprintf("%s/%d", strings.ToLower(d.Transport), port)
		service[sport]++

		if service[sport] > 1 {
			continue
		}

		portInfo = append(portInfo, fmt.Sprintf("%s (%s)", sport, strings.ToLower(d.ServiceName)))
	}

	return fmt.Sprintf("OS: %s, open: %s", strings.Join(nonEmpty(os, osvendor), ", "), strings.Join(portInfo, ", "))
}

func (c censys) Json() ([]byte, error) {
	return json.Marshal(c)
}
