module github.com/wavesoftware/passless-operator

go 1.14

require (
	github.com/fatih/color v1.7.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/magefile/mage v1.9.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	github.com/wavesoftware/go-ensure v1.0.0
	golang.org/x/crypto v0.0.0-20200414173820-0848c9571904
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)
