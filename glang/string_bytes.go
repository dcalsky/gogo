package glang

import "unsafe"

func B2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func S2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
