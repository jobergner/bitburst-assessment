package get

import "assessment/pkg/object"

type Getter interface {
	Get(int) (object.Object, error)
}
