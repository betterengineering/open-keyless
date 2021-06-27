#!/bin/bash

set -e

echo "Determining latest version"
LATEST=$(curl --silent "https://api.github.com/repos/betterengineering/open-keyless/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

echo "Fetching open-keyless-controller version ${LATEST}"
mkdir -p "/tmp/open-keyless/${LATEST}"
curl -Lo "/tmp/open-keyless/${LATEST}/open-keyless-controller_${LATEST}_armhf.deb" "https://github.com/betterengineering/open-keyless/releases/download/${LATEST}/open-keyless-controller_${LATEST}_armhf.deb"

# Check config file.
CONFIG_FILE=/etc/open-keyless-controller/config.yml
if [ -f "$CONFIG_FILE" ]; then
    echo "$CONFIG_FILE exists, skipping update."
else 
    echo "$CONFIG_FILE does not exist, installing."
    curl -Lo "/tmp/open-keyless/${LATEST}/config.yml" "https://github.com/betterengineering/open-keyless/releases/download/${LATEST}/config.yml"
    mkdir -p /etc/open-keyless-controller
    cp "/tmp/open-keyless/${LATEST}/config.yml" "${CONFIG_FILE}"
fi

echo "Installing open-keyless-controller version ${LATEST}"
sudo dpkg -i "/tmp/open-keyless/${LATEST}/open-keyless-controller_${LATEST}_armhf.deb"

echo "Restarting open-keyless-controller"
# TODO - figure out if we can rely on dpkg to do this for us.
systemctl daemon reload
systemctl enable open-keyless-controller
systemctl restart open-keyless-controller