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
	"strconv"
	"strings"
	"time"
)

func main() {
	// Retrieve and validate environment variables.
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")
	if cloudflareAPIKey == "" {
		log.Fatalln("'CLOUDFLARE_API_KEY' env variable missing")
	}

	cloudflareDNSRecord := os.Getenv("CLOUDFLARE_DNS_RECORD")
	if cloudflareDNSRecord == "" {
		log.Fatalln("'CLOUDFLARE_DNS_RECORD' env variable missing")
	}

	cloudflareDNSRecordType := os.Getenv("CLOUDFLARE_DNS_RECORD_TYPE")
	if cloudflareDNSRecordType == "" {
		cloudflareDNSRecordType = "A"
	} else if cloudflareDNSRecordType != "A" && cloudflareDNSRecordType != "AAAA" {
		log.Fatalf(
			"'CLOUDFLARE_DNS_RECORD_TYPE' must be either 'A' or 'AAAA'. You entered '%v'",
			cloudflareDNSRecordType,
		)
	}

	pollingInterval := os.Getenv("POLLING_INTERVAL")
	if pollingInterval == "" {
		pollingInterval = "3"
	} else {
		parsedInteger, errorOccurred := strconv.Atoi(pollingInterval)
		if errorOccurred != nil {
			log.Fatalln(errorOccurred)
		}
		if parsedInteger < 1 || parsedInteger > 59 {
			log.Fatalf("'POLLING_INTERVAL' must have a value between '1' and '59'")
		}
	}

	var proxiedSetting *bool
	proxied := os.Getenv("PROXIED")
	if proxied == "" {
		proxiedSetting = booleanPointer(false)
	} else {
		proxied, errorOccurred := strconv.ParseBool(proxied)
		proxiedSetting = booleanPointer(proxied)
		if errorOccurred != nil {
			log.Fatalln("'PROXIED' must either be 'true' or 'false'")
		}
	}

	// Most API calls require a Context.
	ctx := context.Background()

	// Parse the zone name.
	cloudflareZoneName, errorOccurred := parseZoneName(cloudflareDNSRecord)
	if errorOccurred != nil {
		log.Fatalf(
			"'CLOUDFLARE_DNS_RECORD' must be 'foo.bar.com'; %s",
			errorOccurred,
		)
	}

	// Initialize a new cron job.
	cronjob := cron.New()

	// Schedule the new cron job.
	cronEntryID, errorOccurred := cronjob.AddFunc(
		fmt.Sprintf("@every %sm", pollingInterval), func() {
			currentPublicIP, errorOccured := queryPublicIP(cloudflareDNSRecordType)
			if errorOccured != nil {
				log.Fatalf("Failed to retrieve public ip: %v", errorOccured)
			}

			// Create a Cloudflare API client using the API key.
			apiClient, errorOccurred := cloudflare.NewWithAPIToken(cloudflareAPIKey)
			if errorOccurred != nil {
				log.Fatalln(errorOccurred)
			}

			// Retrieve the Cloudflare Zone ID.
			zoneID, errorOccurred := apiClient.ZoneIDByName(cloudflareZoneName)
			if errorOccurred != nil {
				log.Fatalln(
					"Failed to retrieve Cloudflare Zone ID:",
					errorOccurred,
				)
			}

			// Retrieve the DNS records for the specified Zone ID.
			cloudflareZoneID := cloudflare.ZoneIdentifier(zoneID)
			zoneRecord, _, errorOccurred := apiClient.ListDNSRecords(
				ctx,
				cloudflareZoneID,
				cloudflare.ListDNSRecordsParams{Name: cloudflareDNSRecord},
			)
			if errorOccurred != nil {
				log.Fatalln("Failed to list records:", errorOccurred)
			}

			// This section of code checks if a DNS record exists. If it doesn't
			// exist, a new DNS record is created. The new record includes the
			// current public IP address, the specified record type, hostname
			// (name), a comment with a creation timestamp, and a setting for
			// whether it should be proxied or not.
			if len(zoneRecord) == 0 {
				timeStamp := time.Now().Format(time.DateTime)
				comment := fmt.Sprintf("na.DDNS [%s]", timeStamp)
				_, errorOccurred = apiClient.CreateDNSRecord(
					ctx, cloudflareZoneID, cloudflare.CreateDNSRecordParams{
						Type:      cloudflareDNSRecordType,
						Name:      cloudflareDNSRecord,
						Content:   currentPublicIP,
						Comment:   comment,
						Proxiable: true,
						Proxied:   proxiedSetting,
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
				// In the event that the DNS record exists, this part of the
				// code verifies whether the IP address (content) of the
				// record differs from the current public IP address. If
				// they are different, the DNS record is updated with the
				// current IP and a timestamp in the record's comment.
			} else if zoneRecord[0].Content != currentPublicIP {
				timeStamp := time.Now().Format(time.Stamp)
				comment := stringPointer(
					fmt.Sprintf(
						"na.DDNS [%s]",
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
				// If none of the above conditions apply, it logs that the
				// existing record is valid and doesn't need any changes.
				log.Printf(
					"Record Valid: '%s' is already resolving to '%s'",
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

	// Start the cron job and log a message stating so.
	cronjob.Start()
	log.Println("na.DDNS started")

	// This runs the program indefinitely.
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

// queryPublicIP retrieves the public address of the local machine
func queryPublicIP(recordType string) (string, error) {
	var protocol string
	if recordType == "A" {
		protocol = "ipv4"
	} else if recordType == "AAAA" {
		protocol = "ipv6"
	}
	url := fmt.Sprintf("https://%s.icanhazip.com", protocol)
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

	// Read the entire response body from the HTTP request and stores it.
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
		return "", fmt.Errorf("you entered '%s'", cloudflareDNSRecord)
	}

	zoneName := strings.Join(
		[]string{splitDNSRecord[1], splitDNSRecord[2]},
		".",
	)

	return zoneName, nil
}
