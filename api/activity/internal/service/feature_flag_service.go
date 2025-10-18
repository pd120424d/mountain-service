package service

import (
	"context"

	"github.com/pd120424d/mountain-service/api/shared/firestorex"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// FeatureFlagService provides access to runtime-togglable flags.
// Currently there is only one flag: whether to use Postgres for listing activities.
// Default is false (use Firestore read model).

type FeatureFlagService interface {
	GetUsePostgresForActivities(ctx context.Context) (bool, error)
	SetUsePostgresForActivities(ctx context.Context, value bool, actor string) error
}

// In-memory fallback implementation (per-pod, non-persistent)
type inMemoryFeatureFlag struct {
	logger utils.Logger
	value  bool
}

func NewInMemoryFeatureFlag(logger utils.Logger, defaultValue bool) FeatureFlagService {
	return &inMemoryFeatureFlag{logger: logger.WithName("featureFlag.inmem"), value: defaultValue}
}

func (f *inMemoryFeatureFlag) GetUsePostgresForActivities(ctx context.Context) (bool, error) {
	return f.value, nil
}

func (f *inMemoryFeatureFlag) SetUsePostgresForActivities(ctx context.Context, value bool, actor string) error {
	f.value = value
	f.logger.Infof("Feature flag updated (in-mem): use_postgres_for_activities=%v by=%s", value, actor)
	return nil
}

type firestoreFeatureFlag struct {
	client firestorex.Client
	logger utils.Logger
}

type activitySourceDoc struct {
	UsePostgres bool   `firestore:"use_postgres"`
	UpdatedBy   string `firestore:"updated_by"`
}

func NewFirestoreFeatureFlag(client firestorex.Client, logger utils.Logger) FeatureFlagService {
	return &firestoreFeatureFlag{client: client, logger: logger.WithName("featureFlag.firestore")}
}

func (f *firestoreFeatureFlag) GetUsePostgresForActivities(ctx context.Context) (bool, error) {
	if f.client == nil {
		return false, nil
	}
	doc := f.client.Collection("config").Doc("activity_source")
	ds, err := doc.Get(ctx)
	if err != nil {
		// If document does not exist or any error, default to false (Firestore)
		f.logger.Warnf("Feature flag read failed, defaulting to false: %v", err)
		return false, nil
	}
	var cfg activitySourceDoc
	if err := ds.DataTo(&cfg); err != nil {
		f.logger.Warnf("Feature flag decode failed, defaulting to false: %v", err)
		return false, nil
	}
	return cfg.UsePostgres, nil
}

func (f *firestoreFeatureFlag) SetUsePostgresForActivities(ctx context.Context, value bool, actor string) error {
	if f.client == nil {
		return nil
	}
	doc := f.client.Collection("config").Doc("activity_source")
	_, err := doc.Set(ctx, map[string]interface{}{
		"use_postgres": value,
		"updated_by":   actor,
		"updated_at":   firestorex.ServerTimestamp(),
	})
	if err != nil {
		f.logger.Errorf("Failed to update feature flag in Firestore: %v", err)
		return err
	}
	f.logger.Infof("Feature flag updated (firestore): use_postgres_for_activities=%v by=%s", value, actor)
	return nil
}
