package mocks

//go:generate bazelisk run --config=mayberemote //:mockery   -- --name Transport  --srcpkg=go.skia.org/infra/perf/go/notify --output ${PWD}

//go:generate bazelisk run --config=mayberemote //:mockery   -- --name Notifier  --srcpkg=go.skia.org/infra/perf/go/notify --output ${PWD}
