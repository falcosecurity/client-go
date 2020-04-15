# Release Process

When we release we do the following process:

1. We decide together (usually in the #falco channel in [slack](https://sysdig.slack.com)) what's the [next version](#About-versioning) to tag
2. A person with repository rights does the tag
3. The same person runs commands in their machine following the "Release commands" section below
4. The tag is live on [Github](https://github.com/falcosecurity/client-go/releases) with the changelog attached

## Release commands

Just tag the [version](#About-versioning). For example:

```bash
git tag -a v0.1.0-rc.0 -m "v0.1.0-rc.0"
git push origin v0.1.0-rc.0
```

The [goreleaser](https://goreleaser.com/ci/) will run on CircleCI and do the magic :)

## About versioning

Basically, this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Furthermore, although the `go-client` version is NOT paired with the Falco version directly, the major version MUST be incremented when any backwards incompatible changes are introduced to the proto definitions. The only notable exception is [major version zero](https://semver.org/spec/v2.0.0.html#spec-item-4), indeed, for initial development just the minor version MUST be incremented.