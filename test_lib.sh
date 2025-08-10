#!/bin/bash

declare counterName
declare NO_COLOR
declare -a test_commands

counterName="${counterName:-"goenv-test"}"

GOOS_ENV=${GOOS:-$(go env GOOS)}
GOARCH_ENV=${GOARCH:-$(go env GOARCH)}
BIN_PATH="./bin/goenv-$GOOS_ENV-$GOARCH_ENV"

function red() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;31m${1}\033[0m"; else echo "${1}"; fi; }
function purple() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;34m${1}\033[0m";  else echo "${1}"; fi;}
function green() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;32m${1}\033[0m";  else echo "${1}"; fi;}
function magenta() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;35m${1}\033[0m";  else echo "${1}"; fi;}
function yellow() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[1;33m${1}\033[0m";  else echo "${1}"; fi;}
function cyan() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[1;36m${1}\033[0m";  else echo "${1}"; fi;}
function safe_exit(){
  echo -e "Fatal Error: ${1}"
  exit "${2:-1}"
}

function reset_sample_env(){
  {
    echo "AWS_REGION=us-west-2"
    echo "OUTPUT=json"
    echo "HOSTNAME=localhost"
    echo "DBUSER=readonly"
    echo "DBPASS=readonly"
    echo "DATABASE=test_data"
  } | tee "sample.env" > /dev/null
}

function build_goenv_binary() {
  echo "Building goenv binary..."

  [ -z "${BIN_PATH}" ] && { echo "invalid BIN_PATH in runtime"; exit 1; }

  [ -f "${BIN_PATH}" ] && rm -rf "${BIN_PATH}" && echo "Clean successful: ${BIN_PATH}"

  if ! go build -ldflags "-s -w" -o "$BIN_PATH" .; then
    echo "ERROR: Failed to build goenv binary. Aborting tests."
    exit 1
  fi
  echo "Build successful: $BIN_PATH"
  echo
}

function run() {
  local cmd="${1}"
  local -i testNo
  testNo=$(counter -name "${counterName}" -add || safe_exit "failed to increase counter")
  local -i exitCode
  local output

  if [[ "${cmd}" == -file* ]]; then
    # the user is providing a custom file
    cmd="$BIN_PATH ${1}"
  elif [[ "${cmd}" == -raw* ]]; then
    cmd="${1}"
    cmd=${cmd#"-raw "}
  else
    cmd="$BIN_PATH -file sample.env ${1}"
  fi

  # Capture both stdout and stderr together
  output=$(eval "${cmd}" 2>&1)
  exitCode=$?

  local prefix
  prefix="$(magenta "$(whoami)")@$(yellow goenv.git):$(purple "$(basename .)")"

  if (( exitCode == 0 )); then
    printf "%s ⚡ %s ⇒  %s\n" "${prefix}" "$(cyan "Test #${testNo}")" "$(green "${cmd}")"
    [[ -n "$output" ]] && printf "%s\n" "$output"
  else
    printf "%s ⚡ %s ⇒  %s\n" "${prefix}" "$(cyan "Test #${testNo}")" "$(red "${cmd}")"
    safe_exit "Command failed with exit code ${exitCode}:\noutput: ${output}"
  fi
}

function main(){
  build_goenv_binary
  reset_sample_env

  # Iterate through tests using array indices
  for cmd in "${test_commands[@]}"; do
    [[ -z $cmd || $cmd == \#* ]] && continue
    run "${cmd}"
  done

  rm sample.env
  rm non_existent_file.env
  rm new.env
  echo "All $(counter -name "${counterName}") tests PASS!"
}

