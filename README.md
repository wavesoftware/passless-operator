# PassLess Operator for Kubernetes

[![Build Status](https://travis-ci.com/wavesoftware/passless-operator.svg?branch=master)](https://travis-ci.com/wavesoftware/passless-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/wavesoftware/passless-operator)](https://goreportcard.com/report/github.com/wavesoftware/passless-operator)

PassLess Operator implements a concept of secret management without the need of
credentials storage for usage in Kubernetes. It utilize 
[master password algorithm](https://en.wikipedia.org/wiki/Master_Password) to 
achieve that.

Using PassLess you can avoid storing passwords in your GitOps workflow or the 
need of secure data source like LDAP, or HashiCorp's Vault.

## Installation

Deploy operator with:

```bash
kubectl apply -f https://github.com/wavesoftware/passless-operator/releases/download/v0.2.0/passless.yaml
```

## Usage

This is an example passless resource:

```yaml
---
apiVersion: wavesoftware.pl/v1alpha1
kind: PassLess
metadata:
  name: example
spec:
  db-password:
    version: 1
    scope: alnum
    length: 16
```

It define a Kubernetes secret's template. If you have Passless operator running,
and you create such a resource, operator will create a secret for you. It will 
also, update it if you change passless resource!

```
$ kubectl get passless
NAME      AGE
example   12s

$ kubectl get secrets
NAME      TYPE     DATA   AGE   
example   Opaque   1      10s

$ kubectl get secret example -o jsonpath='{.data.db-password}' | base64 -d
eoXdlNHgrtaxoO34
```

Each PassLess specification element defines a secret element to be created, and 
each password generation can be influenced by providing an options:

### Parameters

* `version` - A sequential password number. Changing the password should be 
   done by advancing this number. Default value is `1`.
 * `scope` - A definition of scope that the password will be generated from. It 
   may be one of (defaults to `alnum`): 
    * `num` for numeric passwords,
    * `alpha` for alphabet based passwords, both big and small caps,
    * `alnum` for alphanumeric passwords, both big and small caps,
    * `human` for letters and numbers that are easy to distinguish by human,
    * `keys` for passwords that can be typed by keyboard, so letters, and 
      numbers, and special characters,
    * `utf8` these passwords contain utf-8 characters, so also a characters 
      that aren't easy to type by keyboard,
    * `list:` followed by list of chars that might be used. Ex.: 
      `list:abcdef1234567890!$`,
 * `length` - A length of password to be generated in number of signs. Default 
   value is `16`.

### Password generator configuration

PassLess created secrets use Master Password algorithm. The master key is derived  
from `default-token` secret form `kube-system` namespace. Site key uses the name
of passless resource and namespace in which it is created.

Above means that the password for the same parameters given will the same, but 
different if created in other namespace or with other name. All passwords will 
change if `default-token` secret form `kube-system` namespace is changed.

For now, this behavior isn't configurable, but it's good idea for future 
features.

## Contributing

Contributions are welcome!

To contribute, follow the standard [git flow](http://danielkummer.github.io/git-flow-cheatsheet) of:

1. Fork it
1. Create your feature branch (`git checkout -b feature/my-new-feature`)
1. Commit your changes (`git commit -am 'Add some feature'`)
1. Push to the branch (`git push origin feature/my-new-feature`)
1. Create new Pull Request

Even if you can't contribute code, if you have an idea for an improvement 
please open an [issue](https://github.com/wavesoftware/passless-operator/issues).
