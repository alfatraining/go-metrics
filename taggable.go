package metrics

// Taggable provides the interface for metrics to have metric-level tags.
type Taggable interface {
	AddTags(tags map[string]string)
	GetTags() map[string]string
}
