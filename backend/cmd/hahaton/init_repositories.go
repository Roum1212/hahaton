package hahaton

import (
	"context"

	"github.com/Roum1212/hahaton/internal/config"
	repository_city_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/city/postgresql"
	repository_cooperation_invitation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/cooperation-invitation/postgresql"
	repository_cooperation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/cooperation/postgresql"
	repository_coordinator_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/coordinator/postgresql"
	repository_event_chat_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/event-chat/postgresql"
	repository_event_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/event/postgresql"
	repository_participation_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/participation/postgresql"
	repository_reward_order_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/reward-order/postgresql"
	repository_reward_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/reward/postgresql"
	repository_session_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/session/postgresql"
	repository_volunteer_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/volunteer/postgresql"
	repository_wallet_postgresql "github.com/Roum1212/hahaton/internal/infra/repository/wallet/postgresql"
)

type repositorySet struct {
	catalog *repository_city_postgresql.Catalog

	chats        *repository_event_chat_postgresql.EventChats
	cooperations repository_cooperation_postgresql.Cooperations
	coordinators repository_coordinator_postgresql.Coordinators
	events       repository_event_postgresql.Events
	invitations  repository_cooperation_invitation_postgresql.CooperationInvitations

	orders         repository_reward_order_postgresql.RewardOrders
	participations repository_participation_postgresql.Participations
	rewards        repository_reward_postgresql.Rewards
	sessions       repository_session_postgresql.Sessions
	volunteers     repository_volunteer_postgresql.Volunteers
	wallets        repository_wallet_postgresql.Wallets
}

func initRepositories(
	ctx context.Context,
	cfg config.Config,
	databaseClients databaseClientSet,
	externalClients externalClientSet,
) (repositorySet, error) {
	catalog, err := repository_city_postgresql.NewCatalog(ctx, databaseClients.postgresql)
	if err != nil {
		return repositorySet{}, err
	}

	if err := initializeCities(ctx, cfg, catalog); err != nil {
		return repositorySet{}, err
	}

	return repositorySet{
		catalog: catalog,

		chats: repository_event_chat_postgresql.NewEventChats(
			databaseClients.postgresql,
			externalClients.encryption,
		),
		cooperations: repository_cooperation_postgresql.NewCooperations(databaseClients.postgresql),
		coordinators: repository_coordinator_postgresql.NewCoordinators(databaseClients.postgresql),
		events:       repository_event_postgresql.NewEvents(databaseClients.postgresql),
		invitations: repository_cooperation_invitation_postgresql.NewCooperationInvitations(
			databaseClients.postgresql,
			externalClients.encryption,
		),

		orders:         repository_reward_order_postgresql.NewRewardOrders(databaseClients.postgresql),
		participations: repository_participation_postgresql.NewParticipations(databaseClients.postgresql),
		rewards:        repository_reward_postgresql.NewRewards(databaseClients.postgresql),
		sessions:       repository_session_postgresql.NewSessions(databaseClients.postgresql),
		volunteers:     repository_volunteer_postgresql.NewVolunteers(databaseClients.postgresql),
		wallets:        repository_wallet_postgresql.NewWallets(databaseClients.postgresql),
	}, nil
}
