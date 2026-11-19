package wails_test

import (
	"testing"

	"github.com/macrox-pro/route-deck/internal/adapter/inbound/wails"
)

func TestMenuItemEnabled(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		command    string
		editorOpen bool
		cliOwned   bool
		want       bool
	}{
		{name: "about always", command: "about", editorOpen: true, want: true},
		{name: "settings always", command: "settings", editorOpen: true, want: true},
		{name: "new-profile always", command: "new-profile", editorOpen: true, want: true},
		{name: "new-proxy always", command: "new-proxy", editorOpen: true, want: true},
		{name: "setup-guides always", command: "setup-guides", editorOpen: true, want: true},
		{name: "apply disabled while editor", command: "apply", editorOpen: true, want: false},
		{name: "apply enabled no editor", command: "apply", editorOpen: false, want: true},
		{name: "recheck disabled while editor", command: "recheck", editorOpen: true, want: false},
		{name: "add-application disabled while editor", command: "add-application", editorOpen: true, want: false},
		{name: "navigate disabled while editor", command: "navigate", editorOpen: true, want: false},
		{name: "toggle-sidebar disabled while editor", command: "toggle-sidebar", editorOpen: true, want: false},
		{name: "stop needs owned and no editor", command: "stop", editorOpen: false, cliOwned: true, want: true},
		{name: "stop blocked when not owned", command: "stop", editorOpen: false, cliOwned: false, want: false},
		{name: "stop blocked while editor even if owned", command: "stop", editorOpen: true, cliOwned: true, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := wails.MenuItemEnabled(tc.command, tc.editorOpen, tc.cliOwned)
			if got != tc.want {
				t.Fatalf("MenuItemEnabled(%q, editorOpen=%v, cliOwned=%v) = %v, want %v",
					tc.command, tc.editorOpen, tc.cliOwned, got, tc.want)
			}
		})
	}
}

func TestGitHubRepoURL(t *testing.T) {
	t.Parallel()
	if wails.GitHubRepoURL != "https://github.com/macrox-pro/route-deck" {
		t.Fatalf("GitHubRepoURL = %q", wails.GitHubRepoURL)
	}
}

func TestEditorDialogs_AnyOpen(t *testing.T) {
	t.Parallel()
	var d *wails.EditorDialogs
	if d.AnyOpen() {
		t.Fatal("nil EditorDialogs AnyOpen should be false")
	}
	d = &wails.EditorDialogs{}
	if d.AnyOpen() {
		t.Fatal("empty EditorDialogs AnyOpen should be false")
	}
}
