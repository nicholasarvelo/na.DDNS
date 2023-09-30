# LFGoD2NS-Cloudflare (Might Change :hurtrealbad:)
This Go application serves as a Dynamic DNS (DDNS) client that interacts with the Cloudflare DNS service. It runs as a containerized service and periodically polls your public IP address. When your public IP changes, the application automatically updates the corresponding DNS record on Cloudflare to ensure that your domain always resolves to the correct IP address.

## Features (Work in Progress)

- **Dynamic DNS Updates:** Automatically updates your Cloudflare DNS records based on changes to your public IP address.
- **Proxy Status :** Choose if record is 'Proxied' or 'DNS Only'.
- **Error Handling:** Error handling and logging for better tracking of any issues.

## Prerequisites

Before using this application, make sure you have the following prerequisites:

- [Docker](https://www.docker.com/) installed on your system.
- A Cloudflare account and API key.
- Familiarity with environment variables.

## Getting Started

1. Clone this repository to your local machine.

2. Set the required environment variables in a `.env` file or in your Docker Compose configuration:

   ```env
   CLOUDFLARE_API_KEY=your_cloudflare_api_key
   CLOUDFLARE_DNS_RECORD=your_dns_record
   CLOUDFLARE_DNS_RECORD_TYPE=your_dns_record_type
   ```

3. Build and run the application as a Docker container:

   ```bash
   docker-compose up -d
   ```

   The application will start polling for IP address changes and updating the Cloudflare DNS record accordingly.

4. Monitor the logs to see the status of DNS record updates:

   ```bash
   docker-compose logs -f
   ```

## Contributing

Feel free to contribute to this project by creating issues or pull requests. Your feedback and contributions are greatly appreciated!

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.