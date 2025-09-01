package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestOutboxRepository_NewOutboxRepository(t *testing.T) {
	t.Parallel()

	t.Run("it creates repository successfully", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		assert.NotNil(t, repo)
	})
}

func TestOutboxRepository_GetUnpublishedEvents(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when getting unpublished events", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "aggregate_id", "event_data", "published", "created_at", "published_at"}).
			AddRow(1, "activity-1", `{"id": 1}`, false, now, nil).
			AddRow(2, "activity-2", `{"id": 2}`, false, now, nil)

		mock.ExpectQuery(`SELECT \* FROM "outbox_events" WHERE published = \$1 ORDER BY created_at ASC LIMIT \$2`).
			WithArgs(false, 10).
			WillReturnRows(rows)

		events, err := repo.GetUnpublishedEvents(10)
		assert.NoError(t, err)
		assert.Len(t, events, 2)
		assert.Equal(t, uint(1), events[0].ID)
		assert.Equal(t, "activity-1", events[0].AggregateID)
		assert.False(t, events[0].Published)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns error when database query fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		mock.ExpectQuery(`SELECT \* FROM "outbox_events" WHERE published = \$1 ORDER BY created_at ASC LIMIT \$2`).
			WithArgs(false, 5).
			WillReturnError(assert.AnError)

		events, err := repo.GetUnpublishedEvents(5)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "failed to get unpublished events")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns empty slice when no unpublished events exist", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		rows := sqlmock.NewRows([]string{"id", "aggregate_id", "event_data", "published", "created_at", "published_at"})

		mock.ExpectQuery(`SELECT \* FROM "outbox_events" WHERE published = \$1 ORDER BY created_at ASC LIMIT \$2`).
			WithArgs(false, 10).
			WillReturnRows(rows)

		events, err := repo.GetUnpublishedEvents(10)
		assert.NoError(t, err)
		assert.Len(t, events, 0)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOutboxRepository_MarkAsPublished(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when marking event as published", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "outbox_events" SET "published"=\$1,"published_at"=\$2 WHERE id = \$3`).
			WithArgs(true, sqlmock.AnyArg(), uint(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.MarkAsPublished(1)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns error when database update fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "outbox_events" SET "published"=\$1,"published_at"=\$2 WHERE id = \$3`).
			WithArgs(true, sqlmock.AnyArg(), uint(1)).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err = repo.MarkAsPublished(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to mark event as published")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOutboxRepository_MarkOutboxEventAsPublished(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when marking outbox event as published", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		event := &models.OutboxEvent{
			ID:        1,
			Published: false,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "outbox_events" SET "aggregate_id"=\$1,"event_data"=\$2,"published"=\$3,"created_at"=\$4,"published_at"=\$5 WHERE "id" = \$6`).
			WithArgs(event.AggregateID, event.EventData, true, sqlmock.AnyArg(), sqlmock.AnyArg(), event.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.MarkOutboxEventAsPublished(event)
		assert.NoError(t, err)
		assert.True(t, event.Published)
		assert.NotNil(t, event.PublishedAt)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("it returns error when database update fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		logger := utils.NewTestLogger()
		repo := NewOutboxRepository(logger, gormDB)

		event := &models.OutboxEvent{
			ID:        1,
			Published: false,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "outbox_events" SET "aggregate_id"=\$1,"event_data"=\$2,"published"=\$3,"created_at"=\$4,"published_at"=\$5 WHERE "id" = \$6`).
			WithArgs(event.AggregateID, event.EventData, true, sqlmock.AnyArg(), sqlmock.AnyArg(), event.ID).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err = repo.MarkOutboxEventAsPublished(event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to mark outbox event as published")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
