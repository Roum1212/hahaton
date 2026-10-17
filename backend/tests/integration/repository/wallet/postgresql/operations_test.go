//go:build integration

package repository_wallet_postgresql_test

import (
	"context"

	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_wallet_operation "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation"
	model_wallet_operation_type "github.com/Roum1212/hahaton/internal/domain/model/wallet-operation-type"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Operations_AccrualPurchaseRefund() {
	ctx := x.T().Context()
	operations, err := x.repository.ListOperations(ctx)
	x.Require().NoError(err)
	x.Require().Empty(operations)
	wallet := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, wallet))
	coordinatorID := repository_fixture.CreateCoordinator(x.T(), x.database).GetUserID()
	participationID, err := model_participation_id.GenerateParticipationID()
	x.Require().NoError(err)
	result, err := model_event_result.NewEventResult(participationID, 180, false)
	x.Require().NoError(err)
	accrualID, err := model_operation_id.GenerateOperationID()
	x.Require().NoError(err)
	purchaseID, err := model_operation_id.GenerateOperationID()
	x.Require().NoError(err)
	refundID, err := model_operation_id.GenerateOperationID()
	x.Require().NoError(err)
	orderID, err := model_reward_order_id.GenerateRewardOrderID()
	x.Require().NoError(err)

	var expected []model_wallet_operation.Operation

	err = x.repository.Modify(ctx, wallet.GetUserID(), func(ctx context.Context, current *aggregate_wallet.Wallet) error {
		accrual, err := current.Accrue(aggregate_wallet_command.NewAccrue(
			accrualID, coordinatorID, result, "Участие", repository_fixture.Now(),
		))
		if err != nil {
			return err
		}

		purchase, err := current.Purchase(aggregate_wallet_command.NewPurchase(
			purchaseID, coordinatorID, orderID, "Сувенир", 20, repository_fixture.Now(),
		))
		if err != nil {
			return err
		}

		refund, err := current.Refund(aggregate_wallet_command.NewRefund(
			refundID, coordinatorID, orderID, "Сувенир", 20, repository_fixture.Now(),
		))
		if err != nil {
			return err
		}

		expected = []model_wallet_operation.Operation{
			accrual,
			purchase,
			refund,
		}
		for _, operation := range expected {
			if err := x.repository.AppendOperation(ctx, operation); err != nil {
				return err
			}
		}

		return nil
	})
	x.Require().NoError(err)
	got, err := x.repository.Load(ctx, wallet.GetUserID())
	x.Require().NoError(err)
	x.Require().Equal(int64(30), got.GetCoinBalance())
	x.Require().Equal(int64(180), got.GetConfirmedMinutes())
	x.Require().Equal(uint64(3), got.GetLastSequence())
	operations, err = x.repository.ListOperations(ctx)
	x.Require().NoError(err)
	x.Require().Len(operations, 3)

	for _, operation := range operations {
		found := false

		for _, want := range expected {
			if operation.GetOperationID() != want.GetOperationID() {
				continue
			}

			found = true

			x.Require().Equal(want.GetOperationType(), operation.GetOperationType())
			x.Require().Equal(want.GetCoordinatorID(), operation.GetCoordinatorID())
			x.Require().Equal(want.GetVolunteerID(), operation.GetVolunteerID())
			x.Require().Equal(want.GetSource(), operation.GetSource())
			x.Require().Equal(want.GetCoinDelta(), operation.GetCoinDelta())
			x.Require().Equal(want.GetMinuteDelta(), operation.GetMinuteDelta())
			x.Require().Equal(want.GetCoinBalanceAfter(), operation.GetCoinBalanceAfter())
			x.Require().Equal(want.GetConfirmedMinutesAfter(), operation.GetConfirmedMinutesAfter())
			x.Require().Equal(want.GetSequence(), operation.GetSequence())
			x.Require().Equal(want.GetDescription(), operation.GetDescription())
			x.Require().WithinDuration(want.GetCreatedAt(), operation.GetCreatedAt(), 0)
		}

		x.Require().True(found)
	}

	x.Require().Equal(model_wallet_operation_type.ParticipationAccrual, expected[0].GetOperationType())
	x.Require().Equal(model_wallet_operation_type.RewardPurchase, expected[1].GetOperationType())
	x.Require().Equal(model_wallet_operation_type.RewardRefund, expected[2].GetOperationType())
	x.Require().Error(x.repository.AppendOperation(ctx, expected[0]))
}
