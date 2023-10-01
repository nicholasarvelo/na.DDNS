# na.ddns - Cloudflare Dynamic DNS Client

**na.ddns** is a Docker application that updates Cloudflare DNS records with your dynamic public IP address at regular intervals. It leverages the Cloudflare API to manage your DNS records dynamically, ensuring your services are always accessible even with changing IP addresses.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Docker installed on your host system.
- Cloudflare API key and DNS record information.
- Basic knowledge of Docker and Cloudflare.

## Usage

1. This section will come, eventually. 

```shell
¯\_(ツ)_/¯
```

## Environment Variables

- `CLOUDFLARE_API_KEY`: Your Cloudflare API key.
- `CLOUDFLARE_DNS_RECORD`: The DNS record you want to update.
- `CLOUDFLARE_DNS_RECORD_TYPE`: The type of DNS record (e.g., A, AAAA).
- `POLLING_INTERVAL`: The interval at which to check for IP updates in minutes.