#!/usr/bin/env bash
set -euo pipefail

release_directory=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
deployment_directory=$(dirname -- "$(dirname -- "$release_directory")")
revision=$(basename -- "$release_directory")
env_file=${1:?Pass the path to the server .env}
previous_compose="$deployment_directory/current/docker-compose.yaml"
if [[ ! -f "$previous_compose" ]]; then
  previous_compose="$(dirname -- "$env_file")/docker-compose.yaml"
fi

previous_backend=$(docker ps -aq --filter label=com.docker.compose.project=hahaton \
  --filter label=com.docker.compose.service=backend --filter label=com.docker.compose.oneoff=False | head -1)
previous_frontend=$(docker ps -aq --filter label=com.docker.compose.project=hahaton \
  --filter label=com.docker.compose.service=frontend --filter label=com.docker.compose.oneoff=False | head -1)
if [[ -n "$previous_backend" ]]; then
  previous_backend=$(docker inspect --format '{{.Image}}' "$previous_backend")
fi
if [[ -n "$previous_frontend" ]]; then
  previous_frontend=$(docker inspect --format '{{.Image}}' "$previous_frontend")
fi

export BACKEND_IMAGE="hahaton-backend:$revision"
export FRONTEND_IMAGE="hahaton-frontend:$revision"
compose=(docker compose --project-name hahaton --env-file "$env_file" -f "$release_directory/docker-compose.yaml")
applications_stopped=false

restore_applications() {
  local exit_code=$?
  trap - ERR
  if [[ "$applications_stopped" == true && -n "$previous_backend" && -n "$previous_frontend" ]]; then
    echo 'Deployment failed; attempting to restart the previous application images.' >&2
    BACKEND_IMAGE="$previous_backend" FRONTEND_IMAGE="$previous_frontend" \
      docker compose --project-name hahaton --env-file "$env_file" -f "$previous_compose" \
      up -d --no-build --pull never --no-deps --wait --wait-timeout 180 backend frontend || true
  fi
  echo "Deployment failed. Database migrations are not rolled back automatically: $release_directory" >&2
  exit "$exit_code"
}

trap restore_applications ERR

test -r "$env_file"
"${compose[@]}" config --quiet
docker load --input "$release_directory/images.tar.gz"
"${compose[@]}" up -d --no-build --pull never --wait --wait-timeout 120 postgres

applications_stopped=true
"${compose[@]}" stop frontend backend
umask 077
"${compose[@]}" exec -T postgres sh -c 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' \
  > "$release_directory/database.dump"
"${compose[@]}" run --rm --no-deps --pull never --entrypoint /app/goose backend up
"${compose[@]}" up -d --no-build --pull never --wait --wait-timeout 180 backend frontend

ln -sfn "$release_directory" "$deployment_directory/current"
rm -f "$release_directory/images.tar.gz"
echo "Deployed $revision; application healthchecks passed."
