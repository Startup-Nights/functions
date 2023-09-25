#!/bin/bash

echo -e "=> deploying functions...\n"

until doctl serverless deploy .
do
    echo -e "\n=> failed to deploy functions, trying again...\n"
    sleep 1
done

echo -e "\n=> functions deployed successfully"
