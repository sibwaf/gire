#!/bin/sh

# Create user for Git SSH server

if [ -z "$GIRE_GIT_USER" ]; then
    git_user_name="git"
else
    git_user_name="$GIRE_GIT_USER"
fi

adduser -D -s /bin/sh "$git_user_name"

git_user_home=$(su "$git_user_name" -s /bin/sh -c "echo \$HOME")

# Start SSH-agent and import identity keys for current user for pulling repositories to Gire

eval $(ssh-agent) > /dev/null

for f in $(find /keys -type f -maxdepth 1); do
    ssh-add "$f"
done

# Configure known_hosts for current user for pulling repositories to Gire

mkdir -p "/keys/trusted"
mkdir -p "$HOME/.ssh"

for i in $(seq 1 99); do
    entry=$(eval echo "\$GIRE_SSH_TRUSTED_$i")
    if [ -z "$entry" ]; then
        continue
    fi
    
    echo "$entry" >> "$HOME/.ssh/known_hosts"
done
for f in $(find /keys/trusted -type f -maxdepth 1); do
    cat "$f" >> "$HOME/.ssh/known_hosts"
done

# Configure authorized_keys for the git user for pulling repositories from Gire

mkdir -p "/keys/authorized"
mkdir -p "$git_user_home/.ssh"

for i in $(seq 1 99); do
    entry=$(eval echo "\$GIRE_SSH_AUTHORIZED_$i")
    if [ -z "$entry" ]; then
        continue
    fi
    
    echo "$entry" >> "$git_user_home/.ssh/authorized_keys"
done
for f in $(find /keys/authorized -type f -maxdepth 1); do
    cat "$f" >> "$git_user_home/.ssh/authorized_keys"
done

# Expose configuration environment variables to SSH connections

echo "$GIRE_REPOSITORY_PATH" > /run/gire_repository_path

# Start SSHD in background

sshd_keys=""

for f in $(find /keys -type f -maxdepth 1); do
    sshd_keys="$sshd_keys -h $f"
done

su "$git_user_name" -s /bin/sh -c "/usr/sbin/sshd $sshd_keys"

# Start the application

/app/gire $@
