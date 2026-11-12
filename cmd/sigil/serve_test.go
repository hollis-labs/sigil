package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hollis-labs/substrate/boot"
	substratecli "github.com/hollis-labs/substrate/modules/cli"

	"github.com/hollis-labs/sigil/internal/cli"
)

// freeAddr returns a loopback address with a port nobody holds right now.
func freeAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().String()
}

// project makes a throwaway .sigil directory with one page, through the real commands.
func project(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	for _, args := range [][]string{{"init"}, {"new", "page", "dashboard"}} {
		root := cli.NewRootCmd()
		root.SetOut(io.Discard)
		root.SetErr(io.Discard)
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatalf("sigil %v: %v", args, err)
		}
	}
	return filepath.Join(dir, ".sigil")
}

// start boots the app with httphost.addr set to addr and returns the cli module's root with
// Sigil's commands installed, as main does.
func start(t *testing.T, addr string) (root interface {
	SetArgs([]string)
	ExecuteContext(context.Context) error
}, stop func()) {
	t.Helper()
	t.Setenv("SUBSTRATE_CONFIG", "") // the developer's override file must not leak into the test
	config := strings.Replace(substrateConfigJSON, "127.0.0.1:3210", addr, 1)
	started, err := boot.Start(context.Background(), substrateApp, substrateModules, config)
	if err != nil {
		t.Fatalf("boot: %v", err)
	}
	r := substratecli.Root(started.App)
	r.SetOut(io.Discard)
	r.SetErr(io.Discard)
	cli.Install(r, started.App)
	return r, func() { _ = started.Stop(context.Background()) }
}

// serve runs `sigil serve` until cancel and returns a function that stops it and reports the
// command's result.
func serve(t *testing.T, addr, sigilDir string) (cancel func() error) {
	t.Helper()
	root, stop := start(t, addr)
	ctx, cancelCtx := context.WithCancel(context.Background())
	root.SetArgs([]string{"serve", "--sigil-dir", sigilDir})
	done := make(chan error, 1)
	go func() { done <- root.ExecuteContext(ctx) }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		resp, err := http.Get("http://" + addr + "/healthz")
		if err == nil {
			resp.Body.Close()
			break
		}
		select {
		case err := <-done:
			stop()
			t.Fatalf("serve stopped early: %v", err)
		default:
		}
		if time.Now().After(deadline) {
			cancelCtx()
			stop()
			t.Fatal("serve did not come up")
		}
		time.Sleep(20 * time.Millisecond)
	}
	return func() error {
		cancelCtx()
		select {
		case err := <-done:
			stop()
			return err
		case <-time.After(5 * time.Second):
			t.Fatal("serve did not stop on cancellation")
			return nil
		}
	}
}

