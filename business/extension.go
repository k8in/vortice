package business

import (
	"vortice/container"
	"vortice/object"
)

// TagBizKind is a constant string used to identify the business kind in tagged data.
const TagBizKind = "biz_kind"

var (
	// TagExtension is a predefined Tag used to mark components or properties as extensions with the key "biz_kind" and value "extension".
	TagExtension = object.NewTag(TagBizKind, "extension")
	// TagAbility is a predefined Tag used to mark components or properties as abilities with the key "biz_kind" and value "ability".
	TagAbility = object.NewTag(TagBizKind, "ability")
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
