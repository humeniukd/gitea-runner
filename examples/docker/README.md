### Run `gitea-runner` in a Docker Container

Prefer Docker-in-Docker for container jobs so workflows use a dedicated Docker daemon instead of controlling the host daemon:

```sh
docker run --privileged \
  -e GITEA_INSTANCE_URL=http://192.168.8.18:3000 \
  -e GITEA_RUNNER_REGISTRATION_TOKEN=<runner_token> \
  -v "$PWD/data:/data" \
  -v runner-docker:/var/lib/docker \
  --name my_runner gitea/runner:nightly-dind
```

The `/data` volume preserves the runner's registration. The separate `/var/lib/docker` volume preserves the bundled daemon's images across container recreation.

The `dind` flavour requires a privileged container and is not a sandbox for untrusted code. Consider `nightly-dind-rootless` to run the bundled daemon as an unprivileged user; it has additional host requirements and rootless Docker limitations. See the [image flavours](../../README.md#image-flavours) for details and the correct image-cache path for each flavour.

If trusted workflows need the host Docker daemon, use the basic image with its socket mounted:

```sh
docker run -e GITEA_INSTANCE_URL=http://192.168.8.18:3000 -e GITEA_RUNNER_REGISTRATION_TOKEN=<runner_token> -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD/data:/data" --name my_runner gitea/runner:nightly
```

Access to the host socket gives workflows control over the host Docker daemon and can let them modify the host. The runner's automatic cleanup also removes resources created through its job Docker proxy; see [Docker from a job](../../README.md#docker-from-a-job-gitea_docker_workspace).
