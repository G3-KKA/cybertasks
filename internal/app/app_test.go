package app

import (
	"cybertask/config"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunShutdown(t *testing.T) {
	// t.Parallel() Not Parallelable test because of syscall.
	const (
		shutdownTimeout = time.Microsecond * 1500 // 1.5 sec.
		failTimeout     = time.Second * 5
	)

	cfg, err := config.Get()
	assert.NoError(t, err)

	app, err := New(cfg)
	assert.NoError(t, err)

	errchan := make(chan error, 1)

	shutdown := func() {
		time.Sleep(shutdownTimeout)
		err2 := syscall.Kill(os.Getpid(), syscall.SIGINT)
		assert.NoError(t, err2)
	}
	runner := func() {
		errchan <- app.Run()
		close(errchan)
	}
	go runner()
	go shutdown()

	timer := time.NewTimer(failTimeout)
	select {
	case <-timer.C:
		t.Fatalf("app shutdown has taken more than %s", failTimeout.String())
	case err = <-errchan:
		assert.NoError(t, err)
	}
}
