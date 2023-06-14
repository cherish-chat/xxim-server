package utils

type xMap struct {
}

var Map = &xMap{}

func (x *xMap) SS2SA(ss map[string]string) map[string]any {
	m := make(map[string]any)
	for k, v := range ss {
		m[k] = v
	}
	return m
}
