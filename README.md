# Gire

Gire ("git reflector", read as "gear") is a minimalistic read-only Git repository mirroring service for self-hosting.

## Configruation

```
GIRE_GIT_USER="git"
GIRE_SCAN_CRON="@daily"
GIRE_SOURCES_PATH="sources.yaml"
GIRE_REPOSITORY_PATH="repositories"
```

## Setting up sources

`sources.yaml`
```yaml
- groupName: gh-source-1
  url: https://github.com/sibwaf
  type: github
  include:
    - ".*"
  exclude:
    - "test"

- groupName: _
  type: repository
  url: git@github.com:sibwaf/gire.git
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
