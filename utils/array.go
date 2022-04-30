package utils

func FilterByGeneric[V any](elems []V, predicate func(V) bool) []V {
	var r []V
	for _, e := range elems {
		if predicate(e) {
			r = append(r, e)
		}
	}
	return r
}

func FilterOneByGeneric[V any](elems []V, predicate func(V) bool) V {
	var r V
	for _, e := range elems {
		if predicate(e) {
			return e
		}
	}
	return r
}
