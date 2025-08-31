package googleadapter

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/pd120424d/mountain-service/api/shared/firestorex"
)

type clientAdapter struct{ c *firestore.Client }

type collectionRefAdapter struct{ c *firestore.CollectionRef }

type queryAdapter struct{ q firestore.Query }

type iteratorAdapter struct{ it *firestore.DocumentIterator }

type docRefAdapter struct{ d *firestore.DocumentRef }

type snapshotAdapter struct{ snap *firestore.DocumentSnapshot }

type writeResultAdapter struct{ r *firestore.WriteResult }

type updateAdapter struct{ u firestore.Update }

func NewClientAdapter(c *firestore.Client) firestorex.Client { return &clientAdapter{c: c} }

func (a *clientAdapter) Collection(name string) firestorex.CollectionRef {
	return &collectionRefAdapter{c: a.c.Collection(name)}
}

func (a *collectionRefAdapter) Doc(id string) firestorex.DocumentRef {
	return &docRefAdapter{d: a.c.Doc(id)}
}

func (a *collectionRefAdapter) Where(field, op string, value interface{}) firestorex.Query {
	return &queryAdapter{q: a.c.Where(field, op, value)}
}

func (a *collectionRefAdapter) OrderBy(field string, dir firestorex.Direction) firestorex.Query {
	if dir == firestorex.Desc {
		return &queryAdapter{q: a.c.OrderBy(field, firestore.Desc)}
	}
	return &queryAdapter{q: a.c.OrderBy(field, firestore.Asc)}
}

func (a *collectionRefAdapter) Limit(n int) firestorex.Query { return &queryAdapter{q: a.c.Limit(n)} }

func (a *collectionRefAdapter) Documents(ctx context.Context) firestorex.DocumentIterator {
	it := a.c.Documents(ctx)
	return &iteratorAdapter{it: it}
}

// Query methods
func (a *queryAdapter) Where(field, op string, value interface{}) firestorex.Query {
	return &queryAdapter{q: a.q.Where(field, op, value)}
}
func (a *queryAdapter) OrderBy(field string, dir firestorex.Direction) firestorex.Query {
	if dir == firestorex.Desc {
		return &queryAdapter{q: a.q.OrderBy(field, firestore.Desc)}
	}
	return &queryAdapter{q: a.q.OrderBy(field, firestore.Asc)}
}
func (a *queryAdapter) Limit(n int) firestorex.Query { return &queryAdapter{q: a.q.Limit(n)} }
func (a *queryAdapter) Documents(ctx context.Context) firestorex.DocumentIterator {
	it := a.q.Documents(ctx)
	return &iteratorAdapter{it: it}
}

// Iterator
func (a *iteratorAdapter) Next() (firestorex.DocumentSnapshot, error) {
	ds, err := a.it.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, firestorex.Done
		}
		return nil, err
	}
	return snapshotWrapper{ds: ds}, nil
}

func (a *iteratorAdapter) Stop() { a.it.Stop() }

type snapshotWrapper struct{ ds *firestore.DocumentSnapshot }

func (s snapshotWrapper) DataTo(v interface{}) error { return s.ds.DataTo(v) }
func (s snapshotWrapper) ID() string                 { return s.ds.Ref.ID }

func iteratorDone() error {
	return fmt.Errorf("iterator done")
}

// Doc ref
func (a *docRefAdapter) Set(ctx context.Context, data interface{}) (*firestorex.WriteResult, error) {
	_, err := a.d.Set(ctx, data)
	if err != nil {
		return nil, err
	}
	return &firestorex.WriteResult{}, nil
}
func (a *docRefAdapter) Update(ctx context.Context, updates []firestorex.Update) (*firestorex.WriteResult, error) {
	gus := make([]firestore.Update, 0, len(updates))
	for _, u := range updates {
		gus = append(gus, firestore.Update{Path: u.Path, Value: adaptValue(u.Value)})
	}
	_, err := a.d.Update(ctx, gus)
	if err != nil {
		return nil, err
	}
	return &firestorex.WriteResult{}, nil
}
func (a *docRefAdapter) Delete(ctx context.Context) (*firestorex.WriteResult, error) {
	_, err := a.d.Delete(ctx)
	if err != nil {
		return nil, err
	}
	return &firestorex.WriteResult{}, nil
}

func adaptValue(v interface{}) interface{} {
	switch t := v.(type) {
	case firestorex.IncrementSentinel:
		return firestore.Increment(t.N)
	case firestorex.ServerTimestampSentinel:
		return firestore.ServerTimestamp
	default:
		return v
	}
}
