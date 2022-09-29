# NR_IoT_Hub

go run main.go NR_INSIGHTS_KEY RPM_ID

Endpoint for Metrics:

GET http://localhost:4590/metric?id=id123&name=metric_name&type=gauge&value=10

Endpoint for Logs:

POST http://localhost:4590/log?id=id123

Endpoint for Ping:
GET http://localhost:4590/ping?type=ping&id=123

Log data in POST Body

##Running in Docker

Forward port 4590 and set INSIGHTS_KEY and ACCOUNT_ID as environment variable

docker build -t <image_tag> . && docker run -p 0.0.0.0:4590:4590 --env INSIGHTS_KEY=<insights_key> ACCOUNT_ID=1234 <image_tag> 