package sampler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsgonz/mule-runtime-master-env-log-receiver/envlogstatsreceiver/internal/stats/scraper"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileBasedSampler(t *testing.T) {
	t.Run("retrieves the sum of the received and transmit from a file.", func(t *testing.T) {
		sampler := NewFileBasedSampler("testdata/test1.data", &BreakLineScraper{})

		want := uint64(1030)

		got, err := sampler.Sample()

		assert.NoError(t, err, "Error on sampling")
		assert.Equal(t, want, got, "Received unexpected result")
	})

	t.Run("when a file does not exists an error is raised", func(t *testing.T) {
		const wanted = "open nonExistingFile.data: no such file or directory"
		sampler := NewFileBasedSampler("nonExistingFile.data", &BreakLineScraper{})

		_, err := sampler.Sample()

		assert.Error(t, err, "Expected an error, but err was nil")
		assert.EqualError(t, err, wanted, "Received unexpected error message")
	})

	t.Run("when the scraper fails, the error is propagated to the sampler", func(t *testing.T) {
		alwaysFailScraper := &AlwaysFailScraper{Error: "Test error"}
		sampler := NewFileBasedSampler("testdata/test1.data", alwaysFailScraper)

		_, err := sampler.Sample()

		assert.Error(t, err, "Expected an error, but err was nil")
		assert.Equal(t, alwaysFailScraper.Error, err.Error(), "Received unexpected error message")
	})
}

func TestFileBasedDeltaSampler(t *testing.T) {
	t.Run("retrieves the delta of the sum of the received and transmit from a file.", func(t *testing.T) {
		tempFile, err := ioutil.TempFile("", "TestFileBasedDeltaSampler-*.txt")
		assert.NoError(t, err, "Error on creating temporary file")
		tempFilePath := tempFile.Name()
		defer os.Remove(tempFilePath)

		err = addValuesToTempFile(tempFile, 200, 400)
		assert.NoError(t, err, "Error on adding values to temp file")

		sampler := NewFileBasedDeltaSampler(tempFilePath, &BreakLineScraper{}, &TestStorage{})

		want := uint64(600)

		got, err := sampler.Sample()

		assert.NoError(t, err, "Error on sampling")
		assert.Equal(t, want, got, "Received unexpected result")

		tempFile, err = os.OpenFile(tempFile.Name(), os.O_WRONLY|os.O_TRUNC, 0666)
		assert.NoError(t, err, "Error on opening temp file")

		err = addValuesToTempFile(tempFile, 400, 600)
		assert.NoError(t, err, "Error on adding values to temp file")

		want = uint64(400)

		got, err = sampler.Sample()

		assert.NoError(t, err, "Error on sampling")
		assert.Equal(t, want, got, "Received unexpected result")
	})
}

func addValuesToTempFile(tempFile *os.File, readBytes uint64, transmitBytes uint64) error {
	// Write the numbers to the file, each on a new line
	content := fmt.Sprintf("%d\n%d", readBytes, transmitBytes)
	if _, err := tempFile.Write([]byte(content)); err != nil {
		return err
	}

	// Close the file
	if err := tempFile.Close(); err != nil {
		return err
	}
	return nil
}

// Scrapers for testing

type AlwaysFailScraper struct {
	Error string
}

func (s *AlwaysFailScraper) Scrape(io.Reader) (scraper.NetworkStats, error) {
	return scraper.NetworkStats{}, errors.New(s.Error)
}

type BreakLineScraper struct{}

func (s *BreakLineScraper) Scrape(r io.Reader) (scraper.NetworkStats, error) {
	// Read all content from the reader
	content, err := io.ReadAll(r)
	if err != nil {
		return scraper.NetworkStats{}, fmt.Errorf("failed to read content: %w", err)
	}

	// Split the content by newline
	parts := bytes.SplitN(content, []byte("\n"), 2)
	if len(parts) < 2 {
		return scraper.NetworkStats{}, errors.New("expected two parts separated by a newline")
	}

	// Convert the first part to uint64
	firstPart, err := strconv.ParseUint(string(parts[0]), 10, 64)
	if err != nil {
		return scraper.NetworkStats{}, fmt.Errorf("failed to parse first part as uint64: %w", err)
	}

	// Convert the second part to uint64
	secondPart, err := strconv.ParseUint(string(parts[1]), 10, 64)
	if err != nil {
		return scraper.NetworkStats{}, fmt.Errorf("failed to parse second part as uint64: %w", err)
	}

	return scraper.NetworkStats{ReceivedBytes: firstPart, TransmittedBytes: secondPart}, nil
}

// Storage for testing

type TestStorage struct {
	LastCount uint64
}

func (s *TestStorage) Load() (uint64, error) {
	return s.LastCount, nil
}

func (s *TestStorage) Save(value uint64) error {
	s.LastCount = value
	return nil
}
