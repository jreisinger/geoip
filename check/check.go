// Package check defines how to check an IP address. It provides types and
// functions that are useful when writing checks.
package check

import (
	"net"
)

// Type is the type of a check.
type Type string

const (
	TypeInfo    Type = "Info" // provides generic information about the IP address
	TypeSec     Type = "Sec"  // says whether the IP address is considered malicious
	TypeInfoSec Type = "InfoSec"
)

// Check checks an IP address providing generic and/or security information.
type Check func(ipaddr net.IP) (Result, error)

// Result is the results of a check.
type Result struct {
	Name            string // check name
	Type            Type   // check type
	Info            Info   // provided by Info check
	IPaddrMalicious bool   // provided by Sec check
}

// Info is some generic information provided by an Info check.
type Info interface {
	String() string
	JsonString() (string, error)
}

// EmptyInfo is returned by checks that don't provide generic information about
// an IP address.
type EmptyInfo struct {
}

func (EmptyInfo) String() string {
	return Na("")
}

func (EmptyInfo) JsonString() (string, error) {
	return "{}", nil
}

// Na returns "n/a" if s is empty.
func Na(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

// NonEmpty returns strings that are not empty.
func NonEmpty(strings ...string) []string {
	var ss []string
	for _, s := range strings {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
