# na.DDNS - Dynamic DNS Client for Cloudflare

One of my favorite features of Google Domains is their native support for Dynamic DNS at no additional cost. It's a shame that all good things eventually [come to an end](https://www.theverge.com/2023/6/16/23763340/google-domains-sunset-sell-squarespace).

I moved all my domains over to Cloudflare, and that's all well and good; however, they do not offer any support for DDNS. I wasn't happy with the solutions I came across. I wanted something simple to deploy and configure with decent logging that quickly assures if it's working and, if so, what it's doing and, if not, what went wrong. 

So I decided to write my own, which resulted in this Dynamic DNS client for Cloudflare, written in Go and designed to run as a Docker container.

## What Does It Do?

* na.DDNS runs as a cron job polling the public ip address of the host system at a user-defined interval.
* If the DNS record does not exist, a new record is created with the current public IP address. Whether it provisions a proxied record or not is also user-defined. 
* If a DNS record already exists and the IP address associated with it is not the same as the current public IP, na.DDNS will update the record with the current IP address.
* If the DNS record's IP matches the public IP, no action is taken.

## Running the Docker Image
### Prerequisites
* A Cloudflare account.
* A user-generated [Cloudflare API token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/).

### Getting Started
1. Pull down the image:
     ```shell
     docker pull steptimeeditor/ns.ddns
     ```
2. Environment variables:
   * `CLOUDFLARE_API_KEY`: (Required) Your Cloudflare API key.
   * `CLOUDFLARE_DNS_RECORD`: (Required) The DNS record you want to update. (e.g. `foo.bar.com`)

3. Additional but optional environment variables:
   * `CLOUDFLARE_DNS_RECORD_TYPE`: (Optional) The DNS record type. Allowed options are 'A' or 'AAAA'. (Defaults to `A`)
   * `POLLING_INTERVAL`: (Optional) The interval (in minutes) between IP checks and updates. Allowed range is `1` through `59`. (Defaults to `3`)
   * `PROXIED`: (Optional) Sets the [Proxy status](https://developers.cloudflare.com/dns/manage-dns-records/reference/proxied-dns-records/) of the DNS record. Allowed options are `true` or `false`. (Defaults to `false`)
     * naDDNS cannot modify the proxied status of an existing record; however, I will be adding this capability soon.

4. Start the container.
5. Feel free to check the container logs for positive affirmation of your life choices up to this point.

