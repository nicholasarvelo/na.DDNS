package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Retrieve Cloudflare API key and DNS record from environment variables
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")
	cloudflareDNSRecord := os.Getenv("CLOUDFLARE_DNS_RECORD")
	cloudflareDNSRecordType := os.Getenv("CLOUDFLARE_DNS_RECORD_TYPE")

	currentPublicIP, errorOccured := queryPublicIP()
	if errorOccured != nil {
		log.Fatalf("Failed to retrieve public ip: %v", errorOccured)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Split the Cloudflare DNS record into its components
	cloudflareZoneName, errorOccurred := splitRecord(cloudflareDNSRecord)
	if errorOccurred != nil {
		log.Fatalln(errorOccurred)
	}

	// Create a Cloudflare API client using the API key
	apiClient, errorOccurred := cloudflare.NewWithAPIToken(cloudflareAPIKey)
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
		cloudflare.ListDNSRecordsParams{Name: cloudflareDNSRecord},
	)
	if errorOccurred != nil {
		log.Println("Failed to list records:", errorOccurred)
	}

	if len(zoneRecord) == 0 {
		yes := booleanPointer(true)
		_, errorOccurred = apiClient.CreateDNSRecord(
			ctx, cloudflareZoneID, cloudflare.CreateDNSRecordParams{
				Type:      cloudflareDNSRecordType,
				Name:      cloudflareDNSRecord,
				Content:   currentPublicIP,
				Comment:   "LFGoD2NS-Cloudflare",
				Proxiable: true,
				Proxied:   yes,
			},
		)
		if errorOccurred != nil {
			log.Printf("Failed to create record: %s", errorOccurred)
		} else {
			log.Printf(
				"Record Created: '%s' is resolving to '%s'",
				cloudflareDNSRecord,
				currentPublicIP,
			)
		}
	} else {
		_, errorOccurred = apiClient.UpdateDNSRecord(
			ctx, cloudflareZoneID, cloudflare.UpdateDNSRecordParams{
				Name:    cloudflareDNSRecord,
				Content: currentPublicIP,
				ID:      zoneRecord[0].ID,
			},
		)
		if errorOccurred != nil {
			fmt.Printf("Unable to update record: %s", errorOccurred)
		} else {
			log.Printf(
				"Record Updated: '%s' is resolving to '%s'",
				cloudflareDNSRecord,
				currentPublicIP,
			)
		}
	}
}

// booleanPointer returns a pointer to a boolean value, primarily aimed at
// enhancing code readability when dealing with structs representing optional
// boolean fields in Go.
func booleanPointer(boolean bool) *bool {
	return &boolean
}

// queryPublicIP retrieves the public IPv4 address of the local machine by making
// a GET request to "https://ipv4.icanhazip.com" and returns it as a string.
func queryPublicIP() (string, error) {
	url := "https://ipv4.icanhazip.com"
	request, errorOccurred := http.Get(url)
	if errorOccurred != nil {
		log.Println(errorOccurred)
	}
	// It's good practice to close the response body after processing the
	// response. This ensures that any associated network resources are released
	// and also prevent resource leaks.
	defer func() {
		if errorOccurred := request.Body.Close(); errorOccurred != nil {
			log.Printf("Failed to close response body: %v", errorOccurred)
		}
	}()

	response, errorOccurred := io.ReadAll(request.Body)
	if errorOccurred != nil {
		log.Println(errorOccurred)
	}

	currentPublicIP := strings.TrimRight(string(response), "\n")

	return currentPublicIP, nil
}

// split function takes a Cloudflare DNS record and splits it into recordName and zoneName.
// It returns these components as strings and an error if the DNS record is invalid.
func splitRecord(cloudflareDNSRecord string) (string, error) {
	splitDNSRecord := strings.Split(cloudflareDNSRecord, ".")
	if len(splitDNSRecord) != 3 {
		return "", fmt.Errorf("invalid DNS Record: %s", cloudflareDNSRecord)
	}

	zoneName := strings.Join(
		[]string{splitDNSRecord[1], splitDNSRecord[2]},
		".",
	)

	return zoneName, nil
}
