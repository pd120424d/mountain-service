package repositories

import (
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
		WillReturnRows(sqlmock.NewRows([]string{"id", "event_type", "aggregate_id", "event_data", "published", "created_at", "published_at"}).
			AddRow(1, "activity.created", "activity-1", `{"x":1}`, false, time.Now(), nil))

	events, err := repo.GetUnpublishedEvents(100)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	assert.NoError(t, mock.ExpectationsWereMet())

	// db error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "outbox_events" WHERE published = $1 ORDER BY created_at ASC LIMIT $2`)).
		WithArgs(false, 50).
		WillReturnError(assert.AnError)

	_, err = repo.GetUnpublishedEvents(50)
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

	err := repo.MarkAsPublished(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// db error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "published"=$1,"published_at"=$2 WHERE id = $3`)).
		WithArgs(true, sqlmock.AnyArg(), 2).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err = repo.MarkAsPublished(2)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOutboxRepository_MarkOutboxEventAsPublished(t *testing.T) {
	logger := utils.NewTestLogger()
	gormDB, mock, sqlDB := newGormWithSQLMock(t)
	defer sqlDB.Close()

	repo := NewOutboxRepository(logger, gormDB)
	event := &models.OutboxEvent{ID: 5, EventType: "activity.created", AggregateID: "activity-5", EventData: `{"x":5}`}

	// success
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "event_type"=$1,"aggregate_id"=$2,"event_data"=$3,"published"=$4,"created_at"=$5,"published_at"=$6 WHERE "id" = $7`)).
		WithArgs(event.EventType, event.AggregateID, event.EventData, true, sqlmock.AnyArg(), sqlmock.AnyArg(), event.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.MarkOutboxEventAsPublished(event)
	assert.NoError(t, err)
	assert.True(t, event.Published)
	assert.NotNil(t, event.PublishedAt)
	assert.NoError(t, mock.ExpectationsWereMet())

	// db error
	event = &models.OutboxEvent{ID: 6, EventType: "activity.created", AggregateID: "activity-6", EventData: `{"x":6}`}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "event_type"=$1,"aggregate_id"=$2,"event_data"=$3,"published"=$4,"created_at"=$5,"published_at"=$6 WHERE "id" = $7`)).
		WithArgs(event.EventType, event.AggregateID, event.EventData, true, sqlmock.AnyArg(), sqlmock.AnyArg(), event.ID).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err = repo.MarkOutboxEventAsPublished(event)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
