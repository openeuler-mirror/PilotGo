#!/bin/bash
set -e -o pipefail
ERR=false
FAIL=false

for file in $(git ls-files | grep "\.go$" | grep -v vendor/); do
  echo -n "Header check: $file... "
  if [[ -z $(cat ${file} | grep  "Copyright (c) KylinSoft Co., Ltd.") && -z $(cat ${file} | grep "PilotGo is licensed under the Mulan PSL v2.") ]]; then
      ERR=true
  fi
  if [ $ERR == true ]; then
    if [[ $# -gt 0 && $1 =~ [[:upper:]fix] ]]; then
      cat ./scripts/licence/boilerplate.go.txt "${file}" > "${file}".new
      mv "${file}".new "${file}"
      echo "$(tput -T xterm setaf 3)FIXING$(tput -T xterm sgr0)"
      ERR=false
    else
      echo "$(tput -T xterm setaf 1)FAIL$(tput -T xterm sgr0)"
      ERR=false
      FAIL=true
    fi
  else
    echo "$(tput -T xterm setaf 2)OK$(tput -T xterm sgr0)"
  fi
done

# If we failed one check, return 1
[ $FAIL == true ] && exit 1 || exit 0