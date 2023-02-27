package grpctest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/berquerant/firehose-test/command"
	"github.com/berquerant/firehose-test/tempdir"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

//go:generate go run github.com/berquerant/goconfig@v0.1.0 -field "Executable string|Port int|HealthWait time.Duration" -option -configOption Option -output grpc_config_generated.go

type Runner struct {
	Dir        *tempdir.Dir
	Conf       *Config
	Conn       *grpc.ClientConn
	RunCommand *exec.Cmd
}

// NewRunner returns a test runner to help combained grpc test.
//
// Options:
//
// - WithExecutable: executable name to be built (default: server)
// - WithHealthWait: wait time for health check (default: 1s)
// - WithPort: port number that the server listens on (default: 17001)
func NewRunner(dir *tempdir.Dir, opt ...Option) *Runner {
	conf := NewConfigBuilder().
		Executable("server").
		Port(17001).
		HealthWait(time.Second).
		Build()
	conf.Apply(opt...)

	return &Runner{
		Dir:  dir,
		Conf: conf,
	}
}

// Init prepares resources for test.
//
// Init does the following:
//
// 1. Build a executable of grpc server
// 2. Set PORT env var
// 3. Start the server
// 4. Wait the server is ready
// 5. Create a connection to server
func (r *Runner) Init(ctx context.Context) error {
	// expect that test code is in the same dir as the target
	executable := r.Dir.Join(r.Conf.Executable.Get())
	if err := command.New("go", "build", "-o", executable).Run(); err != nil {
		return fmt.Errorf("Failed to build %s: %w",
			r.Conf.Executable.Get(),
			err,
		)
	}
	// grpc server of firehose read a port number from env
	os.Setenv("PORT", fmt.Sprint(r.Conf.Port.Get()))
	// start grpc server
	runCommand := command.New(executable)
	if err := runCommand.Start(); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}
	r.RunCommand = runCommand
	// check server is ready
	healthCtx, cancel := context.WithTimeout(ctx, r.Conf.HealthWait.Get())
	defer cancel()
	if err := WaitServerReady(healthCtx, r.Conf.Port.Get()); err != nil {
		return fmt.Errorf("Server is not ready: %w", err)
	}
	// create conn
	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%d", r.Conf.Port.Get()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("Failed to create conn: %w", err)
	}
	r.Conn = conn
	return nil
}

// Close sweeps resources for test.
//
// Close does the following:
//
// 1. Close the connection to server
// 2. Send sigint to server process
// 3. Wait the server is stopped
func (r *Runner) Close() error {
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			return fmt.Errorf("Failed to close conn: %w", err)
		}
	}
	if r.RunCommand != nil {
		// grpc server of firehose gracefully stop when interrupted
		if err := r.RunCommand.Process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("Failed to interrupt %s: %w", r.Conf.Executable.Get(), err)
		}
		if err := r.RunCommand.Wait(); err != nil {
			return fmt.Errorf("Failed to close %s: %w", r.Conf.Executable.Get(), err)
		}
	}
	return nil
}

func WaitServerReady(ctx context.Context, port int) error {
	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("WaitServerReady got conn error: %w", err)
	}
	defer conn.Close()

	var (
		client = healthpb.NewHealthClient(conn)
		req    = &healthpb.HealthCheckRequest{}
	)

	for {
		time.Sleep(200 * time.Millisecond)
		resp, err := client.Check(ctx, req)
		switch {
		case err == nil:
			if resp.Status == healthpb.HealthCheckResponse_SERVING {
				return nil
			}
			log.Printf("WaitServerReady port %d: %v", port, resp.Status)
		case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("WaitServerReady canceled: %w", err)
		default:
			log.Printf("WaitServerReady port %d: %v", port, err)
		}
	}
}
