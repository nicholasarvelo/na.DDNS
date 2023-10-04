# na.DDNS - Cloudflare Dynamic DNS Client

One of my favorite features of Google Domains is their native support for Dynamic DNS at no additional cost. It's a shame that all good things eventually [come to an end](https://www.theverge.com/2023/6/16/23763340/google-domains-sunset-sell-squarespace).

I moved all my domains over to Cloudflare, and that's all and good; however, they do not offer any support for DDNS. I wasn't happy with the solutions I came across, so I decided to write my own.

na.DDNS is a Dynamic DNS client for Cloudflare written in Go and designed to run as a Docker container.


## What does it do?

- The client runs as a cron job polling the public ip address of the host system at a user-defined interval.
- If the DNS record does not exist, it creates a new record with the current public IP address. The record can be set to 'DNS-Only' or 'Proxied'.
- If the DNS record exists and the IP address differs from the current public IP, it updates the record with the current IP.
- If the DNS record's defined IP address matches the current public ip address, it the client takes no additional action.

## Run the Docker Image
### Prerequisites
- A Cloudflare account with your domain's DNS managed by Cloudflare.
- An API key generated from your Cloudflare account with the necessary permissions.

### Getting Started
1. Pull down the image
     ```shell
     docker pull steptimeeditor/na.ddns
     ```
2. Set up environment variables with the required configuration:

   - `CLOUDFLARE_API_KEY`: (Required) Your Cloudflare API key.
   - `CLOUDFLARE_DNS_RECORD`: (Required) The DNS record you want to update.
   - `CLOUDFLARE_DNS_RECORD_TYPE`: (Optional) The DNS record type (Defaults to `A`).
   - `POLLING_INTERVAL`: (Optional) The interval (in minutes) between IP checks and updates (Defaults to `3`).
   - `PROXIED`: (Optional) Set to `true` or `false` to determine whether the DNS record should be proxied through Cloudflare (Defaults to `false`).

3. Run the container.
4. Check the container logs for positive affirmation of your life choices up to this point.
