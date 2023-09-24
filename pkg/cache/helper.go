package cache

import "strings"

// ResolveKey is just a simple method to concat keys with key separator
func ResolveKey(keys ...string) string {
	return strings.Join(keys, KeySeparator)
}

type Resolver struct {
	Prefix []string
}

func (r *Resolver) ResolveKey(keys ...string) string {
	return ResolveKey(append(r.Prefix, keys...)...)
}