package container

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestImageExistsLocally(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	// to help make this test reliable and not flaky, we need to have
	// an image that will exist, and onew that won't exist

	// Test if image exists with specific tag
	invalidImageTag, err := ImageExistsLocally(ctx, "library/alpine:this-random-tag-will-never-exist", "linux/amd64")
	assert.NoError(t, err)
	assert.False(t, invalidImageTag)

	// Test if image exists with specific architecture (image platform)
	invalidImagePlatform, err := ImageExistsLocally(ctx, "alpine:latest", "windows/amd64")
	assert.NoError(t, err)
	assert.False(t, invalidImagePlatform)

	// pull an image
	cli, err := client.NewClientWithOpts(client.FromEnv)
	assert.NoError(t, err)
	cli.NegotiateAPIVersion(context.Background())

	// Chose alpine latest because it's so small
	// maybe we should build an image instead so that tests aren't reliable on dockerhub
	readerDefault, err := cli.ImagePull(ctx, "node:16-buster-slim", image.PullOptions{
		Platform: "linux/amd64",
	})
	assert.NoError(t, err)
	defer readerDefault.Close()
	_, err = io.ReadAll(readerDefault)
	assert.NoError(t, err)

	imageDefaultArchExists, err := ImageExistsLocally(ctx, "node:16-buster-slim", "linux/amd64")
	assert.NoError(t, err)
	assert.True(t, imageDefaultArchExists)
}

func TestImageExistsLocallyQemu(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	if _, ok := os.LookupEnv("NO_QEMU"); ok {
		t.Skip("skipping test because QEMU is disabled")
	}

	ctx := context.Background()

	// pull an image
	cli, err := client.NewClientWithOpts(client.FromEnv)
	assert.NoError(t, err)
	cli.NegotiateAPIVersion(context.Background())

	// Validate if another architecture platform can be pulled
	readerArm64, err := cli.ImagePull(ctx, "node:16-buster-slim", image.PullOptions{
		Platform: "linux/arm64",
	})
	assert.NoError(t, err)
	defer readerArm64.Close()
	_, err = io.ReadAll(readerArm64)
	assert.NoError(t, err)

	imageArm64Exists, err := ImageExistsLocally(ctx, "node:16-buster-slim", "linux/arm64")
	assert.NoError(t, err)
	assert.True(t, imageArm64Exists)
}
