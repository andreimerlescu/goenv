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

declare -a test_commands=(
  "-has -env HOSTNAME"
  "-has -env NON_EXISTENT"
  "-is -env DATABASE -value test_data"
  "-is -env DATABASE -value wrong_data"
  "-print"
  "-json"
  "-yaml"
  "-toml"
  "-ini"
  "-xml"
  "-write -add -env NEW_KEY -value 'a new value'"
  "-has -env NEW_KEY"
  "-write -add -env HOSTNAME -value 'another-host'"
  "-is -env HOSTNAME -value localhost"
  "-write -rm -env OUTPUT"
  "-not -has -env OUTPUT"
  "-file new.env -add -env HELLO -value world -write"
  "-file new.env -is -env HELLO -value world"
  "-v"
  "-file non_existent_file.env -add -env FOO -value bar || echo \"Test success because we expected an error here.\""
  "-file non_existent_file.env -add -env FOO -value bar -write"
)

export test_commands

source ./test_lib.sh || { echo "Missing test_lib.sh"; exit 1; }

main "$@"
