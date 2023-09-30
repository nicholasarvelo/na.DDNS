package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
	"strings"
)

func main() {
	// Retrieve Cloudflare API key and DNS record from environment variables
	CloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")
	CloudflareDNSRecord := os.Getenv("CLOUDFLARE_DNS_RECORD")

	// Most API calls require a Context
	ctx := context.Background()

	// Split the Cloudflare DNS record into its components
	cloudflareRecordName, cloudflareZoneName, errorOccurred := splitRecord(CloudflareDNSRecord)
	if errorOccurred != nil {
		log.Fatalln(errorOccurred)
	}
	// TODO Remove after variable is implemented in function call
	// Print the extracted Cloudflare record name
	fmt.Println(cloudflareRecordName)

	// Create a Cloudflare API client using the API key
	apiClient, errorOccurred := cloudflare.NewWithAPIToken(CloudflareAPIKey)
	if errorOccurred != nil {
		log.Println(errorOccurred)
	}

	// Retrieve the Cloudflare zone ID by zone name
	zoneID, errorOccurred := apiClient.ZoneIDByName(cloudflareZoneName)
	if errorOccurred != nil {
		log.Println("Failed to retrieve Cloudflare Zone ID:", errorOccurred)
	}

	cloudflareZoneID := cloudflare.ZoneIdentifier(zoneID)

	zoneRecord, _, errorOccurred := apiClient.ListDNSRecords(
		ctx,
		cloudflareZoneID,
		cloudflare.ListDNSRecordsParams{Name: CloudflareDNSRecord},
	)
	if errorOccurred != nil {
		log.Println("Failed to list records:", errorOccurred)
	}

	recordExists := len(zoneRecord) != 0
}

// split function takes a Cloudflare DNS record and splits it into recordName and zoneName.
// It returns these components as strings and an error if the DNS record is invalid.
func splitRecord(CloudflareDNSRecord string) (string, string, error) {
	splitDNSRecord := strings.Split(CloudflareDNSRecord, ".")
	if len(splitDNSRecord) != 3 {
		return "", "", fmt.Errorf("invalid DNS Record: %s", CloudflareDNSRecord)
	}

	recordName := splitDNSRecord[0]
	zoneName := strings.Join([]string{splitDNSRecord[1], splitDNSRecord[2]}, ".")

	return recordName, zoneName, nil
}
