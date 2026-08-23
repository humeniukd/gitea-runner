// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubAction struct {
	name    string
	hasPre  bool
	hasPost bool
	calls   map[Stage]int
	err     error
}

func newStub(name string) *stubAction {
	return &stubAction{name: name, calls: map[Stage]int{}}
}

func (s *stubAction) Name() string  { return s.name }
func (s *stubAction) HasPre() bool  { return s.hasPre }
func (s *stubAction) HasPost() bool { return s.hasPost }
func (s *stubAction) Pre(_ *Context) error {
	s.calls[StagePre]++
	return s.err
}

func (s *stubAction) Main(_ *Context) error {
	s.calls[StageMain]++
	return s.err
}

func (s *stubAction) Post(_ *Context) error {
	s.calls[StagePost]++
	return s.err
}

func TestRegisterAndLookup(t *testing.T) {
	defer resetForTest()
	resetForTest()

	a := newStub("checkout")
	Register(a)

	got, ok := Lookup("checkout")
	assert.True(t, ok)
	assert.Same(t, a, got)

	_, ok = Lookup("nope")
	assert.False(t, ok)
}

func TestRegisterDuplicatePanics(t *testing.T) {
	defer resetForTest()
	resetForTest()

	Register(newStub("dup"))

	assert.Panics(t, func() {
		Register(newStub("dup"))
	})
}

func TestRegisterInvalidNamePanics(t *testing.T) {
	defer resetForTest()
	resetForTest()

	for _, bad := range []string{"", "-leading-dash", "Has-Upper", "has/slash", "has@v1", "has space", "has.dot"} {
		t.Run(bad, func(t *testing.T) {
			assert.Panics(t, func() {
				Register(newStub(bad))
			})
		})
	}
}

func TestRegisterValidNames(t *testing.T) {
	defer resetForTest()

	for _, good := range []string{"a", "checkout", "upload-artifact", "x1-y2"} {
		t.Run(good, func(t *testing.T) {
			resetForTest()
			Register(newStub(good))
			_, ok := Lookup(good)
			assert.True(t, ok)
		})
	}
}

func TestParseUses(t *testing.T) {
	for _, tc := range []struct {
		uses string
		name string
		ok   bool
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

func TestRegisterNilPanics(t *testing.T) {
	defer resetForTest()
	resetForTest()

	assert.Panics(t, func() { Register(nil) })
}

func TestNamesSorted(t *testing.T) {
	defer resetForTest()
	resetForTest()

	Register(newStub("zeta"))
	Register(newStub("alpha"))
	Register(newStub("mu"))

	assert.Equal(t, []string{"alpha", "mu", "zeta"}, Names())
}

func TestDispatchStage(t *testing.T) {
	a := newStub("foo")
	for _, s := range []Stage{StagePre, StageMain, StagePost} {
		err := DispatchStage(a, &Context{}, s)
		require.NoError(t, err)
	}
	assert.Equal(t, 1, a.calls[StagePre])
	assert.Equal(t, 1, a.calls[StageMain])
	assert.Equal(t, 1, a.calls[StagePost])

	a.err = errors.New("boom")
	err := DispatchStage(a, &Context{}, StageMain)
	require.Error(t, err)

	err = DispatchStage(a, &Context{}, Stage(99))
	require.Error(t, err)
}

func TestStageString(t *testing.T) {
	assert.Equal(t, "pre", StagePre.String())
	assert.Equal(t, "main", StageMain.String())
	assert.Equal(t, "post", StagePost.String())
	assert.Equal(t, "unknown", Stage(99).String())
}
