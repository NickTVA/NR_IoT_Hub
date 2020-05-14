# NR_IoT_Hub

go run main.go NR_KEY

Endpoint for Metrics:

GET http://localhost:4590/metric?id=id123&name=metric_name&type=gauge&value=10

Endpoint for Logs:

POST http://localhost:4590/log?id=id123


Log data in POST Body