package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	// Retrieve Cloudflare API key and DNS record from environment variables
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")
	cloudflareDNSRecord := os.Getenv("CLOUDFLARE_DNS_RECORD")
	cloudflareDNSRecordType := os.Getenv("CLOUDFLARE_DNS_RECORD_TYPE")
	pollingInterval := os.Getenv("POLLING_INTERVAL")

	// Most API calls require a Context
	ctx := context.Background()

	// Parse the zone name
	cloudflareZoneName, errorOccurred := parseZoneName(cloudflareDNSRecord)
	if errorOccurred != nil {
		log.Fatalln(errorOccurred)
	}

	cronjob := cron.New()

	cronEntryID, errorOccurred := cronjob.AddFunc(
		fmt.Sprintf("@every %sm", pollingInterval), func() {
			currentPublicIP, errorOccured := queryPublicIP()
			if errorOccured != nil {
				log.Fatalf("Failed to retrieve public ip: %v", errorOccured)
			}

			// Create a Cloudflare API client using the API key
			apiClient, errorOccurred := cloudflare.NewWithAPIToken(cloudflareAPIKey)
			if errorOccurred != nil {
				log.Fatalln(errorOccurred)
			}

			// Retrieve the Cloudflare zone ID by zone name
			zoneID, errorOccurred := apiClient.ZoneIDByName(cloudflareZoneName)
			if errorOccurred != nil {
				log.Fatalln(
					"Failed to retrieve Cloudflare Zone ID:",
					errorOccurred,
				)
			}

			cloudflareZoneID := cloudflare.ZoneIdentifier(zoneID)

			zoneRecord, _, errorOccurred := apiClient.ListDNSRecords(
				ctx,
				cloudflareZoneID,
				cloudflare.ListDNSRecordsParams{Name: cloudflareDNSRecord},
			)
			if errorOccurred != nil {
				log.Fatalln("Failed to list records:", errorOccurred)
			}

			if len(zoneRecord) == 0 {
				timeStamp := time.Now().Format(time.Stamp)
				comment := fmt.Sprintf("na.ddns [%s]", timeStamp)
				yes := booleanPointer(true)
				_, errorOccurred = apiClient.CreateDNSRecord(
					ctx, cloudflareZoneID, cloudflare.CreateDNSRecordParams{
						Type:      cloudflareDNSRecordType,
						Name:      cloudflareDNSRecord,
						Content:   currentPublicIP,
						Comment:   comment,
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
			} else if zoneRecord[0].Content != currentPublicIP {
				timeStamp := time.Now().Format(time.Stamp)
				comment := stringPointer(
					fmt.Sprintf(
						"na.ddns [%s]",
						timeStamp,
					),
				)
				_, errorOccurred = apiClient.UpdateDNSRecord(
					ctx, cloudflareZoneID, cloudflare.UpdateDNSRecordParams{
						Name:    cloudflareDNSRecord,
						Content: currentPublicIP,
						Comment: comment,
						ID:      zoneRecord[0].ID,
					},
				)
				if errorOccurred != nil {
					log.Fatalf("Unable to update record: %s", errorOccurred)
				} else {
					log.Printf(
						"Record Updated: '%s' is resolving to '%s'",
						cloudflareDNSRecord,
						currentPublicIP,
					)
				}
			} else {
				log.Printf(
					"No update required: '%s' is already resolving to %s",
					cloudflareDNSRecord,
					currentPublicIP,
				)
			}
		},
	)
	if errorOccurred != nil {
		errorMessage := fmt.Errorf("error with cron entry '%v'", cronEntryID)
		fmt.Printf("%v:%v", errorMessage, errorOccurred)
	}

	cronjob.Start()

	log.Println("na.ddns started")

	runtime.Goexit()
}

// booleanPointer returns a pointer to a boolean value, primarily aimed at
// enhancing code readability when dealing with structs representing optional
// boolean fields in Go.
func booleanPointer(boolean bool) *bool {
	return &boolean
}

// stringPointer returns a pointer to a string value, primarily aimed at
// enhancing code readability when dealing with structs representing optional
// string fields in Go.
func stringPointer(string string) *string {
	return &string
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

// parseZoneName function takes 'cloudflareDNSRecord' and splits the FQDN into
// its components returning the domain name as a string and an error if the FQDN
// is invalid.
func parseZoneName(cloudflareDNSRecord string) (string, error) {
	splitDNSRecord := strings.Split(cloudflareDNSRecord, ".")
	if len(splitDNSRecord) != 3 {
		return "", fmt.Errorf("invalid dns record: %s", cloudflareDNSRecord)
	}

	zoneName := strings.Join(
		[]string{splitDNSRecord[1], splitDNSRecord[2]},
		".",
	)

	return zoneName, nil
}
