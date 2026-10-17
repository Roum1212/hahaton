#!/usr/bin/env bash
set -euo pipefail

revision=${SSH_ORIGINAL_COMMAND:-}
if [[ ! "$revision" =~ ^[0-9a-f]{40}$ ]]; then
  echo 'Expected a full Git commit SHA.' >&2
  exit 1
fi

deployment_directory=/root/hahaton-deploy
env_file=/root/hahaton/.env
mkdir -p "$deployment_directory/releases"
exec 9> "$deployment_directory/deploy.lock"
flock -w 1800 9

release_directory="$deployment_directory/releases/$revision"
umask 077
mkdir -p "$release_directory"
archive=$(mktemp "$deployment_directory/incoming.XXXXXX")
trap 'rm -f "$archive"' EXIT
cat > "$archive"

python3 - "$archive" "$release_directory" <<'PY'
import pathlib
import sys
import tarfile

with tarfile.open(sys.argv[1], "r:") as archive:
    members = archive.getmembers()
    expected = {"images.tar.gz", "docker-compose.yaml", "deploy.sh"}
    if len(members) != len(expected) or {member.name for member in members} != expected:
        raise SystemExit("invalid release files")
    if any(not member.isfile() for member in members):
        raise SystemExit("release must contain regular files")
    for member in members:
        destination = pathlib.Path(sys.argv[2]) / member.name
        with archive.extractfile(member) as source, destination.open("wb") as target:
            while chunk := source.read(1024 * 1024):
                target.write(chunk)
PY

bash "$release_directory/deploy.sh" "$env_file"
