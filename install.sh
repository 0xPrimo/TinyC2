#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "Run this script with sudo."
  exit 1
fi

if [ -z "$1" ]; then
  echo "Usage: sudo ./install.sh /path/to/dist"
  exit 1
fi

TARGET_BIN=$(realpath "$1/link")
TARGET_DIR=$(dirname "$TARGET_BIN")
BIN_NAME=$(basename "$TARGET_BIN")
DEST_PATH="/usr/local/bin/crystal-link"

cat << EOF > "$DEST_PATH"
#!/bin/bash
cd "$TARGET_DIR"
exec "./$BIN_NAME" "\$@"
EOF
chmod +x "$DEST_PATH"

