package scraper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinuxNetworkStatsScraper(t *testing.T) {
	t.Run("Network stats parsed from the file with default interface", func(t *testing.T) {
		const wantedReceivedBytes = 3862937603
		const wantedTransmitBytes = 281882792
		assertExpectedNetUsageBytes("testdata/eth0_test.data", t, wantedReceivedBytes, wantedTransmitBytes, "")
	})

	t.Run("Network stats parsed from the file with second interface in the file", func(t *testing.T) {
		const wantedReceivedBytes = 3862937603
		const wantedTransmitBytes = 281882792
		assertExpectedNetUsageBytes("testdata/eth0_test.data", t, wantedReceivedBytes, wantedTransmitBytes, "eth0")
	})

	t.Run("Network stats parsed from the file with third interface in the file", func(t *testing.T) {
		const wantedReceivedBytes = 2247549264
		const wantedTransmitBytes = 255567044
		assertExpectedNetUsageBytes("testdata/eth0_test.data", t, wantedReceivedBytes, wantedTransmitBytes, "eth1")
	})

	t.Run("Network stats parsed from the file with lo interface in the file", func(t *testing.T) {
		const wantedReceivedBytes = 1982736
		const wantedTransmitBytes = 1982736
		assertExpectedNetUsageBytes("testdata/eth0_test.data", t, wantedReceivedBytes, wantedTransmitBytes, "lo")
	})
}

func assertExpectedNetUsageBytes(testFile string, t *testing.T, wantedReceivedBytes uint64, wantedTransmitBytes uint64, interfaceName string) {
	f, err := os.Open(testFile)

	if err != nil {
		assert.Fail(t, "The following error occurred on retrieving the test file: "+err.Error())
	}

	networkStats, err := NewLinuxNetworkDevicesFileScraperWithInterface(interfaceName).Scrape(f)

	if err != nil {
		assert.Fail(t, "Error on scraping the net stats: "+err.Error())
	}

	assert.Equal(t, wantedReceivedBytes, networkStats.ReceivedBytes, "Error on received bytes")
	assert.Equal(t, wantedTransmitBytes, networkStats.TransmittedBytes, "Error on transmitted bytes")
}
