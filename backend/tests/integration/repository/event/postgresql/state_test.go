//go:build integration

package repository_event_postgresql_test

import (
	"context"
	"time"

	aggregate_event_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event/command"
	model_event_requirements "github.com/Roum1212/hahaton/internal/domain/model/event-requirements"
	model_event_title "github.com/Roum1212/hahaton/internal/domain/model/event-title"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_Save_Completion() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))
	x.Require().NoError(value.Complete(aggregate_event_command.NewComplete(
		"completion-request", "results-fingerprint", repository_fixture.Now().Add(5*time.Hour),
	)))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	assertEvent(x.T(), value, got)
}

func (x *RepositoryTestSuite) TestRepository_Save_UpdateAndClearImage() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value, func(ctx context.Context) error {
		return nil
	}))

	title, err := model_event_title.NewEventTitle("Новое мероприятие")
	x.Require().NoError(err)
	requirements, err := model_event_requirements.NewEventRequirements(nil)
	x.Require().NoError(err)
	x.Require().NoError(value.Update(aggregate_event_command.NewUpdate(
		title, "Новое описание", model_image_path.ImagePath{}, "Новый адрес", value.GetPoint(),
		value.GetSchedule(), requirements, value.GetVolunteerLimit(), repository_fixture.Now(), false,
	)))
	x.Require().NoError(x.repository.Save(ctx, value))
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	assertEvent(x.T(), value, got)
}
