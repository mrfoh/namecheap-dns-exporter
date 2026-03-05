package namecheap

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Format string

const (
	FormatZone     Format = "zone"
	FormatRoute53  Format = "route53"
)

func Export(r io.Reader, w io.Writer, domain string, format Format) error {
	var export DNSExport
	if err := json.NewDecoder(r).Decode(&export); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	records := export.Result.CustomHostRecords.Records

	switch format {
	case FormatRoute53:
		return exportRoute53(w, records, domain)
	default:
		return exportZone(w, records, domain)
	}
}

func exportZone(w io.Writer, records []Record, domain string) error {
	fmt.Fprintf(w, "; Zone file exported from Namecheap DNS for %s\n", domain)
	fmt.Fprintf(w, "$ORIGIN %s.\n\n", domain)

	for _, rec := range records {
		if !rec.IsActive {
			continue
		}
		line := formatZoneRecord(rec)
		fmt.Fprintln(w, line)
	}

	return nil
}

func formatZoneRecord(rec Record) string {
	host := rec.Host
	if host == "@" {
		host = "@"
	}

	data := rec.Data
	ttl := rec.TTL

	switch rec.RecordType {
	case A, AAAA:
		return fmt.Sprintf("%-40s %d IN %-5s %s", host, ttl, rec.RecordType, data)
	case CNAME:
		data = ensureTrailingDot(data)
		return fmt.Sprintf("%-40s %d IN CNAME %s", host, ttl, data)
	case MX:
		data = ensureTrailingDot(data)
		return fmt.Sprintf("%-40s %d IN MX    %d %s", host, ttl, rec.Priority, data)
	case TXT:
		if !strings.HasPrefix(data, "\"") {
			data = fmt.Sprintf("%q", data)
		}
		return fmt.Sprintf("%-40s %d IN TXT   %s", host, ttl, data)
	case NS:
		data = ensureTrailingDot(data)
		return fmt.Sprintf("%-40s %d IN NS    %s", host, ttl, data)
	case SRV:
		target := ""
		if rec.Target != nil {
			target = ensureTrailingDot(*rec.Target)
		} else {
			target = ensureTrailingDot(data)
		}
		return fmt.Sprintf("%-40s %d IN SRV   %d %d %d %s", host, ttl, rec.Priority, rec.Weight, rec.Port, target)
	case CAA:
		return fmt.Sprintf("%-40s %d IN CAA   %s", host, ttl, data)
	case ALIAS:
		data = ensureTrailingDot(data)
		return fmt.Sprintf("; ALIAS (non-standard): %s -> %s", host, data)
	case UrlRedirect:
		return fmt.Sprintf("; URL Redirect: %s -> %s", host, data)
	default:
		return fmt.Sprintf("; Unknown record type %s: %s -> %s", rec.RecordType, host, data)
	}
}

// Route 53 JSON types

type route53ChangeBatch struct {
	Comment string         `json:"Comment"`
	Changes []route53Change `json:"Changes"`
}

type route53Change struct {
	Action            string                `json:"Action"`
	ResourceRecordSet route53ResourceRecord `json:"ResourceRecordSet"`
}

type route53ResourceRecord struct {
	Name            string              `json:"Name"`
	Type            string              `json:"Type"`
	TTL             int64               `json:"TTL"`
	ResourceRecords []route53RecordValue `json:"ResourceRecords"`
}

type route53RecordValue struct {
	Value string `json:"Value"`
}

func exportRoute53(w io.Writer, records []Record, domain string) error {
	var changes []route53Change

	for _, rec := range records {
		if !rec.IsActive {
			continue
		}
		if rec.RecordType == ALIAS || rec.RecordType == UrlRedirect {
			continue
		}

		fqdn := fqdnFromHost(rec.Host, domain)
		value := formatRoute53Value(rec)

		changes = append(changes, route53Change{
			Action: "UPSERT",
			ResourceRecordSet: route53ResourceRecord{
				Name: fqdn,
				Type: string(rec.RecordType),
				TTL:  rec.TTL,
				ResourceRecords: []route53RecordValue{
					{Value: value},
				},
			},
		})
	}

	batch := route53ChangeBatch{
		Comment: fmt.Sprintf("Imported from Namecheap DNS for %s", domain),
		Changes: changes,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(batch)
}

func formatRoute53Value(rec Record) string {
	switch rec.RecordType {
	case MX:
		return fmt.Sprintf("%d %s", rec.Priority, ensureTrailingDot(rec.Data))
	case TXT:
		data := rec.Data
		if !strings.HasPrefix(data, "\"") {
			data = fmt.Sprintf("%q", data)
		}
		return data
	case SRV:
		target := rec.Data
		if rec.Target != nil {
			target = *rec.Target
		}
		return fmt.Sprintf("%d %d %d %s", rec.Priority, rec.Weight, rec.Port, ensureTrailingDot(target))
	case CNAME, NS:
		return ensureTrailingDot(rec.Data)
	default:
		return rec.Data
	}
}

func fqdnFromHost(host, domain string) string {
	if host == "@" {
		return domain + "."
	}
	return host + "." + domain + "."
}

func ensureTrailingDot(s string) string {
	if !strings.HasSuffix(s, ".") {
		return s + "."
	}
	return s
}
