module github.com/tkeel-io/cli

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/briandowns/spinner v1.6.1
	github.com/dapr/cli v1.5.0
	github.com/fatih/color v1.13.0
	github.com/gocarina/gocsv v0.0.0-20210516172204-ca9e8a8ddea8
	github.com/gorilla/websocket v1.4.2
	github.com/gosuri/uitable v0.0.4
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/lib/pq v1.10.3 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/tkeel-io/kit v0.0.0-20211223050802-7dfccfe43fdb
	github.com/tkeel-io/tkeel v0.2.1-0.20220104063042-da4f2efb615a
	github.com/tkeel-io/tkeel-interface/openapi v0.0.0-20211223081012-25aaa61491ab
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	google.golang.org/protobuf v1.27.1
	helm.sh/helm/v3 v3.7.2
	k8s.io/api v0.23.1
	k8s.io/apimachinery v0.23.1
	k8s.io/cli-runtime v0.23.1
	k8s.io/client-go v0.23.1
	k8s.io/helm v2.16.10+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

	k8s.io/client => github.com/kubernetes-client/go v0.0.0-20190928040339-c757968c4c36
)
