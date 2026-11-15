package main

import (
	"bytes"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

// The test binary runs the server itself when the start test asks it to, so
// the test starts the real main, as the image does.
func TestMain(m *testing.M) {
	if os.Getenv("RUN_MAIN_FOR_TEST") == "1" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	t.Run("T.C3 it starts, answers its health check and stops when asked", func(t *testing.T) {
		addr := freeAddress(t)
		var out bytes.Buffer
		cmd := exec.Command(os.Args[0], "-test.run=^$")
		cmd.Env = append(os.Environ(), "RUN_MAIN_FOR_TEST=1", "ADDR="+addr)
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		t.Cleanup(func() { _ = cmd.Process.Kill() })

		url := "http://" + addr + "/api/healthz"
		deadline := time.Now().Add(20 * time.Second)
		for {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					break
				}
				err = &statusError{resp.StatusCode}
			}
			select {
			case werr := <-done:
				t.Fatalf("it stopped before answering %s: %v\n%s", url, werr, out.String())
			default:
			}
			if time.Now().After(deadline) {
				t.Fatalf("no answer from %s within 20 s: %v\n%s", url, err, out.String())
			}
			time.Sleep(100 * time.Millisecond)
		}

		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			t.Fatal(err)
		}
		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("it did not stop cleanly: %v\n%s", err, out.String())
			}
		case <-time.After(20 * time.Second):
			t.Fatalf("it did not stop within 20 s of being asked\n%s", out.String())
		}
	})
}

type statusError struct{ code int }

func (e *statusError) Error() string { return http.StatusText(e.code) }

// freeAddress is an address on the loopback interface that nothing listens on.
func freeAddress(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}
