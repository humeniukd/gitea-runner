// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.yaml.in/yaml/v4"
)

// fieldKind describes what a resolved config path points at, which decides
// whether a path can be set, appended to or listed.
type fieldKind int

const (
	kindScalar fieldKind = iota
	kindSequence
	kindMapping
	kindSection
)

var durationType = reflect.TypeFor[time.Duration]()

// GetValue returns the value stored at path in file, rendered one entry per
// line for sequences and mappings.
func GetValue(file, path string) (string, error) {
	session, err := loadForEdit(file, path)
	if err != nil {
		return "", err
	}

	node, err := lookupNode(session.root, strings.Split(path, "."))
	if err != nil {
		return "", err
	}

	return renderNode(node), nil
}

// SetValue replaces the value at path. A sequence path takes every value, any
// other path takes exactly one.
func SetValue(file, path string, values ...string) error {
	session, err := loadForEdit(file, path)
	if err != nil {
		return err
	}

	switch session.field.kind {
	case kindSequence:
		if len(values) == 0 {
			return fmt.Errorf("%q needs at least one value", path)
		}
		node, err := ensureNode(session.root, strings.Split(path, "."))
		if err != nil {
			return err
		}
		items := make([]*yaml.Node, 0, len(values))
		for _, value := range values {
			scalar, err := scalarNode(session.field.elem, value)
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			items = append(items, scalar)
		}
		node.Kind, node.Tag, node.Style, node.Value, node.Content = yaml.SequenceNode, "!!seq", 0, "", items
	case kindScalar:
		if len(values) != 1 {
			return fmt.Errorf("%q takes exactly one value", path)
		}
		scalar, err := scalarNode(session.field.typ, values[0])
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		node, err := ensureNode(session.root, strings.Split(path, "."))
		if err != nil {
			return err
		}
		node.Kind, node.Tag, node.Style, node.Value, node.Content = scalar.Kind, scalar.Tag, scalar.Style, scalar.Value, nil
	default:
		return fmt.Errorf("%q is a section, set one of its keys instead", path)
	}

	return session.write(file)
}

// AddValue appends values to the sequence at path.
func AddValue(file, path string, values ...string) error {
	session, err := loadForEdit(file, path)
	if err != nil {
		return err
	}
	if session.field.kind != kindSequence {
		return fmt.Errorf("%q is not a list, use `config set` instead", path)
	}
	if len(values) == 0 {
		return fmt.Errorf("%q needs at least one value", path)
	}

	node, err := ensureNode(session.root, strings.Split(path, "."))
	if err != nil {
		return err
	}
	if node.Kind != yaml.SequenceNode {
		node.Kind, node.Tag, node.Style, node.Value, node.Content = yaml.SequenceNode, "!!seq", 0, "", nil
	}

	for _, value := range values {
		scalar, err := scalarNode(session.field.elem, value)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if indexOfValue(node, scalar.Value) >= 0 {
			return fmt.Errorf("%s already contains %q", path, scalar.Value)
		}
		node.Content = append(node.Content, scalar)
	}

	return session.write(file)
}

// RemoveValue drops values from the sequence at path.
func RemoveValue(file, path string, values ...string) error {
	session, err := loadForEdit(file, path)
	if err != nil {
		return err
	}
	if session.field.kind != kindSequence {
		return fmt.Errorf("%q is not a list, use `config set` instead", path)
	}
	if len(values) == 0 {
		return fmt.Errorf("%q needs at least one value", path)
	}

	node, err := lookupNode(session.root, strings.Split(path, "."))
	if err != nil {
		return err
	}
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("%s is not a list in %q", path, file)
	}

	for _, value := range values {
		scalar, err := scalarNode(session.field.elem, value)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		index := indexOfValue(node, scalar.Value)
		if index < 0 {
			return fmt.Errorf("%s does not contain %q", path, scalar.Value)
		}
		node.Content = append(node.Content[:index], node.Content[index+1:]...)
	}

	return session.write(file)
}

func indexOfValue(seq *yaml.Node, value string) int {
	for i, item := range seq.Content {
		if item.Kind == yaml.ScalarNode && item.Value == value {
			return i
		}
	}
	return -1
}

type editSession struct {
	original []byte
	root     *yaml.Node
	field    *fieldInfo
}

