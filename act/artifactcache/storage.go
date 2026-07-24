// Copyright 2023 The Gitea Authors. All rights reserved.
// Copyright 2023 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Storage struct {
	rootDir string
}

func NewStorage(rootDir string) (*Storage, error) {
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return nil, err
	}
	return &Storage{
		rootDir: rootDir,
	}, nil
}

func (s *Storage) Exist(id uint64) (bool, error) {
	name := s.filename(id)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Storage) Write(id uint64, offset int64, reader io.Reader) error {
	name := s.tempName(id, offset)
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return err
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

// WriteBlock stores one staged block of an Azure block blob upload. Blocks may arrive in
// any order and are only put in sequence by OrderBlocks, once the client sends the block
// list that names them.
func (s *Storage) WriteBlock(id uint64, blockID string, reader io.Reader) error {
	name := s.blockName(id, blockID)
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return err
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

// OrderBlocks renames the staged blocks so that Commit concatenates them in the order the
// block list gives, which is the order the client wants them in, not the order they
// arrived in. A block the list does not name is dropped, as Azure drops it.
func (s *Storage) OrderBlocks(id uint64, blockIDs []string) error {
	staged := map[string]string{}
	for i, blockID := range blockIDs {
		name := s.blockName(id, blockID)
		if _, err := os.Stat(name); err != nil {
			return fmt.Errorf("block %q of cache %d was never uploaded: %w", blockID, id, err)
		}
		staged[name] = s.tempName(id, int64(i))
	}

	// Two passes: a block's target name may be another block's current name when the
	// list reorders them, so move everything aside first.
	for name := range staged {
		if err := os.Rename(name, name+".ordering"); err != nil {
			return err
		}
	}
	for name, target := range staged {
		if err := os.Rename(name+".ordering", target); err != nil {
			return err
		}
	}

	// Anything left is a staged block the list dropped.
	entries, err := os.ReadDir(s.tempDir(id))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), blockFilePrefix) {
			_ = os.Remove(filepath.Join(s.tempDir(id), entry.Name()))
		}
	}
	return nil
}

// AppendBlock appends to an append blob, whose parts are always written in order.
func (s *Storage) AppendBlock(id uint64, reader io.Reader) error {
	names, err := s.tempNames(id)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return s.Write(id, int64(len(names)), reader)
}

func (s *Storage) Commit(id uint64, size int64) (int64, error) {
	defer func() {
		_ = os.RemoveAll(s.tempDir(id))
	}()

	name := s.filename(id)
	tempNames, err := s.tempNames(id)
	if err != nil {
		return 0, err
	}

	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return 0, err
	}
	file, err := os.Create(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var written int64
	for _, v := range tempNames {
		f, err := os.Open(v)
		if err != nil {
			return 0, err
		}
		n, err := io.Copy(file, f)
		_ = f.Close()
		if err != nil {
			return 0, err
		}
		written += n
	}

	// If size is less than 0, it means the size is unknown.
	// We can't check the size of the file, just skip the check.
	// It happens when the request comes from old versions of actions, like `actions/cache@v2`.
	if size >= 0 && written != size {
		_ = file.Close()
		_ = os.Remove(name)
		return 0, fmt.Errorf("broken file: %v != %v", written, size)
	}

	return written, nil
}

func (s *Storage) Serve(w http.ResponseWriter, r *http.Request, id uint64) {
	name := s.filename(id)
	http.ServeFile(w, r, name)
}

func (s *Storage) Remove(id uint64) {
	_ = os.Remove(s.filename(id))
	_ = os.RemoveAll(s.tempDir(id))
}

func (s *Storage) filename(id uint64) string {
	return filepath.Join(s.rootDir, fmt.Sprintf("%02x", id%0xff), strconv.FormatUint(id, 10))
}

func (s *Storage) tempDir(id uint64) string {
	return filepath.Join(s.rootDir, "tmp", strconv.FormatUint(id, 10))
}

func (s *Storage) tempName(id uint64, offset int64) string {
	return filepath.Join(s.tempDir(id), fmt.Sprintf("%016x", offset))
}

// blockFilePrefix marks a staged, not yet ordered block, so that Commit's name-ordered
// concatenation never picks one up and OrderBlocks can clean up the leftovers.
const blockFilePrefix = "block-"

func (s *Storage) blockName(id uint64, blockID string) string {
	// The block id is client-chosen (base64), so it is hashed rather than trusted as a
	// path element.
	sum := sha256.Sum256([]byte(blockID))
	return filepath.Join(s.tempDir(id), blockFilePrefix+hex.EncodeToString(sum[:]))
}

func (s *Storage) tempNames(id uint64) ([]string, error) {
	dir := s.tempDir(id)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, v := range files {
		if !v.IsDir() {
			names = append(names, filepath.Join(dir, v.Name()))
		}
	}
	return names, nil
}
