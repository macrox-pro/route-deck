package proxyhealth_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/macrox-pro/route-deck/internal/adapter/outbound/proxyhealth"
)

func TestLocalPortProberOpen(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	port := ln.Addr().(*net.TCPAddr).Port

	p := &proxyhealth.LocalPortProber{Timeout: 300 * time.Millisecond}
	open, err := p.Open(context.Background(), port)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if !open {
		t.Fatal("want open true")
	}

	closed, err := p.Open(context.Background(), 1)
	if err != nil {
		t.Fatalf("Open closed: %v", err)
	}
	if closed {
		t.Fatal("want open false for unused port")
	}
}

func TestLocalPortProberInvalidPort(t *testing.T) {
	p := &proxyhealth.LocalPortProber{}
	_, err := p.Open(context.Background(), 0)
	if err == nil {
		t.Fatal("want error for port 0")
	}
}
