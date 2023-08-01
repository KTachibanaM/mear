#!/usr/bin/env bash
set -e

USERNAME=KTachibanaM
REPO=mear

if command -v curl >/dev/null 2>&1; then
  echo "curl is installed"
else
  echo "curl is not installed"
  exit 1
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
  os=darwin
  if [[ "$(uname -m)" == "x86_64" ]]; then
    arch=amd64
  elif [[ "$(uname -m)" == "arm64" ]]; then
    arch=arm64
  else
    echo "unsupported cpu architecture $(uname -m). you may find binary for your os/architecture at https://mear.cloud/releases"
    exit 1
  fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
  os=linux
  if [[ "$(uname -m)" == "x86_64" ]]; then
    arch=amd64
  else
    echo "unsupported cpu architecture $(uname -m). you may find binary for your os/architecture at https://mear.cloud/releases"
    exit 1
  fi
else
  echo "unsupported os $OSTYPE. you may find binary for your os/architecture at https://mear.cloud/releases"
  exit 1
fi

version=$(curl --silent "https://api.github.com/repos/${USERNAME}/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
version="${version:1}"

echo "installing mear-host ${version} to /usr/local/bin/..."
sudo curl --silent -L "https://github.com/${USERNAME}/${REPO}/releases/download/v${version}/mear-host_${version}_${os}_${arch}" -o "/usr/local/bin/mear-host"
sudo chmod +x "/usr/local/bin/mear-host"

echo "installing mear ${version} to /usr/local/bin/..."
sudo curl --silent -L "https://github.com/${USERNAME}/${REPO}/releases/download/v${version}/mear_${version}_${os}_${arch}" -o "/usr/local/bin/mear"
sudo chmod +x "/usr/local/bin/mear"
