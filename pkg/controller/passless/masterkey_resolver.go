package passless

import (
	"context"
	"errors"
	"math/rand"

	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/passless-operator/pkg/masterpassword"
	"github.com/wavesoftware/passless-operator/pkg/masterpassword/secrets"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type resolver struct {
	client client.Client
	ctx    context.Context
}

func (r *resolver) MasterKey() []byte {
	list := &corev1.SecretList{}
	err := r.client.List(r.ctx, list, client.InNamespace("kube-system"))
	ensure.NoError(err)
	secret, err := findDefaultToken(list)
	ensure.NoError(err)
	result := make([]byte, 0)
	for _, k := range secrets.Keys(secret) {
		bytes := secret.Data[k]
		result = append(result, bytes...)
	}
	source := rand.NewSource(42)
	rr := rand.New(source)
	rr.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

func findDefaultToken(list *corev1.SecretList) (*corev1.Secret, error) {
	for _, sec := range list.Items {
		if sec.Type == "kubernetes.io/service-account-token" &&
			sec.Annotations["kubernetes.io/service-account.name"] == "default" {
			return &sec, nil
		}
	}
	return &corev1.Secret{}, errors.New("can't find default token")
}

func newResolver(c client.Client) masterpassword.MasterKeyResolver {
	return &resolver{
		client: c,
		ctx:    context.TODO(),
	}
}
