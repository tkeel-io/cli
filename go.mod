module github.com/tkeel-io/cli

go 1.16

require (
	github.com/briandowns/spinner v1.6.1
	github.com/dapr/cli v1.4.0
	github.com/dapr/dapr v1.3.1 // indirect
	github.com/fatih/color v1.10.0
	github.com/gocarina/gocsv v0.0.0-20210516172204-ca9e8a8ddea8
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	helm.sh/helm/v3 v3.5.3
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/cli-runtime v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/helm v2.16.10+incompatible
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

	k8s.io/client => github.com/kubernetes-client/go v0.0.0-20190928040339-c757968c4c36
)
