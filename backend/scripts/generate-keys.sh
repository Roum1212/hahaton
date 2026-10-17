#!/usr/bin/env sh
set -eu

printf 'SESSION_TOKEN_V4_SYMMETRIC_KEY=%s\n' "$(openssl rand -hex 32)"
printf 'DATA_ENCRYPTION_KEY=%s\n' "$(openssl rand -hex 32)"
printf 'ADMIN_API_KEY=%s\n' "$(openssl rand -hex 32)"
