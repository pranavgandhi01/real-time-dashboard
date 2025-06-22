module mock-data-service

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/confluentinc/confluent-kafka-go v1.9.2
)

replace mock-data-service/pkg => ../../pkg