// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubAction struct{ name string }

func (s *stubAction) Name() string        { return s.name }
func (s *stubAction) HasPre() bool        { return false }
func (s *stubAction) HasPost() bool       { return false }
func (s *stubAction) Pre(*Context) error  { return nil }
func (s *stubAction) Main(*Context) error { return nil }
func (s *stubAction) Post(*Context) error { return nil }

func TestRegister(t *testing.T) {
	t.Cleanup(resetForTest)
	resetForTest()

	a := &stubAction{name: "checkout"}
	Register(a)
	got, ok := Lookup("checkout")
	assert.True(t, ok)
	assert.Same(t, a, got)
	_, ok = Lookup("nope")
	assert.False(t, ok)

	assert.Panics(t, func() { Register(&stubAction{name: "checkout"}) })
	assert.Panics(t, func() { Register(nil) })
	assert.Panics(t, func() { Register(&stubAction{name: "-bad"}) })
	assert.Panics(t, func() { Register(&stubAction{name: "Has-Upper"}) })
	assert.Panics(t, func() { Register(&stubAction{name: "has/slash"}) })

	resetForTest()
	for _, name := range []string{"a", "checkout", "upload-artifact", "x1-y2"} {
		Register(&stubAction{name: name})
	}
	assert.Equal(t, []string{"a", "checkout", "upload-artifact", "x1-y2"}, Names())
}

func TestParseUses(t *testing.T) {
	for _, tc := range []struct {
		uses, name string
		ok         bool
	}{
		{uses: "builtin:checkout", name: "checkout", ok: true},
		{uses: "checkout"},
		{uses: "builtin:", ok: true},
		{uses: "builtin:has/slash", name: "has/slash", ok: true},
		{uses: "builtin:Has-Upper", name: "Has-Upper", ok: true},
	} {
		t.Run(tc.uses, func(t *testing.T) {
			name, ok := ParseUses(tc.uses)
			assert.Equal(t, tc.name, name)
			assert.Equal(t, tc.ok, ok)
		})
	}
}
