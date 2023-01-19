#!/bin/sh

case "$SSH_ORIGINAL_COMMAND" in
(git-upload-pack*)
    base_dir=$(cat /run/gire_repository_path)
    if [ -z "$base_dir" ]; then
        base_dir="$HOME"
    fi

    sh -c "cd $base_dir; $SSH_ORIGINAL_COMMAND"
    ;;
(git-receive-pack*)
    exit 1
    ;;
(*)
    exit 1
    ;;
esac
