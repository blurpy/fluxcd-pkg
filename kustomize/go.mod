module github.com/fluxcd/pkg/kustomize

go 1.18

replace github.com/fluxcd/pkg/apis/kustomize => ../apis/kustomize

require (
	github.com/drone/envsubst v1.0.3
	github.com/fluxcd/pkg/apis/kustomize v0.7.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/onsi/gomega v1.24.1
	github.com/otiai10/copy v1.7.0
	k8s.io/api v0.25.4
	k8s.io/apiextensions-apiserver v0.25.4
	k8s.io/apimachinery v0.25.4
	sigs.k8s.io/controller-runtime v0.13.1
	sigs.k8s.io/kustomize/api v0.12.1
	sigs.k8s.io/kustomize/kyaml v0.13.9
	sigs.k8s.io/yaml v1.3.0
)

// Fix CVE-2022-28948
replace gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.10.0 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	go.starlark.net v0.0.0-20221028183056-acb66ad56dd2 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/oauth2 v0.2.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/term v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/time v0.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/client-go v0.25.4 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221110221610-a28e98eb7c70 // indirect
	k8s.io/utils v0.0.0-20221108210102-8e77b1f39fe2 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
