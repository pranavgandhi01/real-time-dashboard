package scaling

import (
	"context"
	"fmt"
	"os"
	"time"
	"real-time-dashboard/config"
	flightlog "real-time-dashboard/log"
	"github.com/segmentio/kafka-go"
)

type ScalableConsumer struct {
	reader   *kafka.Reader
	nodeID   string
	groupID  string
}

func NewScalableConsumer(cfg *config.Config) *ScalableConsumer {
	nodeID := getNodeID()
	groupID := fmt.Sprintf("%s-%s", cfg.Kafka.GroupID, getDeploymentID())
	
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Kafka.BrokerAddress},
		Topic:   cfg.Kafka.Topic,
		GroupID: groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait: 1 * time.Second,
		StartOffset: kafka.LastOffset,
		Partition: -1, // Auto-assign partitions
	})
	
	flightlog.LogInfo("Scalable consumer created: NodeID=%s, GroupID=%s", nodeID, groupID)
	
	return &ScalableConsumer{
		reader:  reader,
		nodeID:  nodeID,
		groupID: groupID,
	}
}

func (sc *ScalableConsumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := sc.reader.ReadMessage(ctx)
			if err != nil {
				flightlog.LogError("Consumer read error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}
			
			if err := handler(m.Value); err != nil {
				flightlog.LogError("Message handler error: %v", err)
			}
		}
	}
}

func (sc *ScalableConsumer) Close() error {
	return sc.reader.Close()
}

func getNodeID() string {
	if nodeID := os.Getenv("NODE_ID"); nodeID != "" {
		return nodeID
	}
	hostname, _ := os.Hostname()
	return hostname
}

func getDeploymentID() string {
	if deployID := os.Getenv("DEPLOYMENT_ID"); deployID != "" {
		return deployID
	}
	return "default"
}