package firestoretest

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/firestorex"
)

type Fake struct {
	mu          sync.RWMutex
	collections map[string][]doc
}

type doc struct {
	id   string
	data map[string]interface{}
}

type coll struct {
	f          *Fake
	key        string
	flt        []filter
	ord        *order
	startAfter interface{}
	lim        int
}

type filter struct {
	field, op string
	value     interface{}
}

type order struct {
	field string
	desc  bool
}

type iter struct {
	items []doc
	idx   int
}

type docRef struct {
	f       *Fake
	key, id string
}

type writeResult struct{}

func NewFake() *Fake { return &Fake{collections: map[string][]doc{}} }

func (f *Fake) WithCollection(name string, docs []map[string]interface{}) *Fake {
	f.mu.Lock()
	defer f.mu.Unlock()
	arr := make([]doc, 0, len(docs))
	for i, m := range docs {
		id := fmt.Sprintf("%d", i+1)
		arr = append(arr, doc{id: id, data: cloneMap(m)})
	}
	f.collections[name] = arr
	return f
}

func (f *Fake) Collection(name string) firestorex.CollectionRef { return &coll{f: f, key: name} }

func (c *coll) Doc(id string) firestorex.DocumentRef { return &docRef{f: c.f, key: c.key, id: id} }

func (c *coll) Where(field, op string, value interface{}) firestorex.Query {
	c.flt = append(c.flt, filter{field: field, op: op, value: value})
	return c
}

func (c *coll) OrderBy(field string, dir firestorex.Direction) firestorex.Query {
	c.ord = &order{field: field, desc: dir == firestorex.Desc}
	return c
}

func (c *coll) StartAfter(v interface{}) firestorex.Query { c.startAfter = v; return c }

func (c *coll) Limit(n int) firestorex.Query { c.lim = n; return c }

func (c *coll) Documents(ctx context.Context) firestorex.DocumentIterator {
	c.f.mu.RLock()
	defer c.f.mu.RUnlock()
	src := c.f.collections[c.key]
	items := make([]doc, 0, len(src))
	// filter
	for _, d := range src {
		if matchFilters(d.data, c.flt) {
			items = append(items, d)
		}
	}
	// order
	if c.ord != nil {
		sort.Slice(items, func(i, j int) bool {
			vi := items[i].data[c.ord.field]
			vj := items[j].data[c.ord.field]
			less := fmt.Sprint(vi) < fmt.Sprint(vj)
			if c.ord.desc {
				return !less
			}
			return less
		})
	}
	// startAfter
	if c.ord != nil && c.startAfter != nil {
		cursor := fmt.Sprint(c.startAfter)
		filtered := make([]doc, 0, len(items))
		for _, d := range items {
			val := fmt.Sprint(d.data[c.ord.field])
			if c.ord.desc {
				if val < cursor {
					filtered = append(filtered, d)
				}
			} else {
				if val > cursor {
					filtered = append(filtered, d)
				}
			}
		}
		items = filtered
	}

	// limit
	if c.lim > 0 && c.lim < len(items) {
		items = items[:c.lim]
	}
	return &iter{items: items, idx: 0}
}

func matchFilters(m map[string]interface{}, flts []filter) bool {
	for _, f := range flts {
		if f.op != "==" {
			return false
		}
		if !reflect.DeepEqual(m[f.field], f.value) {
			return false
		}
	}
	return true
}

func (it *iter) Next() (firestorex.DocumentSnapshot, error) {
	if it.idx >= len(it.items) {
		return nil, firestorex.Done
	}
	d := it.items[it.idx]
	it.idx++
	return snap{doc: d}, nil
}
func (it *iter) Stop() {}

type snap struct{ doc doc }

func (s snap) DataTo(v interface{}) error { return firestorex.SnapshotDataTo(s.doc.data, v) }
func (s snap) ID() string                 { return s.doc.id }

func (r *docRef) Get(ctx context.Context) (firestorex.DocumentSnapshot, error) {
	r.f.mu.RLock()
	defer r.f.mu.RUnlock()
	arr := r.f.collections[r.key]
	for _, d := range arr {
		if d.id == r.id {
			return snap{doc: d}, nil
		}
	}
	return nil, errors.New("firestoretest: doc not found")
}

func (r *docRef) Set(ctx context.Context, data interface{}) (*firestorex.WriteResult, error) {
	m, ok := toMap(data)
	if !ok {
		return nil, errors.New("firestoretest: set expects map or struct")
	}
	r.f.mu.Lock()
	defer r.f.mu.Unlock()
	arr := r.f.collections[r.key]
	found := false
	for i := range arr {
		if arr[i].id == r.id {
			arr[i].data = applySentinels(m)
			found = true
			break
		}
	}
	if !found {
		arr = append(arr, doc{id: r.id, data: applySentinels(m)})
	}
	r.f.collections[r.key] = arr
	return &firestorex.WriteResult{}, nil
}

func (r *docRef) Update(ctx context.Context, updates []firestorex.Update) (*firestorex.WriteResult, error) {
	r.f.mu.Lock()
	defer r.f.mu.Unlock()
	arr := r.f.collections[r.key]
	for i := range arr {
		if arr[i].id == r.id {
			for _, u := range updates {
				arr[i].data[u.Path] = applyValue(arr[i].data[u.Path], u.Value)
			}
			r.f.collections[r.key] = arr
			return &firestorex.WriteResult{}, nil
		}
	}
	return nil, errors.New("firestoretest: doc not found")
}

func (r *docRef) Delete(ctx context.Context) (*firestorex.WriteResult, error) {
	r.f.mu.Lock()
	defer r.f.mu.Unlock()
	arr := r.f.collections[r.key]
	for i := range arr {
		if arr[i].id == r.id {
			arr = append(arr[:i], arr[i+1:]...)
			r.f.collections[r.key] = arr
			return &firestorex.WriteResult{}, nil
		}
	}
	return nil, nil
}

func applySentinels(in map[string]interface{}) map[string]interface{} {
	out := cloneMap(in)
	for k, v := range out {
		out[k] = applyValue(out[k], v)
	}
	return out
}

func applyValue(prev, v interface{}) interface{} {
	switch t := v.(type) {
	case firestorex.IncrementSentinel:
		if prev == nil {
			return t.N
		}
		if n, ok := prev.(int64); ok {
			return n + t.N
		}
		if n, ok := prev.(int); ok {
			return int64(n) + t.N
		}
		return t.N
	case firestorex.ServerTimestampSentinel:
		return time.Now().UTC().Format(time.RFC3339)
	default:
		return v
	}
}

func toMap(v interface{}) (map[string]interface{}, bool) {
	switch t := v.(type) {
	case map[string]interface{}:
		return t, true
	default:
		return structToMap(v)
	}
}

func structToMap(v interface{}) (map[string]interface{}, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, false
	}
	tp := rv.Type()
	m := map[string]interface{}{}
	for i := 0; i < rv.NumField(); i++ {
		field := tp.Field(i)
		tag := field.Tag.Get("firestore")
		if tag == "" {
			tag = field.Name
		}
		m[tag] = rv.Field(i).Interface()
	}
	return m, true
}

func cloneMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
