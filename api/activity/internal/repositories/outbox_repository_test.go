package repositories

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newGormWithSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}
	return gormDB, mock, db
}

func TestOutboxRepository_GetUnpublishedEvents(t *testing.T) {
	logger := utils.NewTestLogger()
	gormDB, mock, sqlDB := newGormWithSQLMock(t)
	defer sqlDB.Close()

	repo := NewOutboxRepository(logger, gormDB)

	// success
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "outbox_events" WHERE published = $1 ORDER BY created_at ASC LIMIT $2`)).
		WithArgs(false, 100).
		WillReturnRows(sqlmock.NewRows([]string{"id", "aggregate_id", "event_data", "published", "created_at", "published_at"}).
			AddRow(1, "activity-1", `{"x":1}`, false, time.Now(), nil))

	events, err := repo.GetUnpublishedEvents(context.Background(), 100)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	assert.NoError(t, mock.ExpectationsWereMet())

	// db error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "outbox_events" WHERE published = $1 ORDER BY created_at ASC LIMIT $2`)).
		WithArgs(false, 50).
		WillReturnError(assert.AnError)

	_, err = repo.GetUnpublishedEvents(context.Background(), 50)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOutboxRepository_MarkAsPublished(t *testing.T) {
	logger := utils.NewTestLogger()
	gormDB, mock, sqlDB := newGormWithSQLMock(t)
	defer sqlDB.Close()

	repo := NewOutboxRepository(logger, gormDB)

	// success
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "published"=$1,"published_at"=$2 WHERE id = $3`)).
		WithArgs(true, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.MarkAsPublished(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// db error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "published"=$1,"published_at"=$2 WHERE id = $3`)).
		WithArgs(true, sqlmock.AnyArg(), 2).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err = repo.MarkAsPublished(context.Background(), 2)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOutboxRepository_MarkOutboxEventAsPublished(t *testing.T) {
	logger := utils.NewTestLogger()
	gormDB, mock, sqlDB := newGormWithSQLMock(t)
	defer sqlDB.Close()

	repo := NewOutboxRepository(logger, gormDB)
	event := &models.OutboxEvent{ID: 5, AggregateID: "activity-5", EventData: `{"x":5}`}

	// success
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "aggregate_id"=$1,"event_data"=$2,"published"=$3,"created_at"=$4,"published_at"=$5 WHERE "id" = $6`)).
		WithArgs(event.AggregateID, event.EventData, true, sqlmock.AnyArg(), sqlmock.AnyArg(), event.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.MarkOutboxEventAsPublished(context.Background(), event)
	assert.NoError(t, err)
	assert.True(t, event.Published)
	assert.NotNil(t, event.PublishedAt)
	assert.NoError(t, mock.ExpectationsWereMet())

	// db error
	event = &models.OutboxEvent{ID: 6, AggregateID: "activity-6", EventData: `{"x":6}`}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "aggregate_id"=$1,"event_data"=$2,"published"=$3,"created_at"=$4,"published_at"=$5 WHERE "id" = $6`)).
		WithArgs(event.AggregateID, event.EventData, true, sqlmock.AnyArg(), sqlmock.AnyArg(), event.ID).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err = repo.MarkOutboxEventAsPublished(context.Background(), event)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
