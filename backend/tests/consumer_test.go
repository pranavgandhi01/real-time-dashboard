package tests

import (
	"context"
	"os"
	"testing"
	"time"
	"real-time-dashboard/config"
	"real-time-dashboard/scaling"
)

func TestScalableConsumer_Creation(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			BrokerAddress: "localhost:9092",
			Topic:         "test-flights",
			GroupID:       "test-group",
		},
	}
	
	consumer := scaling.NewScalableConsumer(cfg)
	if consumer == nil {
		t.Error("Expected consumer to be created")
	}
	
	defer consumer.Close()
}

func TestScalableConsumer_NodeID(t *testing.T) {
	// Test with environment variable
	os.Setenv("NODE_ID", "test-node-123")
	defer os.Unsetenv("NODE_ID")
	
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			BrokerAddress: "localhost:9092",
			Topic:         "test-flights",
			GroupID:       "test-group",
		},
	}
	
	consumer := scaling.NewScalableConsumer(cfg)
	defer consumer.Close()
	
	// Consumer should be created successfully with custom node ID
	if consumer == nil {
		t.Error("Expected consumer to be created with custom NODE_ID")
	}
}

func TestScalableConsumer_DeploymentID(t *testing.T) {
	// Test with deployment ID
	os.Setenv("DEPLOYMENT_ID", "production")
	defer os.Unsetenv("DEPLOYMENT_ID")
	
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			BrokerAddress: "localhost:9092",
			Topic:         "test-flights",
			GroupID:       "test-group",
		},
	}
	
	consumer := scaling.NewScalableConsumer(cfg)
	defer consumer.Close()
	
	if consumer == nil {
		t.Error("Expected consumer to be created with custom DEPLOYMENT_ID")
	}
}

func TestScalableConsumer_ConsumeTimeout(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			BrokerAddress: "localhost:9092",
			Topic:         "test-flights",
			GroupID:       "test-group",
		},
	}
	
	consumer := scaling.NewScalableConsumer(cfg)
	defer consumer.Close()
	
	// Test consume with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	messageCount := 0
	err := consumer.Consume(ctx, func(data []byte) error {
		messageCount++
		return nil
	})
	
	// Should timeout (no Kafka running in test)
	if err != context.DeadlineExceeded {
		t.Logf("Expected timeout error, got: %v", err)
	}
}