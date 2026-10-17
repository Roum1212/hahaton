//go:build integration

package integration

import (
	"context"
	"errors"
	"time"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	aggregate_wallet "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet"
	aggregate_wallet_command "github.com/Roum1212/hahaton/internal/domain/aggregate/wallet/command"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_event_result "github.com/Roum1212/hahaton/internal/domain/model/event-result"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
	model_operation_id "github.com/Roum1212/hahaton/internal/domain/model/operation-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	client_postgresql "github.com/Roum1212/hahaton/internal/infra/client/postgresql"
	repository_event "github.com/Roum1212/hahaton/internal/infra/repository/event/postgresql"
	repository_wallet "github.com/Roum1212/hahaton/internal/infra/repository/wallet/postgresql"
)

func checkInfrastructure(x *fixture) {
	c, a := x.coordinator("infra-owner"), x.login("infra-volunteer")
	eventID, err := model_event_id.NewEventID(x.event(c, "Транзакции", 1))
	must(x.t, err)
	userID, err := model_user_id.NewUserID(a.id)
	must(x.t, err)
	participationID, err := model_participation_id.GenerateParticipationID()
	must(x.t, err)
	result, err := model_event_result.NewEventResult(participationID, 180, false)
	must(x.t, err)
	operationID, err := model_operation_id.GenerateOperationID()
	must(x.t, err)

	database := client_postgresql.NewClient(x.pool)
	events, wallets := repository_event.NewEvents(database), repository_wallet.NewWallets(database)
	ctx := x.t.Context()
	wallet, err := wallets.Load(ctx, userID)
	must(x.t, err)

	err = database.InTx(ctx, func(ctx context.Context) error {
		_, err := database.Exec(ctx, "DELETE FROM wallets WHERE volunteer_id = $1", userID.String())
		if err != nil {
			return err
		}

		return wallets.Save(ctx, wallet)
	})
	require(x.t, errors.Is(err, aggregate_wallet.ErrNotFound), "save of missing wallet must fail")
	_, err = wallets.Load(ctx, userID)
	must(x.t, err)

	event, err := events.Load(ctx, eventID)
	must(x.t, err)
	err = database.InTx(ctx, func(ctx context.Context) error {
		if _, err := database.Exec(ctx, "DELETE FROM event_chats WHERE event_id = $1", eventID.String()); err != nil {
			return err
		}

		if _, err := database.Exec(ctx, "DELETE FROM events WHERE id = $1", eventID.String()); err != nil {
			return err
		}

		return events.Save(ctx, event)
	})
	require(x.t, errors.Is(err, aggregate_event.ErrNotFound), "save of missing event must fail")
	_, err = events.Load(ctx, eventID)
	must(x.t, err)
	err = events.Create(ctx, event, func(ctx context.Context) error {
		return nil
	})
	require(x.t, errors.Is(err, aggregate_event.ErrAlreadyExists), "duplicate event must keep domain error")

	stop := errors.New("cancel demo transaction")
	rollback := true
	change := func(ctx context.Context, event *aggregate_event.Event) error {
		if err := event.Cancel(aggregate_event_command.NewCancel(time.Now().UTC())); err != nil {
			return err
		}

		if err := wallets.Modify(ctx, userID, func(ctx context.Context, wallet *aggregate_wallet.Wallet) error {
			operation, err := wallet.Accrue(aggregate_wallet_command.NewAccrue(
				operationID, event.GetCoordinatorID(), result, event.GetTitle().String(), time.Now().UTC(),
			))
			if err != nil {
				return err
			}

			return wallets.AppendOperation(ctx, operation)
		}); err != nil {
			return err
		}

		operations, err := wallets.ListOperations(ctx)
		if err != nil {
			return err
		}

		require(x.t, len(operations) == 1, "query must see writes from transaction context")

		current, err := wallets.Load(ctx, userID)
		if err != nil {
			return err
		}

		require(x.t, current.GetCoinBalance() == 30, "query row must see nested transaction writes")

		if rollback {
			return stop
		}

		return nil
	}

	err = events.Modify(ctx, eventID, change)
	require(x.t, errors.Is(err, stop), "transaction must preserve callback error")
	wallet, err = wallets.Load(ctx, userID)
	must(x.t, err)
	require(x.t, wallet.GetCoinBalance() == 0 && wallet.GetConfirmedMinutes() == 0, "wallet rollback failed")

	operations, err := wallets.ListOperations(ctx)
	must(x.t, err)
	require(x.t, len(operations) == 0, "operation committed before outer transaction")

	event, err = events.Load(ctx, eventID)
	must(x.t, err)
	require(x.t, event.GetStatus() == model_event_status.Published, "event rollback failed")

	rollback = false

	must(x.t, events.Modify(ctx, eventID, change))
	wallet, err = wallets.Load(ctx, userID)
	must(x.t, err)
	require(x.t, wallet.GetCoinBalance() == 30 && wallet.GetConfirmedMinutes() == 180, "wallet commit failed")

	operations, err = wallets.ListOperations(ctx)
	must(x.t, err)
	require(x.t, len(operations) == 1, "operation commit failed")

	event, err = events.Load(ctx, eventID)
	must(x.t, err)
	require(x.t, event.GetStatus() == model_event_status.Canceled, "event commit failed")
}
