package passless

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/passless-operator/pkg/masterpassword"
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
	crt := secret.Data["ca.crt"]
	var dst []byte
	_, err = base64.StdEncoding.Decode(dst, crt)
	ensure.NoError(err)
	return dst
}

func findDefaultToken(list *corev1.SecretList) (corev1.Secret, error) {
	for _, sec := range list.Items {
		if sec.Annotations["kubernetes.io/service-account.name"] == "default" {
			return sec, nil
		}
	}
	return corev1.Secret{}, errors.New("can't find default token")
}

func newResolver(c client.Client) masterpassword.MasterKeyResolver {
	return &resolver{
		client: c,
		ctx:    context.TODO(),
	}
}
