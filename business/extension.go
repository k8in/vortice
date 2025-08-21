package business

import (
	"vortice/container"
	"vortice/object"
)

// TagBizKindKey is a constant string used to identify the business kind in tagged data.
const TagBizKindKey = "kind"

var (
	// TagExtensionKind is a predefined Tag used to mark components or properties as extensions with the key "biz_kind" and value "extension".
	TagExtensionKind = object.NewTag(TagBizKindKey, "extension")
	// TagAbilityKind is a predefined Tag used to mark components or properties as abilities with the key "biz_kind" and value "ability".
	TagAbilityKind = object.NewTag(TagBizKindKey, "ability")
)

type (
	// Extension is a marker interface used to signify that a type can be extended or enriched with additional behaviors.
	Extension interface{}
	// ExtensionObject represents an object that can be managed within a container and extended with additional behaviors.
	ExtensionObject struct {
		Extension
		container.Object
	}
)

// newExtensionObject creates a new ExtensionObject with the given container.Object, allowing for extended behaviors.
func newExtensionObject(obj container.Object) *ExtensionObject {
	return &ExtensionObject{Object: obj}
}
