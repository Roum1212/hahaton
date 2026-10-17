//go:build e2e

package rpc

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/require"
)

type identityProvider struct {
	server        *httptest.Server
	signer        jose.Signer
	addressStatus atomic.Int32
}

func (x *identityProvider) token(t *testing.T, subject string, overrides map[string]any) string {
	t.Helper()

	claims := map[string]any{
		"iss":           x.server.URL,
		"aud":           clientID,
		"sub":           subject,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Hour).Unix(),
		"given_name":    "Иван",
		"family_name":   "Иванов",
		"middle_name":   "Иванович",
		"email":         subject + "@example.com",
		"phone_number":  "+79990000001",
		"date_of_birth": "2000-01-01",
	}
	for name, value := range overrides {
		claims[name] = value
	}

	data, err := json.Marshal(claims)
	require.NoError(t, err)
	signed, err := x.signer.Sign(data)
	require.NoError(t, err)
	token, err := signed.CompactSerialize()
	require.NoError(t, err)

	return token
}

func newIdentityProvider(t *testing.T) *identityProvider {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       privateKey,
	},
		(&jose.SignerOptions{}).WithHeader("kid", "e2e"))
	require.NoError(t, err)

	x := &identityProvider{signer: signer}
	x.addressStatus.Store(http.StatusOK)
	x.server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"issuer":                                x.server.URL,
				"jwks_uri":                              x.server.URL + "/keys",
				"id_token_signing_alg_values_supported": []string{"RS256"},
			}))
		case "/keys":
			require.NoError(t, json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
				Key:       &privateKey.PublicKey,
				KeyID:     "e2e",
				Algorithm: "RS256",
				Use:       "sig",
			}}}))
		case "/addresses":
			require.Equal(t, "Token e2e-dadata", r.Header.Get("Authorization"))

			var request struct {
				Query     string `json:"query"`
				FromBound struct {
					Value string `json:"value"`
				} `json:"from_bound"`
				ToBound struct {
					Value string `json:"value"`
				} `json:"to_bound"`
				Count     int32 `json:"count"`
				Locations []struct {
					Country string `json:"country_iso_code"`
				} `json:"locations"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
			require.Equal(t, "Обнинск", request.Query)
			require.Equal(t, "city", request.FromBound.Value)
			require.Equal(t, "settlement", request.ToBound.Value)
			require.EqualValues(t, 10, request.Count)
			require.Len(t, request.Locations, 1)
			require.Equal(t, "RU", request.Locations[0].Country)
			w.WriteHeader(int(x.addressStatus.Load()))
			_, err := w.Write([]byte(`{
 "suggestions": [{
  "value": "Обнинск",
  "unrestricted_value": "Калужская обл, г Обнинск",
  "data": {
   "fias_id": "obninsk", "city_fias_id": "obninsk", "city_with_type": "г Обнинск",
   "region_with_type": "Калужская обл", "geo_lat": "55.1", "geo_lon": "36.6", "qc_geo": "4"
  }
 }]
}`))
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(x.server.Close)

	return x
}
