//go:build integration

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	//

	// Generate a random DNS challenge key for testing
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		t.Fatalf("Failed to generate random challenge key: %v", err)
	}
	randomKey := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Generate a random DNS name suffix for testing (hex encoded to ensure DNS-safe characters)
	randomSuffix := hex.EncodeToString(randomBytes[:8])

	fixture := acmetest.NewFixture(&deSECDNSProviderSolver{},
		acmetest.SetResolvedZone(zone),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/desec"),
		acmetest.SetDNSChallengeKey(randomKey),
		acmetest.SetResolvedFQDN(fmt.Sprintf("cert-manager-dns01-tests-%s.%s", randomSuffix, zone)),
		acmetest.SetPropagationLimit(time.Minute*15),
	)

	// need to uncomment and RunConformance delete runBasic and runExtended once https://github.com/cert-manager/cert-manager/pull/4835 is merged
	// fixture.RunConformance(t)
	fixture.RunBasic(t)
	fixture.RunExtended(t)
}
