package v1alpha1

import (
	"encoding/base64"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wavesoftware/passless-operator/pkg/masterpassword"
)

// CreateSecret will created a secret that corresponds to
func (in *PassLess) CreateSecret(generator masterpassword.Generator) *corev1.Secret {
	data := in.createData(generator)
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:        in.Name,
			Namespace:   in.Namespace,
			Labels:      in.copyMap(in.Labels),
			Annotations: in.copyMap(in.Annotations),
		},
		Data: data,
		Type: corev1.SecretTypeOpaque,
	}
	return secret
}

func (in *PassLess) createData(generator masterpassword.Generator) map[string][]byte {
	data := make(map[string][]byte, len(in.Spec))
	for name, entry := range in.Spec {
		secret := generator.Generate(in.identity(name), entry.Scope, entry.Version, entry.Length)
		dst := base64.StdEncoding.EncodeToString([]byte(secret))
		data[name] = []byte(dst)
	}
	return data
}

func (in *PassLess) identity(name string) string {
	return in.Namespace + "/" + name
}

func (in *PassLess) copyMap(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
