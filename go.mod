module github.com/tilebox/tilebox-go

go 1.23.1

require (
	connectrpc.com/connect v1.18.1
	github.com/avast/retry-go/v4 v4.6.0
	github.com/google/uuid v1.6.0
	github.com/hashicorp/go-retryablehttp v0.7.7
	github.com/lmittmann/tint v1.0.7
	github.com/oklog/ulid/v2 v2.1.0
	github.com/paulmach/orb v0.11.1
	github.com/samber/lo v1.47.0
	github.com/samber/slog-multi v1.2.4
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/contrib/bridges/otelslog v0.10.0
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.11.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.34.0
	go.opentelemetry.io/otel/metric v1.35.0
	go.opentelemetry.io/otel/sdk v1.35.0
	go.opentelemetry.io/otel/sdk/log v0.11.0
	go.opentelemetry.io/otel/trace v1.35.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
	pgregory.net/rapid v1.1.0
)

require (
	buf.build/gen/go/bufbuild/bufplugin/protocolbuffers/go v1.36.5-20250121211742-6d880cc6cc8d.1 // indirect
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.5-20250307204501-0409229c3780.1 // indirect
	buf.build/gen/go/bufbuild/registry/connectrpc/go v1.18.1-20250116203702-1c024d64352b.1 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.36.5-20250116203702-1c024d64352b.1 // indirect
	buf.build/gen/go/pluginrpc/pluginrpc/protocolbuffers/go v1.36.4-20241007202033-cf42259fcbfc.1 // indirect
	buf.build/go/bufplugin v0.8.0 // indirect
	buf.build/go/protoyaml v0.3.1 // indirect
	buf.build/go/spdx v0.2.0 // indirect
	cel.dev/expr v0.21.2 // indirect
	connectrpc.com/otelconnect v0.7.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/bufbuild/buf v1.51.0 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/bufbuild/protoplugin v0.0.0-20250106231243-3a819552c9d9 // indirect
	github.com/bufbuild/protovalidate-go v0.9.3-0.20250317160558-38a17488914d // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.16.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/cli v27.5.1+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v28.0.2+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.8.2 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/felixge/fgprof v0.9.5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/cel-go v0.24.1 // indirect
	github.com/google/go-containerregistry v0.20.3 // indirect
	github.com/google/pprof v0.0.0-20250302191652-9094ed2288e7 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jdx/go-netrc v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/mount v0.3.4 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/moby/sys/reexec v0.1.0 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.3.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/onsi/ginkgo/v2 v2.22.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/profile v1.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.50.1 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.4.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tetratelabs/wazero v1.9.0 // indirect
	github.com/vbatts/tar-split v0.12.1 // indirect
	go.lsp.dev/jsonrpc2 v0.10.0 // indirect
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2 // indirect
	go.lsp.dev/protocol v0.12.0 // indirect
	go.lsp.dev/uri v0.3.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.59.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.0 // indirect
	go.opentelemetry.io/otel/log v0.11.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.35.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.uber.org/zap/exp v0.3.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250228200357-dead58393ab7 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250224174004-546df14abb99 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250224174004-546df14abb99 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
	pluginrpc.com/pluginrpc v0.5.0 // indirect
)

tool github.com/bufbuild/buf/cmd/buf
