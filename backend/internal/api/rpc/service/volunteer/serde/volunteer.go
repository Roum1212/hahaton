package rpc_service_volunteer_serde

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"

	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	model_date_of_birth "github.com/Roum1212/hahaton/internal/domain/model/date-of-birth"
	model_wallet_operation "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation"
	model_wallet_operation_type "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-type"
)

func AggregateToPBVolunteer(
	volunteer *aggregate_volunteer.Volunteer,
	wallet *aggregate_wallet.Wallet,
) *volunteer_v1.Volunteer {
	name := volunteer.GetFullName()

	return &volunteer_v1.Volunteer{
		Id:               volunteer.GetUserID().String(),
		AvatarUrl:        volunteer.GetAvatarPath().GetURL(),
		DateOfBirth:      DateToPB(volunteer.GetDateOfBirth()),
		FirstName:        name.GetFirstName(),
		LastName:         name.GetLastName(),
		MiddleName:       name.GetMiddleName(),
		EmailAddress:     volunteer.GetEmailAddress().String(),
		PhoneNumber:      volunteer.GetPhoneNumber().String(),
		CoinBalance:      wallet.GetCoinBalance(),
		ConfirmedMinutes: wallet.GetConfirmedMinutes(),
	}
}

func CoordinatorVolunteerToPB(
	volunteer *aggregate_volunteer.Volunteer,
	showContacts bool,
) *volunteer_v1.CoordinatorVolunteer {
	name := volunteer.GetFullName()

	result := &volunteer_v1.CoordinatorVolunteer{
		Id:         volunteer.GetUserID().String(),
		AvatarUrl:  volunteer.GetAvatarPath().GetURL(),
		FirstName:  name.GetFirstName(),
		LastName:   name.GetLastName(),
		MiddleName: name.GetMiddleName(),
	}
	if showContacts {
		email, phone := volunteer.GetEmailAddress().String(), volunteer.GetPhoneNumber().String()
		result.EmailAddress = &email
		result.PhoneNumber = &phone
	}

	return result
}

func DateFromPB(value *date.Date) (model_date_of_birth.DateOfBirth, error) {
	if value == nil {
		return model_date_of_birth.DateOfBirth{}, errors.New("full date of birth is required")
	}

	if value.GetYear() == 0 || value.GetMonth() == 0 || value.GetDay() == 0 {
		return model_date_of_birth.DateOfBirth{}, errors.New("full date of birth is required")
	}

	return model_date_of_birth.NewDateOfBirth(int(value.GetYear()), time.Month(value.GetMonth()), int(value.GetDay()))
}

//nolint:gosec // The domain date is complete, nonfuture and validated before serialization.
func DateToPB(value model_date_of_birth.DateOfBirth) *date.Date {
	return &date.Date{
		Year:  int32(value.GetYear()),
		Month: int32(value.GetMonth()),
		Day:   int32(value.GetDay()),
	}
}

func OperationToPB(value model_wallet_operation.Operation) *volunteer_v1.VolunteerOperation {
	result := &volunteer_v1.VolunteerOperation{
		Id:                    value.GetOperationID().String(),
		CoordinatorId:         value.GetCoordinatorID().String(),
		VolunteerId:           value.GetVolunteerID().String(),
		Description:           value.GetDescription(),
		Type:                  OperationTypeToPB(value.GetOperationType()),
		CoinDelta:             value.GetCoinDelta(),
		MinuteDelta:           value.GetMinuteDelta(),
		CoinBalanceAfter:      value.GetCoinBalanceAfter(),
		ConfirmedMinutesAfter: value.GetConfirmedMinutesAfter(),
		Sequence:              value.GetSequence(),
		CreatedAt:             timestamppb.New(value.GetCreatedAt()),
	}
	switch value.GetOperationType() {
	case model_wallet_operation_type.ParticipationAccrual:
		result.Source = &volunteer_v1.VolunteerOperation_ParticipationId{
			ParticipationId: value.GetSource().GetParticipationID().String(),
		}

	case model_wallet_operation_type.RewardPurchase, model_wallet_operation_type.RewardRefund:
		result.Source = &volunteer_v1.VolunteerOperation_RewardOrderId{
			RewardOrderId: value.GetSource().GetRewardOrderID().String(),
		}
	}

	return result
}

func OperationTypeFromPB(value volunteer_v1.VolunteerOperationType) (model_wallet_operation_type.OperationType, error) {
	switch value {
	case volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_PARTICIPATION_ACCRUAL:
		return model_wallet_operation_type.NewOperationType(model_wallet_operation_type.ParticipationAccrual.Int32())
	case volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_REWARD_PURCHASE:
		return model_wallet_operation_type.NewOperationType(model_wallet_operation_type.RewardPurchase.Int32())
	case volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_REWARD_REFUND:
		return model_wallet_operation_type.NewOperationType(model_wallet_operation_type.RewardRefund.Int32())
	case volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_UNSPECIFIED:
		return 0, errors.New("unspecified operation type")
	default:
		return 0, fmt.Errorf("unknown operation type: %d", value)
	}
}

func OperationTypeToPB(value model_wallet_operation_type.OperationType) volunteer_v1.VolunteerOperationType {
	switch value {
	case model_wallet_operation_type.ParticipationAccrual:
		return volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_PARTICIPATION_ACCRUAL
	case model_wallet_operation_type.RewardPurchase:
		return volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_REWARD_PURCHASE
	case model_wallet_operation_type.RewardRefund:
		return volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_REWARD_REFUND
	default:
		return volunteer_v1.VolunteerOperationType_VOLUNTEER_OPERATION_TYPE_UNSPECIFIED
	}
}
