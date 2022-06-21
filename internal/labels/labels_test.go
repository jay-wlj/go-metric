package labels

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func Test_Labels(t *testing.T) {
	labels := Labels{
		{Key: "ip", Value: attribute.StringValue("1.1.1.1")},
		{Key: "zone", Value: attribute.StringValue("sh")},
		{Key: "type", Value: attribute.StringValue("counter")},
		{Key: "a", Value: attribute.StringValue("b")},
	}
	beforeHash := labels.Hash()
	afterHash := labels.Hash()
	if labels[0].Key != "a" {
		t.Error("sort doesn't work well")
	}
	if labels[3].Key != "zone" {
		t.Error("sort doesn't work well")
	}
	if beforeHash != afterHash {
		t.Error("hash changes")
	}
}
