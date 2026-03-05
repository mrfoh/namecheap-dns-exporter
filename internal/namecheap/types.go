package namecheap

import (
	"encoding/json"
	"fmt"
)

type DNSExport struct {
	Result Result `json:"Result"`
}

type Result struct {
	CustomHostRecords CustomHostRecords `json:"CustomHostRecords"`
}

type CustomHostRecords struct {
	MaxAllowedRecords int64    `json:"MaxAllowedRecords"`
	ReadOnly          bool     `json:"ReadOnly"`
	Records           []Record `json:"Records"`
}

type RecordType string

const (
	A           RecordType = "A"
	CNAME       RecordType = "CNAME"
	MX          RecordType = "MX"
	TXT         RecordType = "TXT"
	AAAA        RecordType = "AAAA"
	NS          RecordType = "NS"
	UrlRedirect RecordType = "URL Redirect"
	SRV         RecordType = "SRV"
	CAA         RecordType = "CAA"
	ALIAS       RecordType = "ALIAS"
)

func recordTypeFromInt(value int) (RecordType, error) {
	switch value {
	case 1:
		return A, nil
	case 2:
		return CNAME, nil
	case 3:
		return MX, nil
	case 4:
		return TXT, nil
	case 5:
		return AAAA, nil
	case 6:
		return NS, nil
	case 7:
		return UrlRedirect, nil
	case 8:
		return SRV, nil
	case 9:
		return CAA, nil
	case 10:
		return ALIAS, nil
	default:
		return "", fmt.Errorf("unknown record type: %d", value)
	}
}

type Record struct {
	HostID     int64      `json:"HostId"`
	ReadOnly   int        `json:"ReadOnly"`
	RecordType RecordType `json:"RecordType"`
	Host       string     `json:"Host"`
	Data       string     `json:"Data"`
	TTL        int64      `json:"Ttl"`
	Priority   int64      `json:"Priority"`
	IsActive   bool       `json:"IsActive"`
	IsDynDNS   bool       `json:"IsDynDns"`
	Service    *string    `json:"Service"`
	Protocol   *string    `json:"Protocol"`
	Target     *string    `json:"Target"`
	Weight     int64      `json:"Weight"`
	Port       int64      `json:"Port"`
}

func (r *Record) UnmarshalJSON(data []byte) error {
	type Alias Record
	aux := &struct {
		RecordType int `json:"RecordType"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	rt, err := recordTypeFromInt(aux.RecordType)
	if err != nil {
		return err
	}
	r.RecordType = rt
	return nil
}
