#!/bin/bash

echo "-> deploying functions..."

until doctl serverless deploy . --env data.env
do
    echo "-> failed to deploy functions, trying again..."
    sleep 1
done

echo "-> functions deployed successfully"
