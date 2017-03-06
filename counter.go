package metrics

import "sync/atomic"

// Counter holds an int64 value that can be incremented and decremented.
type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
	Snapshot() Counter

	Taggable
}

// GetOrRegisterCounter returns an existing Counter or constructs and registers
// a new StandardCounter.
func GetOrRegisterCounter(name string, r Registry) Counter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewCounter).(Counter)
}

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	if UseNilMetrics {
		return NilCounter{}
	}
	return &StandardCounter{count: 0}
}

// NewRegisteredCounter constructs and registers a new StandardCounter.
func NewRegisteredCounter(name string, r Registry) Counter {
	c := NewCounter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// CounterSnapshot is a read-only copy of another Counter.
type CounterSnapshot struct {
	count int64
	tags  map[string]string
}

// Clear panics.
func (c *CounterSnapshot) Clear() {
	panic("Clear called on a CounterSnapshot")
}

// Count returns the count at the time the snapshot was taken.
func (c *CounterSnapshot) Count() int64 { return c.count }

// Dec panics.
func (c *CounterSnapshot) Dec(int64) {
	panic("Dec called on a CounterSnapshot")
}

// Inc panics.
func (c *CounterSnapshot) Inc(int64) {
	panic("Inc called on a CounterSnapshot")
}

// Snapshot returns the snapshot.
func (c *CounterSnapshot) Snapshot() Counter { return c }

// AddTags panics.
func (c *CounterSnapshot) AddTags(tags map[string]string) {
	panic("AddTags called on a CounterSnapshot")
}

// GetTags returns the tags attached to this snapshot.
func (c *CounterSnapshot) GetTags() map[string]string { return c.tags }

// NilCounter is a no-op Counter.
type NilCounter struct{}

// Clear is a no-op.
func (NilCounter) Clear() {}

// Count is a no-op.
func (NilCounter) Count() int64 { return 0 }

// Dec is a no-op.
func (NilCounter) Dec(i int64) {}

// Inc is a no-op.
func (NilCounter) Inc(i int64) {}

// Snapshot is a no-op.
func (NilCounter) Snapshot() Counter { return NilCounter{} }

// AddTags is a no-op.
func (NilCounter) AddTags(tags map[string]string) {}

// GetTags is a no-op.
func (NilCounter) GetTags() map[string]string { return nil }

// StandardCounter is the standard implementation of a Counter and uses the
// sync/atomic package to manage a single int64 value.
type StandardCounter struct {
	count int64
	tags  map[string]string
}

// Clear sets the counter to zero.
func (c *StandardCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

// Count returns the current count.
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}

// Snapshot returns a read-only copy of the counter.
func (c *StandardCounter) Snapshot() Counter {
	if len(c.tags) == 0 {
		return &CounterSnapshot{count: c.Count()}
	}
	tagsCopy := map[string]string{}
	for k, v := range c.tags {
		tagsCopy[k] = v
	}
	return &CounterSnapshot{count: c.Count(), tags: tagsCopy}
}

// AddTags satisfies the Taggable interface and adds metric-level tags.
func (c *StandardCounter) AddTags(tags map[string]string) {
	if c.tags == nil {
		c.tags = tags
		return
	}
	for k, tag := range tags {
		c.tags[k] = tag
	}
}

// GetTags satisfies the Taggable interface.
func (c *StandardCounter) GetTags() map[string]string {
	return c.tags
}
