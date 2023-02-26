package main_test

import (
	"context"
	"testing"
	"time"

	"github.com/berquerant/firehose-test/grpctest"
	"github.com/berquerant/firehose-test/tempdir"
	"github.com/stretchr/testify/assert"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
)

func TestRun(t *testing.T) {
	dir := tempdir.New("grpctest")
	defer dir.Close()

	runner := grpctest.NewRunner(dir, grpctest.WithHealthWait(3*time.Second))
	defer runner.Close()
	assert.Nil(t, runner.Init(context.TODO()))

	client := grpchealth.NewHealthClient(runner.Conn)
	resp, err := client.Check(context.TODO(), &grpchealth.HealthCheckRequest{})
	assert.Nil(t, err)
	assert.Equal(t, grpchealth.HealthCheckResponse_SERVING, resp.Status)
}
