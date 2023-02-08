# Gire

Gire ("git reflector", read as "gear") is a minimalistic read-only Git repository mirroring service for self-hosting.

## Quickstart

```sh
# Clone and build

git clone git@github.com:sibwaf/gire

cd gire

docker build -t gire .

# Prepare configuration

mkdir -p run

ssh-keygen -t ed25519 -f run/gire_ed25519 -P ""
sudo chown git:git run/gire_ed25519
sudo chmod 755 run/gire_ed25519

ssh-keyscan github.com > run/github.pub

ssh-add -L > run/me.pub

echo "- url: git@github.com:sibwaf/gire" > run/sources.yaml

# Run in background

docker run --rm -d --name gire \
    -v "$(pwd)/run/repositories:/app/repositories" \
    -v "$(pwd)/run/sources.yaml:/app/sources.yaml:ro" \
    -v "$(pwd)/run/gire_ed25519:/keys/id:ro" \
    -v "$(pwd)/run/github.pub:/keys/trusted/github:ro" \
    -v "$(pwd)/run/me.pub:/keys/authorized/me:ro" \
    -e GIRE_CRON="0 * * * * *" \
    -p 122:22 \
    gire

# Watch the logs to catch the moment when repository gets mirrored

docker logs -f gire

# Try cloning it back from Gire

git clone git@[localhost:122]:_/gire run/gire
```

## Configruation

Gire is configured using only environment variables as it is supposed to be always running in a container.

Available options with their default values:
```sh
# Unix user inside the container which will be used for providing access
# to Gire's repositories over SSH: GIRE_GIT_USER@address:repository
#
# Used only on container start - see entrypoint.sh
GIRE_GIT_USER="git"

# Cron expression for source scanning / repository updates
# Format: seconds minutes hours day-of-month month day-of-week
#
# See https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc#hdr-CRON_Expression_Format
GIRE_CRON="@daily"

# Path to source configuration file
GIRE_SOURCES_PATH="sources.yaml"

# Path to directory where all repositories will be stored
GIRE_REPOSITORY_PATH="repositories"
```

Default working directory in the container is `/app`

## Setting up sources

`sources.yaml`
```yaml
- url: git@github.com:sibwaf/gire.git

- url: https://github.com/sibwaf
  type: github
  groupName: gh-source-1
  authToken: 12345
  include:
    - ".*"
  exclude:
    - "test"
```

All properties are optional excluding `url`.

- `url` - URL to the repository, behavior differs based on `type`
- `type`
  - `type: repository` - default value, supports any `url` value which can be used in `git clone $URL`
  - `type: github` - integration with GitHub to effortlessly mirror multiple repositories at once
- `groupName` - subdirectory of `GIRE_REPOSITORY_PATH` to store all repositories from this source; for `type: repository` default value is `_`
- `authToken` - authentication token for integrations, supports environment variable expansion: `authToken: $GITHUB_TOKEN`
- `include` - regular expressions for filtering URLs in multi-repository sources; empty `include` means "take everything"
- `exclude` - regular expressions for excluding URLs in multi-repository sources; takes priority over `include`

All expressions in `include` / `exclude` lists are joined using `OR`:
```yaml
# If the source provides URLs test0, test1, test2:
# test0 - ignore (missing in includes)
# test1 - mirror (present in includes, missing in excludes)
# test2 - ignore (exclude takes priority)

include:
  - "test1"
  - "test2"
exclude:
  - "test2"
```

## Setting up SSH keys

The only authentication Gire supports is by using SSH keys.

### Host keys

Mount your ***private*** key files to `/keys` in the container. Those *must be* owned by the git user (see GIRE_GIT_USER) and have correct permissions for SSH (755 or less).

Example: `-v /home/USER/gire_key:/keys/id_ed25519`

You can mount multiple files (`id_rsa`, `id_ed25519`, ...) at the same time - all of them will be used when starting sshd.

### Trusted keys / known_hosts

Note: you can get the key using `ssh-keyscan HOSTNAME` which will result in one or more lines of format `HOSTNAME ALGORIGHM KEY`.

Cloning repositories into Gire with SSH requires having the server key in the `$HOME/.ssh/known_hosts`. Add them using:
1. Environment variables `GIRE_SSH_TRUSTED_1` .. `GIRE_SSH_TRUSTED_99`: `-e GIRE_SSH_TRUSTED_1="example.com ssh-rsa ABCDE"`
2. File mounts into `/keys/trusted`. Each file might contain any number of keys, filenames don't matter.

All of those keys will be appended into `$HOME/.ssh/known_hosts` at start.

### Authorized keys / authorized_keys

Note: this section assumes default value `git` for the `GIRE_GIT_USER` configuration variable. If you change it, the `/home/git` path should be replaced with the suitable home path. It's handled automatically when using the following configuration options.

Cloning repositories from Gire requires having client keys in the `/home/git/.ssh/authorized_keys`. Add them using:
1. Environment variables `GIRE_SSH_AUTHORIZED_1` .. `GIRE_SSH_AUTHORIZED_99`: `-e GIRE_SSH_AUTHORIZED_1="ssh-rsa ABCDE user@example.com"`
2. File mounts into `/keys/authorized`. Each file might contain any number of keys, filenames don't matter.

All of those keys will be appended into `/home/git/.ssh/authorized_keys` at start.

## Integrations

### GitHub

Sources with `type: github` will use the provided URL to list all available repositories and synchronize all of them. Both users and organizations are supported.

To allow *listing* private repositories and/or raise the request rate limit you can provide your personal access token using the `authToken` property. Notice: Gire will pull all repositories using SSH, so you still need to have SSH keys configured to pull private repositories.

The URL is also used as prefix filter, so if you have permissions for repositories `u1/aaaa` and `u2/bbbb`:
- `url: https://github.com/u1` will synchronize `u1/aaaa` only
- `url: https://github.com/u2` will synchronize `u2/bbbb` only

By default, `groupName` value will be automatically extracted from the URL: `https://github.com/USERNAME` -> `USERNAME`. You can override it by providing your own value in the source configuration.

See also:
- [https://docs.github.com/en/rest/repos/repos#list-repositories-for-a-user]
- [https://docs.github.com/en/rest/repos/repos#list-repositories-for-the-authenticated-user]
