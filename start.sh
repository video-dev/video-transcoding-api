#!/bin/bash

# Logfile location
LOG_FILE=/home/forge/logs/video-transcoding-api.log
mkdir -p /home/forge/logs/

# Change to ssm dir
cd /home/forge/ssm-cli
# Run pipenv
pipenv install
# Gets all video-transcoding-api secrets and stores them in secrets/video-transcoding-api-stage.env
pipenv run python spr get -e stage -r ECS -s video-transcoding-api --secure
# Export all the secrets as env variables
export $(cat secrets/video-transcoding-api-stage.env | xargs)

# Change to project dir
cd /home/forge/video-transcoding-api.sportsrecruits.com
# Include forge secrets
export $(cat .env | xargs)
echo $HTTP_PORT

# Run the service
echo "Starting video-transcoding-api service and using AWS access Key ID: $(echo $MEDIACONVERT_AWS_ACCESS_KEY_ID)" >> $LOG_FILE
./video-transcoding-api