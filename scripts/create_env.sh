#! /bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TEMPLATE_FILE=$SCRIPT_DIR/template.env
ENV_FILE=$SCRIPT_DIR/../.env

cp $TEMPLATE_FILE $ENV_FILE
