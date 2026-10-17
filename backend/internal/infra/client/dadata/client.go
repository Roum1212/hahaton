package client_dadata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	service_suggest_addresses "github.com/Roum1212/hahaton/internal/app/service/suggest-addresses"
	model_address_suggestion "github.com/Roum1212/hahaton/internal/domain/model/address-suggestion"
	model_city "github.com/Roum1212/hahaton/internal/domain/model/city"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
)

const (
	maxResponseBytes = 4 << 20
)

type Client struct {
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

func (x Client) FindCities(ctx context.Context, query string, count int32) ([]model_city.City, error) {
	rows, err := x.request(ctx, query, "city", "settlement", count)
	if err != nil {
		return nil, err
	}

	cities := make([]model_city.City, 0, len(rows))
	for _, row := range rows {
		city, err := parseCity(row)
		if err != nil {
			return nil, fmt.Errorf("%w: decode city: %w", service_suggest_addresses.ErrUnavailable, err)
		}

		cities = append(cities, city)
	}

	return cities, nil
}

func (x Client) GetCity(ctx context.Context, id string) (model_city.City, bool, error) {
	cities, err := x.FindCities(ctx, id, 1)
	if err != nil {
		return model_city.City{}, false, err
	}

	for _, city := range cities {
		if city.GetID() == id {
			return city, true, nil
		}
	}

	return model_city.City{}, false, nil
}

func (x Client) GetVersion() string {
	return "dadata"
}

func (x Client) Suggest(
	ctx context.Context,
	query string,
	fromBound string,
	toBound string,
	count int32,
) ([]model_address_suggestion.AddressSuggestion, error) {
	rows, err := x.request(ctx, query, fromBound, toBound, count)
	if err != nil {
		return nil, err
	}

	values := make([]model_address_suggestion.AddressSuggestion, 0, len(rows))
	for _, row := range rows {
		value, err := parseSuggestion(row)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", service_suggest_addresses.ErrUnavailable, err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (x Client) request(
	ctx context.Context,
	query string,
	fromBound string,
	toBound string,
	count int32,
) (_ []map[string]json.RawMessage, err error) {
	if len(x.apiKey) == 0 {
		return nil, service_suggest_addresses.ErrUnavailable
	}

	body, err := json.Marshal(map[string]any{
		"query":      query,
		"from_bound": map[string]string{"value": fromBound},
		"to_bound":   map[string]string{"value": toBound},
		"count":      count,
		"locations":  []map[string]string{{"country_iso_code": "RU"}},
	})
	if err != nil {
		return nil, fmt.Errorf("encode address request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, x.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create address request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Token "+x.apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := x.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", service_suggest_addresses.ErrUnavailable, err)
	}

	defer func() {
		closeError := response.Body.Close()
		if err == nil && closeError != nil {
			err = fmt.Errorf("%w: close address response: %w", service_suggest_addresses.ErrUnavailable, closeError)
		}
	}()

	switch response.StatusCode {
	case http.StatusForbidden:
		return nil, service_suggest_addresses.ErrQuota
	case http.StatusTooManyRequests:
		return nil, service_suggest_addresses.ErrQuota
	case http.StatusOK:
	default:
		return nil, fmt.Errorf(
			"%w: address response status %d",
			service_suggest_addresses.ErrUnavailable,
			response.StatusCode,
		)
	}

	var document map[string]json.RawMessage
	if err := json.NewDecoder(io.LimitReader(response.Body, maxResponseBytes)).Decode(&document); err != nil {
		return nil, fmt.Errorf("%w: decode address response: %w", service_suggest_addresses.ErrUnavailable, err)
	}

	raw, ok := document["suggestions"]
	if !ok {
		return nil, fmt.Errorf("%w: address suggestions are missing", service_suggest_addresses.ErrUnavailable)
	}

	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("%w: decode address suggestions: %w", service_suggest_addresses.ErrUnavailable, err)
	}

	return rows, nil
}

func parseCity(row map[string]json.RawMessage) (model_city.City, error) {
	var data map[string]json.RawMessage
	if err := json.Unmarshal(row["data"], &data); err != nil {
		return model_city.City{}, err
	}

	name := stringValue(data["city"])
	if value := stringValue(data["settlement"]); len(value) != 0 {
		name = value
	}

	if len(name) == 0 && stringValue(data["region_type"]) == "г" {
		name = stringValue(data["region"])
	}

	return model_city.NewCity(
		stringValue(data["fias_id"]),
		name,
		stringValue(data["region_with_type"]),
		parsePoint(data["geo_lat"], data["geo_lon"]),
		nil,
	)
}

func parsePoint(latitudeRaw, longitudeRaw json.RawMessage) *model_geo_point.GeoPoint {
	latitudeText, longitudeText := stringPointer(latitudeRaw), stringPointer(longitudeRaw)
	if latitudeText == nil || longitudeText == nil {
		return nil
	}

	latitude, latitudeError := strconv.ParseFloat(*latitudeText, 64)

	longitude, longitudeError := strconv.ParseFloat(*longitudeText, 64)
	if latitudeError != nil || longitudeError != nil {
		return nil
	}

	point, err := model_geo_point.NewGeoPoint(latitude, longitude)
	if err != nil {
		return nil
	}

	return &point
}

func parseSuggestion(row map[string]json.RawMessage) (model_address_suggestion.AddressSuggestion, error) {
	var data map[string]json.RawMessage
	if err := json.Unmarshal(row["data"], &data); err != nil {
		return model_address_suggestion.AddressSuggestion{}, err
	}

	point := parsePoint(data["geo_lat"], data["geo_lon"])

	var qcGeo *int32

	rawQC := data["qc_geo"]

	value := string(rawQC)
	if text := stringPointer(rawQC); text != nil {
		value = *text
	}

	if number, err := strconv.ParseInt(value, 10, 32); err == nil && number >= 0 && number <= 5 {
		parsed := int32(number)
		qcGeo = &parsed
	}

	valueText, unrestricted := "", ""
	if value := stringPointer(row["value"]); value != nil {
		valueText = *value
	}

	if value := stringPointer(row["unrestricted_value"]); value != nil {
		unrestricted = *value
	}

	return model_address_suggestion.NewAddressSuggestion(
		valueText,
		unrestricted,
		stringPointer(data["fias_id"]),
		stringPointer(data["city_fias_id"]),
		stringPointer(data["settlement_fias_id"]),
		stringPointer(data["city_with_type"]),
		stringPointer(data["region_with_type"]),
		stringPointer(data["settlement_with_type"]),
		point,
		qcGeo,
	)
}

func stringPointer(raw json.RawMessage) *string {
	var value *string
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}

	return value
}

func stringValue(raw json.RawMessage) string {
	value := stringPointer(raw)
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}

func NewClient(apiKey, endpoint string, httpClient *http.Client) (Client, error) {
	address, err := url.Parse(endpoint)
	if err != nil {
		return Client{}, fmt.Errorf("parse dadata endpoint: %w", err)
	}

	if (address.Scheme != "http" && address.Scheme != "https") || len(address.Host) == 0 {
		return Client{}, errors.New("dadata endpoint must be an absolute http or https url")
	}

	if httpClient == nil || httpClient.Timeout <= 0 {
		return Client{}, errors.New("dadata http client must have a positive timeout")
	}

	return Client{
		apiKey:     apiKey,
		endpoint:   endpoint,
		httpClient: httpClient,
	}, nil
}