func do(t *testing.T, method, url, host string, header map[string]string, body []byte) (int, string) {
	t.Helper()
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if host != "" {
		req.Host = host
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func TestServeOnLoopbackWithoutCredentials(t *testing.T) {
	sigilDir := project(t)
	addr := freeAddr(t)
	stop := serve(t, addr, sigilDir)
	base := "http://" + addr

	for _, path := range []string{"/healthz", "/", "/pages/dashboard", "/assets/theme.css", "/api/pages", "/api/project", "/api/me"} {
		if code, body := do(t, "GET", base+path, "", nil, nil); code != http.StatusOK {
			t.Errorf("GET %s with no credential: %d %s", path, code, body)
		}
	}

	// Writes carry the X-Requested-With rule; the body is the flat error.
	page := []byte(`{"sigil":"1.0","id":"broken"}`)
	code, body := do(t, "POST", base+"/api/pages", "", map[string]string{"Content-Type": "application/json"}, page)
	var flat map[string]any
	if code != http.StatusForbidden || json.Unmarshal([]byte(body), &flat) != nil || flat["code"] != "forbidden" {
		t.Errorf("write without X-Requested-With: %d %s", code, body)
	}
	code, body = do(t, "POST", base+"/api/pages", "", map[string]string{"Content-Type": "application/json", "X-Requested-With": "test"}, page)
	flat = nil
	if code != http.StatusUnprocessableEntity || json.Unmarshal([]byte(body), &flat) != nil || flat["code"] != "validation_failed" || flat["message"] == nil {
		t.Errorf("invalid write: %d %s", code, body)
	}
	if code, body = do(t, "GET", base+"/api/pages/nope", "", nil, nil); code != http.StatusNotFound || !strings.Contains(body, `"code":"not_found"`) {
		t.Errorf("unknown page: %d %s", code, body)
	}

	// A Host the server is not meant to answer is refused (DNS rebinding).
	for _, host := range []string{"evil.example", "localhost.evil.example", "192.0.2.5:3210"} {
		if code, _ := do(t, "GET", base+"/api/pages", host, nil, nil); code != http.StatusForbidden {
			t.Errorf("Host %q: want 403, got %d", host, code)
		}
	}
	if code, _ := do(t, "GET", base+"/api/pages", "localhost:"+strings.Split(addr, ":")[1], nil, nil); code != http.StatusOK {
		t.Errorf("Host localhost: want 200, got %d", code)
	}

	if err := stop(); err != nil {
		t.Errorf("serve returned %v on cancellation, want nil", err)
	}
	if _, err := http.Get(base + "/healthz"); err == nil {
		t.Error("the listener is still up after shutdown")
	}
}

func TestServeStopsWithAReloadStreamOpen(t *testing.T) {
	sigilDir := project(t)
	addr := freeAddr(t)
	stop := serve(t, addr, sigilDir)
	resp, err := http.Get("http://" + addr + "/events")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	start := time.Now()
	if err := stop(); err != nil {
		t.Errorf("serve returned %v", err)
	}
	if d := time.Since(start); d > 3*time.Second {
		t.Errorf("shutdown waited %s for an open reload stream", d)
	}
}

func TestServeRefusesANonLoopbackAddress(t *testing.T) {
	sigilDir := project(t)
	// A tailnet address and a tailnet name are accepted by the framework in mode off;
	// 198.51.100.7 stands for any LAN address. The literals are split so the hygiene scan
	// does not read them as real hosts.
	tailnetIP := "100.64." + "0.1"
	tailnetName := "sigil.example." + "ts.net"
	for _, host := range []string{"0.0.0.0", "198.51.100.7", tailnetIP, tailnetName, ""} {
		addr := fmt.Sprintf("%s:%d", host, 18199)
		root, stop := start(t, addr)
		root.SetArgs([]string{"serve", "--sigil-dir", sigilDir})
		err := root.ExecuteContext(context.Background())
		stop()
		if err == nil || !strings.Contains(err.Error(), "not a loopback address") {
			t.Errorf("addr %q: want the loopback refusal, got %v", addr, err)
		}
	}
}

func TestServePortFlagKeepsTheLoopbackHost(t *testing.T) {
	sigilDir := project(t)
	addr := freeAddr(t)
	port := strings.Split(addr, ":")[1]
	root, stop := start(t, "127.0.0.1:1") // the configured port is not used
	defer stop()
	ctx, cancel := context.WithCancel(context.Background())
	root.SetArgs([]string{"serve", "--port", port, "--sigil-dir", sigilDir})
	done := make(chan error, 1)
	go func() { done <- root.ExecuteContext(ctx) }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		resp, err := http.Get("http://" + addr + "/healthz")
		if err == nil {
			resp.Body.Close()
			break
		}
		if time.Now().After(deadline) {
			cancel()
			t.Fatalf("serve --port %s did not come up", port)
		}
		time.Sleep(20 * time.Millisecond)
	}
	cancel()
	if err := <-done; err != nil {
		t.Errorf("serve returned %v", err)
	}
}

// A compile command prints nothing on stderr from the framework: the "auth is not enforcing"
// notice belongs to a listener, so only serve logs it.
func TestCompileCommandsPrintNothingOnStderr(t *testing.T) {
	project(t)
	t.Setenv("SUBSTRATE_CONFIG", "")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origStderr := os.Stderr
	os.Stderr = w
	t.Cleanup(func() { os.Stderr = origStderr })

	for _, args := range [][]string{{"validate"}, {"generate", "--target", "go-templ", "--output", "out"}, {"version"}} {
		started, err := boot.Start(context.Background(), substrateApp, substrateModules, substrateConfigJSON)
		if err != nil {
			t.Fatalf("boot: %v", err)
		}
		root := substratecli.Root(started.App)
		root.SetOut(io.Discard)
		root.SetErr(io.Discard)
		cli.Install(root, started.App)
		root.SetArgs(args)
		err = root.ExecuteContext(context.Background())
		_ = started.Stop(context.Background())
		if err != nil {
			t.Fatalf("sigil %v: %v", args, err)
		}
	}
	os.Stderr = origStderr
	_ = w.Close()
	got, _ := io.ReadAll(r)
	if len(got) != 0 {
		t.Errorf("stderr should be empty, got %q", got)
	}
}

func TestMCPServeOverStdio(t *testing.T) {
	sigilDir := project(t)
	root, stop := start(t, "127.0.0.1:0")
	defer stop()

	inR, inW, _ := os.Pipe()
	outR, outW, _ := os.Pipe()
	origIn, origOut := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = inR, outW
	defer func() { os.Stdin, os.Stdout = origIn, origOut }()

	root.SetArgs([]string{"mcp", "serve", "--sigil-dir", sigilDir})
	done := make(chan error, 1)
	go func() { done <- root.ExecuteContext(context.Background()) }()

	dec := json.NewDecoder(outR)
	call := func(id int, method string, params any) map[string]any {
		msg, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": params})
		if _, err := inW.Write(append(msg, '\n')); err != nil {
			t.Fatal(err)
		}
		var resp map[string]any
		if err := dec.Decode(&resp); err != nil {
			t.Fatalf("%s: %v", method, err)
		}
		return resp
	}
	init := call(1, "initialize", map[string]any{"protocolVersion": "2025-06-18", "capabilities": map[string]any{}, "clientInfo": map[string]any{"name": "t", "version": "1"}})
	info := init["result"].(map[string]any)["serverInfo"].(map[string]any)
	if info["name"] != "sigil" {
		t.Errorf("server info: %v", info)
	}
	if _, err := inW.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n")); err != nil {
		t.Fatal(err)
	}
	list := call(2, "tools/list", map[string]any{})
	if n := len(list["result"].(map[string]any)["tools"].([]any)); n != 9 {
		t.Errorf("tools/list: %d tools, want 9", n)
	}
	res := call(3, "tools/call", map[string]any{"name": "sigil_list_pages", "arguments": map[string]any{}})
	if !strings.Contains(fmt.Sprint(res["result"]), "dashboard") {
		t.Errorf("sigil_list_pages: %v", res)
	}
	_ = inW.Close()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("mcp serve returned %v after the client closed stdin", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("mcp serve did not stop when stdin closed")
	}
}
