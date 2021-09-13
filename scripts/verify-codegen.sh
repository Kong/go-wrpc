#!/bin/bash

function is_dirty_repo() {
  DIRTY=$(git status --porcelain 2> /dev/null | wc -l)
  if [[ $DIRTY -eq 0 ]];
  then
    return 0
  else
    return 1
  fi
}

is_dirty_repo
if [[ $? -ne 0 ]];
then
  echo "Dirty repository, please commit changes and try again"
  exit 1
fi

buf generate
if [[ $? -ne 0 ]];
then
  echo "failed to run 'buf generate'"
  exit 1
fi

is_dirty_repo
if [[ $? -ne 0 ]];
then
  echo "generated code is not up to date"
  exit 1
fi

