//go:build integration

package integration

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//nolint:paralleltest // These scenarios exercise environment configuration and a shared HTTP transport.
func TestBackendContracts(t *testing.T) {
	seen := make(map[string]bool)

	scenarios := []struct {
		name string
		run  func(x *fixture)
	}{
		{
			"identity",
			checkIdentity,
		},
		{
			"cooperation",
			checkCooperation,
		},
		{
			"events",
			checkEvents,
		},
		{
			"rewards",
			checkRewards,
		},
		{
			"chat",
			checkChat,
		},
		{
			"geo_images",
			checkGeoImages,
		},
		{
			"runtime",
			checkRuntime,
		},
		{
			"optional_fields",
			checkOptionalFields,
		},
		{
			"api_conversions",
			checkAPIConversions,
		},
		{
			"infrastructure",
			checkInfrastructure,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			x := newFixture(t)
			x.seen = seen
			scenario.run(x)
		})
	}

	if len(seen) == 0 {
		t.Skip("PostgreSQL integration scenarios were not run")
	}

	files, err := filepath.Glob("../../proto/*/v1/*_service.proto")
	must(t, err)

	count := 0

	for _, file := range files {
		data, err := os.Open(filepath.Clean(file))
		must(t, err)

		scanner := bufio.NewScanner(data)

		var pkg, service string

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "package ") {
				pkg = strings.TrimSuffix(strings.TrimPrefix(line, "package "), ";")
			}

			if strings.HasPrefix(line, "service ") {
				service = strings.Fields(line)[1]
			}

			if strings.HasPrefix(line, "rpc ") {
				name, _, _ := strings.Cut(strings.TrimPrefix(line, "rpc "), "(")

				method := pkg + "." + service + "/" + name
				if !seen[method] {
					t.Errorf("RPC not exercised: %s", method)
				}

				count++
			}
		}

		must(t, scanner.Err())
		must(t, data.Close())
	}

	require(t, count > 0, "no proto methods discovered")
	t.Logf("exercised %d proto methods", count)
}
