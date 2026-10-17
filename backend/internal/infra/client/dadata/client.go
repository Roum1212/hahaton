package client_dadata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	service_suggest_addresses "github.com/Roum1212/hahaton/internal/app/service/suggest-addresses"
	model_address_suggestion "github.com/Roum1212/hahaton/internal/domain/model/address-suggestion"
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
)

const (
	endpoint         = "https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"
	maxResponseBytes = 4 << 20
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func (x Client) Suggest(
	ctx context.Context,
	query string,
	fromBound string,
	toBound string,
	count int32,
) (_ []model_address_suggestion.AddressSuggestion, err error) {
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

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
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

func NewClient(apiKey string, httpClient *http.Client) (Client, error) {
	if httpClient == nil || httpClient.Timeout <= 0 {
		return Client{}, errors.New("dadata http client must have a positive timeout")
	}

	return Client{
		apiKey:     apiKey,
		httpClient: httpClient,
	}, nil
}