// loadForEdit validates path against the Config schema and parses file into a
// document node, so every caller fails before touching the file on disk.
func loadForEdit(file, path string) (*editSession, error) {
	if path == "" {
		return nil, errors.New("no config key given")
	}
	field, err := resolvePath(path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file %q does not exist, create one with `gitea-runner config generate > %s`", file, file)
		}
		return nil, err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("parse config file %q: %w", file, err)
	}
	if root.Kind == 0 || len(root.Content) == 0 {
		root = yaml.Node{
			Kind:    yaml.DocumentNode,
			Content: []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}},
		}
	}
	if root.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("config file %q is not a YAML mapping", file)
	}

	return &editSession{original: content, root: &root, field: field}, nil
}

// lookupNode resolves segments in an existing document without creating keys.
func lookupNode(root *yaml.Node, segments []string) (*yaml.Node, error) {
	node := root.Content[0]
	for i, segment := range segments {
		if node.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("%q is not set", strings.Join(segments[:i], "."))
		}
		value := mappingValue(node, segment)
		if value == nil {
			return nil, fmt.Errorf("%q is not set", strings.Join(segments[:i+1], "."))
		}
		node = value
	}
	return node, nil
}

// ensureNode resolves segments, creating the intermediate mappings and the leaf
// placeholder when they are missing.
func ensureNode(root *yaml.Node, segments []string) (*yaml.Node, error) {
	node := root.Content[0]
	for i, segment := range segments {
		if node.Kind != yaml.MappingNode {
			if node.Kind == yaml.ScalarNode && node.Tag == "!!null" {
				node.Kind, node.Tag, node.Style, node.Value = yaml.MappingNode, "!!map", 0, ""
			} else {
				return nil, fmt.Errorf("%q is not a section", strings.Join(segments[:i], "."))
			}
		}
		value := mappingValue(node, segment)
		if value == nil {
			value = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: segment},
				value)
		}
		node = value
	}
	return node, nil
}

func mappingValue(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

func renderNode(node *yaml.Node) string {
	switch node.Kind {
	case yaml.SequenceNode:
		lines := make([]string, 0, len(node.Content))
		for _, item := range node.Content {
			lines = append(lines, item.Value)
		}
		return strings.Join(lines, "\n")
	case yaml.MappingNode:
		lines := make([]string, 0, len(node.Content)/2)
		for i := 0; i+1 < len(node.Content); i += 2 {
			lines = append(lines, node.Content[i].Value+"="+node.Content[i+1].Value)
		}
		return strings.Join(lines, "\n")
	default:
		return node.Value
	}
}

// restoreBlankLines re-inserts the blank lines between top-level sections that
// the encoder drops, so one edit does not reflow the whole file.
func restoreBlankLines(original, generated []byte) []byte {
	spaced := map[string]bool{}
	blank := false
	for line := range strings.Lines(string(original)) {
		line = strings.TrimRight(line, "\r\n")
		switch {
		case strings.TrimSpace(line) == "":
			blank = true
		case strings.HasPrefix(line, "#"): // the block belongs to the key below it
		default:
			if key, ok := topLevelKey(line); ok && blank {
				spaced[key] = true
			}
			blank = false
		}
	}
	if len(spaced) == 0 {
		return generated
	}

	var out []string
	for line := range strings.Lines(string(generated)) {
		line = strings.TrimRight(line, "\r\n")
		if key, ok := topLevelKey(line); ok && spaced[key] {
			start := len(out)
			for start > 0 && strings.HasPrefix(out[start-1], "#") {
				start--
			}
			if start > 0 && strings.TrimSpace(out[start-1]) != "" {
				out = slices.Insert(out, start, "")
			}
		}
		out = append(out, line)
	}
	if bytes.HasSuffix(generated, []byte("\n")) {
		out = append(out, "")
	}

	return []byte(strings.Join(out, "\n"))
}

func topLevelKey(line string) (string, bool) {
	if line == "" || line[0] == ' ' || line[0] == '\t' || line[0] == '#' || line[0] == '-' {
		return "", false
	}
	key, _, ok := strings.Cut(line, ":")
	return key, ok
}

// write re-encodes the edited document, rejects it if it no longer parses as a
// config, and swaps it in atomically so a failure leaves the original in place.
func (s *editSession) write(file string) error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(s.root); err != nil {
		return err
	}
	if err := encoder.Close(); err != nil {
		return err
	}

	// Unknown keys stay tolerated, as LoadDefault only warns about them.
	var typeErr *yaml.TypeError
	if err := yaml.Unmarshal(buf.Bytes(), &Config{}); err != nil && !errors.As(err, &typeErr) {
		return fmt.Errorf("the edit would produce an invalid config: %w", err)
	}
	content := restoreBlankLines(s.original, buf.Bytes())

	mode := os.FileMode(0o600)
	if info, err := os.Stat(file); err == nil {
		mode = info.Mode().Perm()
	}

	temp, err := os.CreateTemp(filepath.Dir(file), filepath.Base(file)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())

	if _, err := temp.Write(content); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(temp.Name(), mode); err != nil {
		return err
	}

	return os.Rename(temp.Name(), file)
}

