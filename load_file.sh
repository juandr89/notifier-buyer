#!/bin/bash

# Define output YAML file
output_file="./config.yaml"
echo 
# Create the YAML file and write the content
cat <<EOL > $output_file
port: $PORT
api_key: $API_KEY
notification_sender: $SENDER
smtp:
  host: $SMTP_HOST
  port: $SMTP_PORT
  username: $SMTP_USERNAME
  password: $SMTP_PASSWORD
redis:
  host: $REDIS_HOST
  port: $REDIS_PORT
  password: $REDIS_PASSWORD
  db: 1
  tls_enable: $TLS_ENABLE
forecast_service:
  base_url: $FORECAST_URL
  api_key: $FORECAST_API_KEY
EOL

echo "YAML configuration file created at $output_file"

./my-go-app & 

wait $!

