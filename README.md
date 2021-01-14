[![Build Status](https://travis-ci.org/jreisinger/checkip.svg?branch=master)](https://travis-ci.org/jreisinger/checkip)

# checkip

`checkip` is a CLI tool that finds out information about an IP address. 

```
$ checkip 1.1.1.1
AS          13335 | 1.1.1.0 - 1.1.1.255 | CLOUDFLARENET - Cloudflare, Inc. | US
DNS         one.one.one.one.
ThreatCrowd most users have voted this malicious
AbuseIPDB   reported malicious 8 times (32% confidence) | cloudflare.com | Content Delivery Network
GEO         city unknown | Australia | AU
VirusTotal  analysis stats: 0 malicious, 2 suspicious, 87 harmless
```

## Features

* Necessary files (AS, GEO) are automatically downloaded and updated in the background.
* Checks are done concurrently to save time.
* Output is colored to improve readability.
* It's easy to add new checks.

Currently these types of information are provided:

* AS data using TSV file from [iptoasn](https://iptoasn.com/).
* DNS name using [net.LookupAddr](https://golang.org/pkg/net/#LookupAddr) Go function.
* [ThreatCrowd](https://www.threatcrowd.org/) voting about whether the IP address is malicious.
---
* [AbuseIPDB](https://www.abuseipdb.com) confidence score that the IP address is malicious. You need to [register](https://www.abuseipdb.com/register?plan=free) to get the API key (it's free).
* GEOgraphic location using [GeoLite2 City database](https://dev.maxmind.com/geoip/geoip2/geolite2/) file. You need to [register](https://dev.maxmind.com/geoip/geoip2/geolite2/#Download_Access) to get the license key (it's free).
* [VirusTotal](https://developers.virustotal.com/v3.0/reference#ip-object) scanners results. You need to [register](https://www.virustotal.com/gui/join-us) to to get the API key (it's free).

You can store LICENSE/API keys in `~/.checkip.yaml` or in environment variables.

## Installation

Download the latest [release](https://github.com/jreisinger/checkip/releases) for your operating system and architecture. Or clone the repo and run `make install`.

## Development

```
vim main.go
make install # version defaults to "dev" if VERSION envvar is not set

make release # you'll find releases in releases/ directory
```

Builds are done inside Docker container. When you push to GitHub Travis CI will
try and build a release for you and publish it on GitHub.

Check test coverage:

```
go test -coverprofile cover.out ./...
go tool cover -html=cover.out
```
