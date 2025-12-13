package kafka

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

// Consumer handles consuming events from Kafka for analytics
type Consumer struct {
	reader  *kafka.Reader
	enabled bool
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers, topic, groupID string, enabled bool) *Consumer {
	if !enabled {
		log.Info().Msg("Kafka consumer disabled")
		return &Consumer{enabled: false}
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokers},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Info().Str("brokers", brokers).Str("topic", topic).Str("groupId", groupID).Msg("Kafka consumer created")
	return &Consumer{reader: reader, enabled: true}
}

// Close closes the Kafka reader
func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

// Start starts consuming messages and processes them with the handler
func (c *Consumer) Start(ctx context.Context, handler func(GameEvent)) {
	if !c.enabled {
		return
	}

	log.Info().Msg("Starting Kafka consumer")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Kafka consumer stopping")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Error().Err(err).Msg("Error reading message")
				continue
			}

			var event GameEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Error().Err(err).Msg("Error unmarshaling event")
				continue
			}

			handler(event)
		}
	}
}

// ProcessEvent is a handler that processes analytics events from Kafka
// In production, this would aggregate metrics and store to database
func ProcessEvent(event GameEvent) {
	switch event.Type {
	case EventGameStarted:
		log.Info().
			Str("type", string(event.Type)).
			Str("gameId", event.Data["gameId"].(string)).
			Str("player1", event.Data["player1"].(string)).
			Str("player2", event.Data["player2"].(string)).
			Bool("isBot", event.Data["isBotGame"].(bool)).
			Msg("Analytics: Game started")

	case EventGameMove:
		log.Info().
			Str("type", string(event.Type)).
			Str("gameId", event.Data["gameId"].(string)).
			Str("player", event.Data["player"].(string)).
			Interface("column", event.Data["column"]).
			Interface("moveNumber", event.Data["moveNumber"]).
			Msg("Analytics: Move made")

	case EventGameEnded:
		log.Info().
			Str("type", string(event.Type)).
			Str("gameId", event.Data["gameId"].(string)).
			Str("winner", event.Data["winner"].(string)).
			Str("result", event.Data["result"].(string)).
			Interface("totalMoves", event.Data["totalMoves"]).
			Msg("Analytics: Game ended")

	case EventPlayerConnected, EventPlayerDisconnected:
		log.Info().
			Str("type", string(event.Type)).
			Str("username", event.Data["username"].(string)).
			Msg("Analytics: Player event")

	case EventMatchmakingTimeout:
		log.Info().
			Str("type", string(event.Type)).
			Str("username", event.Data["username"].(string)).
			Interface("waitDuration", event.Data["waitDuration"]).
			Msg("Analytics: Matchmaking timeout")

	default:
		log.Debug().
			Str("type", string(event.Type)).
			Interface("data", event.Data).
			Msg("Analytics: Unknown event")
	}
}
