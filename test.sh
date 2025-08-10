#!/bin/bash
# shellcheck disable=SC2086

go install github.com/andreimerlescu/counter@latest || safe_exit "failed to install counter"

[ ! -d .counters ] && { mkdir -p .counters || safe_exit "failed to mkdir .counters"; }

declare counterName
declare COUNTER_DIR
declare COUNTER_USE_FORCE
declare COUNTER_ALWAYS_YES
declare NO_COLOR

COUNTER_DIR=$(realpath .counters)
COUNTER_USE_FORCE=1
COUNTER_ALWAYS_YES=1
NO_COLOR="${NO_COLOR:=false}"
counterName="goenv-test"

export COUNTER_USE_FORCE
export COUNTER_ALWAYS_YES
export COUNTER_DIR
export counterName

counter -name $counterName -reset -yes 1> /dev/null || safe_exit "failed to reset counter"

declare input_file
input_file="${1:-test_cmds.txt}"

declare -a test_commands=()

if [ -n "${input_file}" ] && [ -f "${input_file}" ]; then

  while IFS= read -r cmd; do
      [[ -z $cmd || $cmd == \#* ]] && continue
      test_commands+=("${cmd}")
  done < "$input_file"

fi

export test_commands

source ./test_lib.sh || { echo "Missing test_lib.sh"; exit 1; }

main "$@"
