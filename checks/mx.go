package checks

import (
	"encoding/json"
	"log"
	"net"
	"strings"

	"github.com/jreisinger/checkip/check"
)

type MX struct {
	Records map[string][]string `json:"records"` // domain => MX records
}

func (mx MX) Summary() string {
	var s string
	for domain, mxRecords := range mx.Records {
		for i := range mxRecords {
			mxRecords[i] = strings.TrimSuffix(mxRecords[i], ".")
		}
		s += domain + ": " + strings.Join(mxRecords, ", ")
	}
	return check.Na(s)
}

func (mx MX) JsonString() (string, error) {
	b, err := json.Marshal(mx)
	return string(b), err
}

func CheckMX(ipaddr net.IP) (check.Result, error) {
	names, _ := net.LookupAddr(ipaddr.String()) // NOTE: ignoring error

	// Enrich names with a name with 'www.' removed.
	// [www.csh.ac.at.] => [www.csh.ac.at. csh.ac.at.]
	for _, n := range names {
		t := strings.TrimPrefix(n, "www.")
		names = append(names, t)
	}

	// Enrich names with a domain name from AbuseIPDB.
	// [www.csh.ac.at. csh.ac.at.] = > [www.csh.ac.at. csh.ac.at. aco.net]
	r, _ := CheckAbuseIPDB(ipaddr) // NOTE: ignoring error
	j, _ := r.Info.JsonString()
	sr := strings.NewReader(j)
	decoder := json.NewDecoder(sr)
	var a abuseIPDB
	if err := decoder.Decode(&a); err != nil {
		log.Fatal(err)
	}
	names = append(names, a.Domain)

	var mx MX

	for _, n := range names {
		var mxRecords2 []string
		mxRecords, _ := net.LookupMX(n) // NOTE: ingoring error
		for _, r := range mxRecords {
			mxRecords2 = append(mxRecords2, r.Host)
		}
		if _, ok := mx.Records[n]; !ok {
			mx.Records = make(map[string][]string)
		}
		mx.Records[n] = mxRecords2
	}

	return check.Result{
		Name: "dns mx",
		Type: check.TypeInfo,
		Info: mx,
	}, nil
}