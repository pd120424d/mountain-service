package firestorex

import (
	"context"
	"errors"
)

// Direction for OrderBy
// Asc = 0, Desc = 1 for compactness
const (
	Asc Direction = iota
	Desc
)

type Direction int

// Done is returned by iterators when no more results are available
var Done = errors.New("firestorex: done")

type Client interface {
	Collection(name string) CollectionRef
}

type CollectionRef interface {
	Doc(id string) DocumentRef
	Where(field string, op string, value interface{}) Query
	OrderBy(field string, dir Direction) Query
	Limit(n int) Query
	Documents(ctx context.Context) DocumentIterator
}

type Query interface {
	Where(field string, op string, value interface{}) Query
	OrderBy(field string, dir Direction) Query
	StartAfter(v interface{}) Query
	Limit(n int) Query
	Documents(ctx context.Context) DocumentIterator
}

type DocumentIterator interface {
	Next() (DocumentSnapshot, error)
	Stop()
}

type DocumentSnapshot interface {
	DataTo(v interface{}) error
	ID() string
}

type DocumentRef interface {
	Set(ctx context.Context, data interface{}) (*WriteResult, error)
	Update(ctx context.Context, updates []Update) (*WriteResult, error)
	Delete(ctx context.Context) (*WriteResult, error)
}

type WriteResult struct{}

type Update struct {
	Path  string
	Value interface{}
}

func ServerTimestamp() interface{} { return ServerTimestampSentinel{} }

type ServerTimestampSentinel struct{}

func Increment(n int64) interface{} { return IncrementSentinel{N: n} }

type IncrementSentinel struct{ N int64 }

func IsDone(err error) bool { return err == Done }
