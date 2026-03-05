# namecheap-dns-exporter

A CLI tool to convert a Namecheap DNS JSON export into a BIND zone file or AWS Route 53 JSON changeset.

## Installation

```bash
go install github.com/mrfoh/namecheap-dns-exporter/cmd@latest
```

Or build from source:

```bash
git clone https://github.com/mrfoh/namecheap-dns-exporter.git
cd namecheap-dns-exporter
go build -o namecheap-dns-exporter ./cmd
```

## Usage

### Export to BIND zone file (default)

```bash
namecheap-dns-exporter -d example.com dns-dump.json
```

### Export to AWS Route 53 JSON

```bash
namecheap-dns-exporter -d example.com -f route53 dns-dump.json > changeset.json
```

Then import into Route 53:

```bash
aws route53 change-resource-record-sets \
  --hosted-zone-id Z1234567890 \
  --change-batch file://changeset.json
```

### Override TTL for all records

Useful during DNS migrations when you want low TTLs for quick cutover:

```bash
namecheap-dns-exporter -d example.com -t 300 dns-dump.json
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--domain` | `-d` | (required) | Domain name for the zone file |
| `--format` | `-f` | `zone` | Output format: `zone` or `route53` |
| `--ttl` | `-t` | `0` | Override TTL for all records (in seconds). `0` preserves original TTLs. |

## Getting the Namecheap DNS export

1. Log in to your Namecheap account
2. Go to **Domain List** > select your domain > **Advanced DNS**
3. Open browser developer tools (Network tab)
4. Look for the API call that returns your DNS records as JSON
5. Save the JSON response to a file

## Supported record types

| Type | Zone file | Route 53 |
|------|-----------|----------|
| A | Yes | Yes |
| AAAA | Yes | Yes |
| CNAME | Yes | Yes |
| MX | Yes | Yes |
| TXT | Yes | Yes |
| NS | Yes | Yes |
| SRV | Yes | Yes |
| CAA | Yes | Yes |
| ALIAS | Comment only | Skipped |
| URL Redirect | Comment only | Skipped |

## License

MIT
