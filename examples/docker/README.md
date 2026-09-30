### Run `gitea-runner` in a Docker Container

```sh
docker run --privileged \
  -e GITEA_INSTANCE_URL=http://192.168.8.18:3000 \
  -e GITEA_RUNNER_REGISTRATION_TOKEN=<runner_token> \
  -v "$PWD/data:/data" \
  -v runner-docker:/var/lib/docker \
  --name my_runner gitea/runner:nightly-dind
```

The `/data` volume preserves the runner's registration, and `/var/lib/docker` keeps the images pulled by its bundled Docker daemon.

To use the host's Docker daemon instead, mount its socket into the basic image. This gives jobs control over the host's daemon.

```sh
docker run -e GITEA_INSTANCE_URL=http://192.168.8.18:3000 -e GITEA_RUNNER_REGISTRATION_TOKEN=<runner_token> -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD/data:/data" --name my_runner gitea/runner:nightly
```

See [image flavours](../../README.md#image-flavours) for the rootless variant.
