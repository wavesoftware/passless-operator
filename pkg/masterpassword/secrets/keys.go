package secrets

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
)

// Keys returns a keys from secret in sorted slice
func Keys(sec *corev1.Secret) []string {
	m := sec.Data
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
