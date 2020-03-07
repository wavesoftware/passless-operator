package secrets

import (
	"hash/crc32"

	corev1 "k8s.io/api/core/v1"
)

// Hashcode will calculate a secret's hash code
func Hashcode(secret *corev1.Secret) int {
	code := 31
	keys := Keys(secret)
	for _, k := range keys {
		next := crc32.ChecksumIEEE(secret.Data[k])
		code = code*int(next) + 17
	}
	return code
}
