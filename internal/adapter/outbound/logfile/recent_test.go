package logfile_test

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/logfile"
)

func TestRecentMissingFile(t *testing.T) {
	_, err := logfile.Recent(context.Background(), filepath.Join(t.TempDir(), "missing.log"), 10, "")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRecentInvalidBefore(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")
	if err := os.WriteFile(path, []byte(
		`{"time":"2026-11-17T10:00:00Z","level":"INFO","msg":"ok","component":"main"}`+"\n",
	), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := logfile.Recent(context.Background(), path, 10, "not-a-time")
	if err == nil {
		t.Fatal("expected invalid before error")
	}
}

func TestRecentSkipsBadLinesAndDropsSecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")
	body := "" +
		"not json\n" +
		`{"time":"2026-11-17T10:00:00Z","level":"INFO","msg":"ok","component":"main","password":"secret","token":"t","extra":"x"}` + "\n" +
		"{broken\n" +
		`{"time":"2026-11-17T09:00:00Z","level":"INFO","msg":"older","component":"health"}` + "\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := logfile.Recent(context.Background(), path, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 2 {
		t.Fatalf("len=%d want 2", len(got.Records))
	}
	if got.HasMore {
		t.Fatal("HasMore want false")
	}
	if got.Records[0].Msg != "ok" || got.Records[0].Time != "2026-11-17T10:00:00Z" {
		t.Fatalf("newest = %+v", got.Records[0])
	}
	if got.Records[0].Component != "main" {
		t.Fatalf("component = %q", got.Records[0].Component)
	}
	attrs := map[string]string{}
	for _, attr := range got.Records[0].Attrs {
		attrs[attr.Key] = attr.Value
	}
	if attrs["extra"] != "x" {
		t.Fatalf("attrs = %+v, want extra=x", got.Records[0].Attrs)
	}
	if _, ok := attrs["password"]; ok {
		t.Fatal("password leaked into attrs")
	}
	if _, ok := attrs["token"]; ok {
		t.Fatal("token leaked into attrs")
	}
	if _, ok := attrs["msg"]; ok {
		t.Fatal("msg must not be duplicated in attrs")
	}
}

func TestRecentKeepsPlatformInfoAttrs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")
	if err := os.WriteFile(path, []byte(
		`{"time":"2026-11-17T16:15:48Z","level":"INFO","msg":"Platform Info:","ID":"26B131","Name":"MacOS","Version":"27.2.1","Branding":"MacOS 27.2.1"}`+"\n",
	), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := logfile.Recent(context.Background(), path, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 1 || got.Records[0].Msg != "Platform Info:" {
		t.Fatalf("got %+v", got.Records)
	}
	attrs := map[string]string{}
	for _, attr := range got.Records[0].Attrs {
		attrs[attr.Key] = attr.Value
	}
	if attrs["Name"] != "MacOS" || attrs["Version"] != "27.2.1" || attrs["Branding"] != "MacOS 27.2.1" {
		t.Fatalf("attrs = %+v", got.Records[0].Attrs)
	}
}

func TestRecentCapHasMoreAndBefore(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")
	if err := os.WriteFile(path+".1", []byte(
		`{"time":"2026-11-17T08:00:00Z","level":"INFO","msg":"rot1","component":"main"}`+"\n",
	), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(
		`{"time":"2026-11-17T09:00:00Z","level":"INFO","msg":"a","component":"main"}`+"\n"+
			`{"time":"2026-11-17T10:00:00Z","level":"INFO","msg":"b","component":"main"}`+"\n"+
			`{"time":"2026-11-17T11:00:00Z","level":"INFO","msg":"c","component":"main"}`+"\n",
	), 0o600); err != nil {
		t.Fatal(err)
	}
	first, err := logfile.Recent(context.Background(), path, 2, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Records) != 2 || !first.HasMore {
		t.Fatalf("first = %+v", first)
	}
	if first.Records[0].Msg != "c" || first.Records[1].Msg != "b" {
		t.Fatalf("order = %q, %q", first.Records[0].Msg, first.Records[1].Msg)
	}
	second, err := logfile.Recent(context.Background(), path, 2, first.Records[len(first.Records)-1].Time)
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Records) != 2 {
		t.Fatalf("second len=%d", len(second.Records))
	}
	if second.Records[0].Msg != "a" || second.Records[1].Msg != "rot1" {
		t.Fatalf("second order = %q, %q", second.Records[0].Msg, second.Records[1].Msg)
	}
	if second.HasMore {
		t.Fatal("second HasMore want false")
	}
}

func TestRecentDuplicateMsgKeepsLast(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")
	// encoding/json keeps the last duplicate key.
	line := `{"time":"2026-11-17T10:00:00Z","level":"INFO","msg":"slog","msg":"cli-text","component":"proxybridge","line_class":"redirect_ready"}` + "\n"
	if err := os.WriteFile(path, []byte(line), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := logfile.Recent(context.Background(), path, 5, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 1 {
		t.Fatalf("len=%d", len(got.Records))
	}
	if got.Records[0].Msg != "cli-text" {
		t.Fatalf("msg = %q", got.Records[0].Msg)
	}
	if got.Records[0].LineClass != "redirect_ready" {
		t.Fatalf("line_class = %q", got.Records[0].LineClass)
	}
}

func TestRecentSlogJSONFixture(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{ReplaceAttr: logfile.RedactAttr})
	logger := slog.New(handler)
	logger.Info("bridge line",
		"component", "proxybridge",
		"line_class", "redirect_ready",
		"msg", "Redirecting traffic",
		"password", "must-not-appear",
	)
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte("must-not-appear")) {
		t.Fatalf("password leaked into file: %s", raw)
	}
	got, err := logfile.Recent(context.Background(), path, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 1 {
		t.Fatalf("len=%d body=%s", len(got.Records), raw)
	}
	if got.Records[0].Component != "proxybridge" {
		t.Fatalf("component = %q", got.Records[0].Component)
	}
	if got.Records[0].LineClass != "redirect_ready" {
		t.Fatalf("line_class = %q", got.Records[0].LineClass)
	}
	if got.Records[0].Msg == "" {
		t.Fatal("empty msg")
	}
}

func TestReaderImplementsPort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "route-deck.log")
	if err := os.WriteFile(path, []byte(
		`{"time":"2026-11-17T10:00:00Z","level":"INFO","msg":"hi","component":"main"}`+"\n",
	), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := (&logfile.Reader{}).Recent(context.Background(), path, 1, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 1 || got.Records[0].Msg != "hi" {
		t.Fatalf("got %+v", got)
	}
}