type fieldInfo struct {
	kind fieldKind
	typ  reflect.Type // the leaf type, for scalars
	elem reflect.Type // the element type, for sequences
}

// resolvePath walks the Config struct through the yaml tags of a dotted path,
// so an unknown or mistyped key is rejected before anything is written.
func resolvePath(path string) (*fieldInfo, error) {
	typ := reflect.TypeFor[Config]()
	segments := strings.Split(path, ".")

	for i, segment := range segments {
		switch typ.Kind() {
		case reflect.Struct:
			field, ok := fieldByYAMLName(typ, segment)
			if !ok {
				return nil, fmt.Errorf("unknown config key %q, valid keys here: %s",
					strings.Join(segments[:i+1], "."), strings.Join(yamlNames(typ), ", "))
			}
			typ = field.Type
		case reflect.Map:
			// The segment names a user-defined entry, so the walk ends here.
			if i != len(segments)-1 {
				return nil, fmt.Errorf("%q has no sub-keys", strings.Join(segments[:i+1], "."))
			}
			return &fieldInfo{kind: kindScalar, typ: typ.Elem()}, nil
		default:
			return nil, fmt.Errorf("%q is a value, not a section", strings.Join(segments[:i], "."))
		}
	}

	return describe(typ), nil
}

func describe(typ reflect.Type) *fieldInfo {
	switch {
	case typ.Kind() == reflect.Slice:
		return &fieldInfo{kind: kindSequence, elem: typ.Elem()}
	case typ.Kind() == reflect.Map:
		return &fieldInfo{kind: kindMapping}
	case typ.Kind() == reflect.Struct:
		return &fieldInfo{kind: kindSection}
	default:
		return &fieldInfo{kind: kindScalar, typ: typ}
	}
}

func fieldByYAMLName(typ reflect.Type, name string) (reflect.StructField, bool) {
	for field := range typ.Fields() {
		if yamlName(field) == name {
			return field, true
		}
	}
	return reflect.StructField{}, false
}

func yamlNames(typ reflect.Type) []string {
	names := make([]string, 0, typ.NumField())
	for field := range typ.Fields() {
		if name := yamlName(field); name != "" && name != "-" {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names
}

func yamlName(field reflect.StructField) string {
	name, _, _ := strings.Cut(field.Tag.Get("yaml"), ",")
	if name == "" {
		return strings.ToLower(field.Name)
	}
	return name
}

// scalarNode converts a command line value into a typed YAML scalar, so a bad
// value is reported instead of silently landing in the file as a string.
func scalarNode(typ reflect.Type, value string) (*yaml.Node, error) {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ == durationType {
		duration, err := time.ParseDuration(value)
		if err != nil {
			return nil, fmt.Errorf("%q is not a duration such as 30s, 5m or 3h", value)
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: duration.String()}, nil
	}

	switch typ.Kind() {
	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("%q is not a boolean", value)
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(parsed)}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(value, 10, typ.Bits())
		if err != nil {
			return nil, fmt.Errorf("%q is not a valid %s", value, typ.Kind())
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(parsed, 10)}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(value, 10, typ.Bits())
		if err != nil {
			return nil, fmt.Errorf("%q is not a valid %s", value, typ.Kind())
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatUint(parsed, 10)}, nil
	case reflect.String:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}, nil
	default:
		return nil, fmt.Errorf("unsupported config value type %s", typ)
	}
}
