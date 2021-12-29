#!/usr/bin/env bash

# ------------------------------------------------------------
# Copyright 2021 The tKeel Contributors.
# Licensed under the Apache License.
# ------------------------------------------------------------

# tKeel CLI location
: ${TKEEL_INSTALL_DIR:="/usr/local/bin"}

# sudo is required to copy binary to TKEEL_INSTALL_DIR for linux
: ${USE_SUDO:="false"}

# Http request CLI
TKEEL_HTTP_REQUEST_CLI=curl

# GitHub Organization and repo name to download release
GITHUB_ORG=tkeel-io
GITHUB_REPO=cli

# tKeel CLI filename
TKEEL_CLI_FILENAME=tk

TKEEL_CLI_FILE="${TKEEL_INSTALL_DIR}/${TKEEL_CLI_FILENAME}"

getSystemInfo() {
    ARCH=$(uname -m)
    case $ARCH in
        armv7*) ARCH="arm";;
        aarch64) ARCH="arm64";;
        x86_64) ARCH="amd64";;
    esac

    OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

    # Most linux distro needs root permission to copy the file to /usr/local/bin
    if [[ "$OS" == "linux" || "$OS" == "darwin" ]] && [ "$TKEEL_INSTALL_DIR" == "/usr/local/bin" ]; then
        USE_SUDO="true"
    fi
}

verifySupported() {
    local supported=(darwin-amd64 linux-amd64 linux-arm linux-arm64)
    local current_osarch="${OS}-${ARCH}"

    for osarch in "${supported[@]}"; do
        if [ "$osarch" == "$current_osarch" ]; then
            echo "Your system is ${OS}_${ARCH}"
            return
        fi
    done

    if [ "$current_osarch" == "darwin-arm64" ]; then
        echo "The darwin_arm64 arch has no native binary, however you can use the amd64 version so long as you have rosetta installed"
        echo "Use 'softwareupdate --install-rosetta' to install rosetta if you don't already have it"
        ARCH="amd64"
        return
    fi


    echo "No prebuilt binary for ${current_osarch}"
    exit 1
}

runAsRoot() {
    local CMD="$*"

    if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
        CMD="sudo $CMD"
    fi

    $CMD
}

checkHttpRequestCLI() {
    if type "curl" > /dev/null; then
        TKEEL_HTTP_REQUEST_CLI=curl
    elif type "wget" > /dev/null; then
        TKEEL_HTTP_REQUEST_CLI=wget
    else
        echo "Either curl or wget is required"
        exit 1
    fi
}

checkExistingTKeel() {
    if [ -f "$TKEEL_CLI_FILE" ]; then
        echo -e "\ntKeel CLI is detected:"
        $TKEEL_CLI_FILE --version
        echo -e "Reinstalling tKeel CLI - ${TKEEL_CLI_FILE}...\n"
    else
        echo -e "Installing tKeel CLI...\n"
    fi
}

getLatestRelease() {
    local tkeelReleaseUrl="https://api.github.com/repos/${GITHUB_ORG}/${GITHUB_REPO}/releases"
    local latest_release=""

    if [ "$TKEEL_HTTP_REQUEST_CLI" == "curl" ]; then
        latest_release=$(curl -s $tkeelReleaseUrl | grep \"tag_name\" | grep -v rc | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    else
        latest_release=$(wget -q --header="Accept: application/json" -O - $tkeelReleaseUrl | grep \"tag_name\" | grep -v rc | awk 'NR==1{print $2}' |  sed -n 's/\"\(.*\)\",/\1/p')
    fi

    ret_val=$latest_release
}

downloadFile() {
    LATEST_RELEASE_TAG=$1

    TKEEL_CLI_ARTIFACT="${TKEEL_CLI_FILENAME}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_BASE="https://github.com/${GITHUB_ORG}/${GITHUB_REPO}/releases/download"
    DOWNLOAD_URL="${DOWNLOAD_BASE}/${LATEST_RELEASE_TAG}/${TKEEL_CLI_ARTIFACT}"

    # Create the temp directory
    TKEEL_TMP_ROOT=$(mktemp -dt tkeel-install-XXXXXX)
    ARTIFACT_TMP_FILE="$TKEEL_TMP_ROOT/$TKEEL_CLI_ARTIFACT"

    echo "Downloading $DOWNLOAD_URL ..."
    if [ "$TKEEL_HTTP_REQUEST_CLI" == "curl" ]; then
        curl -SsL "$DOWNLOAD_URL" -o "$ARTIFACT_TMP_FILE"
    else
        wget -q -O "$ARTIFACT_TMP_FILE" "$DOWNLOAD_URL"
    fi

    if [ ! -f "$ARTIFACT_TMP_FILE" ]; then
        echo "failed to download $DOWNLOAD_URL ..."
        exit 1
    fi
}

installFile() {
    tar xf "$ARTIFACT_TMP_FILE" -C "$TKEEL_TMP_ROOT"
    local tmp_root_tkeel_cli="$TKEEL_TMP_ROOT/$TKEEL_CLI_FILENAME"

    if [ ! -f "$tmp_root_tkeel_cli" ]; then
        echo "Failed to unpack tKeel CLI executable."
        exit 1
    fi

    chmod o+x $tmp_root_tkeel_cli
    runAsRoot cp "$tmp_root_tkeel_cli" "$TKEEL_INSTALL_DIR"

    if [ -f "$TKEEL_CLI_FILE" ]; then
        echo "$TKEEL_CLI_FILENAME installed into $TKEEL_INSTALL_DIR successfully."

        $TKEEL_CLI_FILE --version
    else 
        echo "Failed to install $TKEEL_CLI_FILENAME"
        exit 1
    fi
}

fail_trap() {
    result=$?
    if [ "$result" != "0" ]; then
        echo "Failed to install tKeel CLI"
        echo "For support, go to https://tkeel.io"
    fi
    cleanup
    exit $result
}

cleanup() {
    if [[ -d "${TKEEL_TMP_ROOT:-}" ]]; then
        rm -rf "$TKEEL_TMP_ROOT"
    fi
}

installCompleted() {
    echo -e "\nTo get started with tKeel, please visit https://docs.tkeel.io/getting-started/"
}

# -----------------------------------------------------------------------------
# main
# -----------------------------------------------------------------------------
trap "fail_trap" EXIT

getSystemInfo
verifySupported
checkExistingTKeel
checkHttpRequestCLI


if [ -z "$1" ]; then
    echo "Getting the latest tKeel CLI..."
    getLatestRelease
else
    ret_val=v$1
fi

echo "Installing $ret_val tKeel CLI..."

downloadFile $ret_val
installFile
cleanup

installCompleted
