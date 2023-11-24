#!/bin/sh

set -e

envsubst < /app/config-template.yaml > /app/config.yaml

/app/sog-backend
