# checkip

Find out geographic location of an IP address. It takes care of downloading the
GeoLite2-City database if it's not present or outdated. You need to set the
`GEOIP_LICENSE_KEY` environment variable in order to download the database. 

Installation:

```
> go install
```

Usage:

```
> checkip 1.1.1.1
, Australia, AU

> checkip $(dig +short reisinge.net)
Frankfurt am Main, Germany, DE

> checkip $(curl -s util.reisinge.net/addr)
Partizanska lupca, Slovakia, SK
```
