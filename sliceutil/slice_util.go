package sliceutil

func Filter[S ~[]E, E any](x S, filter func(a E) bool) S {
	var res []E
	for _, v := range x {
		if filter(v) {
			res = append(res, v)
		}
	}
	return res
}

func DistinctFunc[S ~[]E, E any](x S, cmp func(a, b E) bool) S {
	var result []E
	for _, item := range x {
		duplicate := false
		for j := range result {
			if cmp(item, result[j]) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, item)
		}
	}
	return result
}

func Map[S ~[]E, E any, R any](x S, fun func(a E) R) []R {
	var r []R
	for _, v := range x {
		r = append(r, fun(v))
	}
	return r
}
