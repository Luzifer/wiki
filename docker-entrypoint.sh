#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail

chown wiki: "${DATA_DIR}"

exec gosu wiki /usr/local/bin/wiki
