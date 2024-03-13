#!/bin/bash

if grep -q "export PATH=\$PATH:$(pwd)" ~/.bashrc; then
  echo "PATH already set"
else
  echo "export PATH="\$PATH:$(pwd)"" >> ~/.bashrc
  export PATH=$PATH:$(pwd)
  echo "Successfully added PATH"
fi

if grep -q "export LATER_PROJECT_DIR=\"$(pwd)/\"" ~/.bashrc; then
  echo "LATER_PROJECT_DIR already set"
else
  echo "export LATER_PROJECT_DIR=\"$(pwd)/\"" >> ~/.bashrc
  echo "Successfully added LATER_PROJECT_DIR"
fi
