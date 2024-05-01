#!/bin/bash

echo -e "=> deploying functions...\n"

args=""

# check if there is an environment file - this is only the case for local 
# deployments but it allows for using the same script to deploy the functions
if [ -f ./data.env ]; then
    args="--env data.env "
fi

until doctl serverless deploy . ${args}
do
    echo -e "\n=> failed to deploy functions, trying again...\n"
    sleep 10
done

echo -e "\n=> functions deployed successfully"
