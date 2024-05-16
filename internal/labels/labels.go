package labels

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/cespare/xxhash/v2"
	"go.opentelemetry.io/otel/attribute"
)

type Labels []attribute.KeyValue

func (l Labels) Len() int           { return len(l) }
func (l Labels) Less(i, j int) bool { return l[i].Key < l[j].Key }
func (l Labels) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Labels) String() string {
	var kvs = make(map[string]string)
	for idx := range l {
		kvs[string(l[idx].Key)] = l[idx].Value.AsString()
	}
	data, _ := json.Marshal(kvs)
	return string(data)
}

func (l Labels) Hash() uint64 {
	sort.Sort(l)

	var builder strings.Builder
	for idx := range l {
		builder.WriteString(string(l[idx].Key))
		builder.WriteString("=")
		builder.WriteString(l[idx].Value.AsString())
		if idx != len(l)-1 {
			builder.WriteString(",")
		}
	}
	return xxhash.Sum64String(builder.String())
}

func (l Labels) Keys() []attribute.Key {
	keys := make([]attribute.Key, 0, len(l))
	for _, label := range l {
		keys = append(keys, label.Key)
	}
	return keys
}

func (l Labels) Map() map[string]string {
	m := make(map[string]string)
	for _, label := range l {
		m[string(label.Key)] = label.Value.AsString()
	}
	return m
}
