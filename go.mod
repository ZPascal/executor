module code.cloudfoundry.org/executor

go 1.27.0

replace code.cloudfoundry.org/bbs/models => ../bbs/models

replace code.cloudfoundry.org/garden => ../garden

require (
	code.cloudfoundry.org/archiver v0.87.0
	code.cloudfoundry.org/bbs/models v0.0.0-00010101000000-000000000000
	code.cloudfoundry.org/bytefmt v0.89.0
	code.cloudfoundry.org/cacheddownloader v0.0.0-20250312193827-23c030d5e4f3
	code.cloudfoundry.org/clock v1.87.0
	code.cloudfoundry.org/diego-logging-client v0.124.0
	code.cloudfoundry.org/durationjson v0.89.0
	code.cloudfoundry.org/eventhub v0.89.0
	code.cloudfoundry.org/garden v0.0.0-00010101000000-000000000000
	code.cloudfoundry.org/go-loggregator/v9 v9.2.1
	code.cloudfoundry.org/lager/v3 v3.86.0
	code.cloudfoundry.org/tlsconfig v0.66.0
	code.cloudfoundry.org/volman v0.0.0-20250910193608-1cc72f1031b7
	code.cloudfoundry.org/workpool v0.0.0-20250911194158-1489753f182e
	github.com/envoyproxy/go-control-plane/envoy v1.39.0
	github.com/fsnotify/fsnotify v1.10.1
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.4
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/onsi/ginkgo/v2 v2.32.2
	github.com/onsi/gomega v1.43.0
	github.com/tedsuo/ifrit v0.0.0-20260908181113-dd353a7daa27
	golang.org/x/time v0.16.0
	google.golang.org/protobuf v1.36.12
)

require (
	cel.dev/expr v0.25.2 // indirect
	code.cloudfoundry.org/bbs/encryption v1.11.0 // indirect
	code.cloudfoundry.org/bbs/format v1.10.0 // indirect
	code.cloudfoundry.org/cfhttp/v2 v2.93.0 // indirect
	code.cloudfoundry.org/dockerdriver v0.106.0 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20260831145205-e8366a756183 // indirect
	code.cloudfoundry.org/goshims v0.113.0 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/bmizerany/pat v0.0.0-20210406213842-e4b6760bdd6f // indirect
	github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 // indirect
	github.com/cyphar/filepath-securejoin v0.7.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20260906184651-6331bc6350fe // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/tedsuo/rata v1.0.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/tools v0.50.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260904194346-d0f1323225a4 // indirect
	google.golang.org/grpc v1.83.2 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
