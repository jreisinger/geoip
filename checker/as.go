package checker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jreisinger/checkip/check"
)

// AS holds information about an Autonomous System from iptoasn.com.
type AS struct {
	Number      int    `json:"-"`
	FirstIP     net.IP `json:"-"`
	LastIP      net.IP `json:"-"`
	Description string `json:"description"`
	CountryCode string `json:"-"`
}

func (a AS) String() string {
	return fmt.Sprintf("AS description: %s", check.Na(a.Description))
}

func (a AS) JsonString() (string, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

// CheckAs fills in AS data for a given IP address. The data is taken from a TSV
// file ip2asn-combined downloaded from iptoasn.com. The file is created or
// updated as needed.
func CheckAs(ipaddr net.IP) check.Result {
	file := "/var/tmp/ip2asn-combined.tsv"
	url := "https://iptoasn.com/data/ip2asn-combined.tsv.gz"

	if err := check.UpdateFile(file, url, "gz"); err != nil {
		return check.Result{Error: check.NewResultError(err)}
	}

	as, err := asSearch(ipaddr, file)
	if err != nil {
		return check.Result{Error: check.NewResultError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))}
	}

	return check.Result{
		CheckName: "iptoasn.com",
		CheckType: check.TypeInfo,
		Data:      as,
	}
}

// search the ippadrr in tsvFile and if found fills in AS data.
func asSearch(ipaddr net.IP, tsvFile string) (AS, error) {
	tsv, err := os.Open(tsvFile)
	if err != nil {
		return AS{}, err
	}

	as := AS{}
	s := bufio.NewScanner(tsv)
	for s.Scan() {
		line := s.Text()
		fields := strings.Split(line, "\t")
		as.FirstIP = net.ParseIP(fields[0])
		as.LastIP = net.ParseIP(fields[1])
		if ipIsBetween(ipaddr, as.FirstIP, as.LastIP) {
			as.Number, err = strconv.Atoi(fields[2])
			if err != nil {
				return AS{}, fmt.Errorf("converting string to int: %v", err)
			}
			as.CountryCode = fields[3]
			as.Description = fields[4]
			return as, nil
		}
	}
	if s.Err() != nil {
		return AS{}, err
	}
	return as, nil
}

func ipIsBetween(ipAddr, firstIPAddr, lastIPAddr net.IP) bool {
	if bytes.Compare(ipAddr, firstIPAddr) >= 0 && bytes.Compare(ipAddr, lastIPAddr) <= 0 {
		return true
	}
	return false
}