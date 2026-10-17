package model_address_suggestion

import (
	model_geo_point "github.com/Roum1212/hahaton/internal/domain/model/geo-point"
)

type AddressSuggestion struct {
	value              string
	unrestrictedValue  string
	fiasID             *string
	cityFiasID         *string
	settlementFiasID   *string
	cityWithType       *string
	regionWithType     *string
	settlementWithType *string
	point              *model_geo_point.GeoPoint
	qcGeo              *int32
}

func (x AddressSuggestion) GetCityFiasID() *string {
	return x.cityFiasID
}

func (x AddressSuggestion) GetCityWithType() *string {
	return x.cityWithType
}

func (x AddressSuggestion) GetFiasID() *string {
	return x.fiasID
}

func (x AddressSuggestion) GetPoint() *model_geo_point.GeoPoint {
	return x.point
}

func (x AddressSuggestion) GetQcGeo() *int32 {
	return x.qcGeo
}

func (x AddressSuggestion) GetRegionWithType() *string {
	return x.regionWithType
}

func (x AddressSuggestion) GetSettlementFiasID() *string {
	return x.settlementFiasID
}

func (x AddressSuggestion) GetSettlementWithType() *string {
	return x.settlementWithType
}

func (x AddressSuggestion) GetUnrestrictedValue() string {
	return x.unrestrictedValue
}

func (x AddressSuggestion) GetValue() string {
	return x.value
}

func (x AddressSuggestion) Validate() error {
	if x.point != nil {
		return x.point.Validate()
	}

	return nil
}

func NewAddressSuggestion(
	value string,
	unrestrictedValue string,
	fiasID *string,
	cityFiasID *string,
	settlementFiasID *string,
	cityWithType *string,
	regionWithType *string,
	settlementWithType *string,
	point *model_geo_point.GeoPoint,
	qcGeo *int32,
) (AddressSuggestion, error) {
	x := AddressSuggestion{
		value:              value,
		unrestrictedValue:  unrestrictedValue,
		fiasID:             fiasID,
		cityFiasID:         cityFiasID,
		settlementFiasID:   settlementFiasID,
		cityWithType:       cityWithType,
		regionWithType:     regionWithType,
		settlementWithType: settlementWithType,
		point:              point,
		qcGeo:              qcGeo,
	}
	if err := x.Validate(); err != nil {
		return AddressSuggestion{}, err
	}

	return x, nil
}
