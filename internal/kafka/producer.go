package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// EventType defines the types of events we publish
type EventType string

const (
	EventGameStarted        EventType = "game.started"
	EventGameMove           EventType = "game.move"
	EventGameEnded          EventType = "game.ended"
	EventPlayerConnected    EventType = "player.connected"
	EventPlayerDisconnected EventType = "player.disconnected"
	EventMatchmakingTimeout EventType = "matchmaking.timeout"
)

// GameEvent represents an event to be published
type GameEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Producer handles publishing events to Kafka
type Producer struct {
	writer  *kafka.Writer
	enabled bool
}

// NewProducer creates a new Kafka producer with optional SASL authentication
func NewProducer(brokers, topic string, enabled bool, username, password string) *Producer {
	if !enabled {
		log.Info().Msg("Kafka producer disabled")
		return &Producer{enabled: false}
	}

	var writer *kafka.Writer

	// If username and password provided, use SASL/PLAIN with TLS
	if username != "" && password != "" {
		mechanism := plain.Mechanism{
			Username: username,
			Password: password,
		}

		dialer := &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
			TLS:           &tls.Config{MinVersion: tls.VersionTLS12},
		}

		transport := &kafka.Transport{
			Dial: dialer.DialFunc,
		}

		writer = &kafka.Writer{
			Addr:         kafka.TCP(brokers),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
			Transport:    transport,
		}

		log.Info().Str("brokers", brokers).Str("topic", topic).Msg("Kafka producer created with SASL/PLAIN authentication")
	} else {
		// For local setup without credentials, don't set Transport field
		writer = &kafka.Writer{
			Addr:         kafka.TCP(brokers),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		}

		log.Info().Str("brokers", brokers).Str("topic", topic).Msg("Kafka producer created (no authentication)")
	}

	return &Producer{writer: writer, enabled: true}
}

// Close closes the Kafka writer
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

// Publish publishes an event to Kafka
func (p *Producer) Publish(ctx context.Context, event GameEvent) error {
	if !p.enabled {
		return nil
	}

	event.ID = uuid.New().String()
	event.Timestamp = time.Now()

	data, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal event")
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
	})

	if err != nil {
		log.Error().Err(err).Str("type", string(event.Type)).Msg("Failed to publish event")
		return err
	}

	log.Debug().Str("type", string(event.Type)).Msg("Event published")
	return nil
}

// PublishGameStarted publishes a game started event
func (p *Producer) PublishGameStarted(ctx context.Context, gameID uuid.UUID, player1, player2 string, isBot bool) {
	p.Publish(ctx, GameEvent{
		Type: EventGameStarted,
		Data: map[string]interface{}{
			"gameId":    gameID.String(),
			"player1":   player1,
			"player2":   player2,
			"isBotGame": isBot,
		},
	})
}

// PublishGameMove publishes a move event
func (p *Producer) PublishGameMove(ctx context.Context, gameID uuid.UUID, player string, column, moveNum int) {
	p.Publish(ctx, GameEvent{
		Type: EventGameMove,
		Data: map[string]interface{}{
			"gameId":     gameID.String(),
			"player":     player,
			"column":     column,
			"moveNumber": moveNum,
		},
	})
}

// PublishGameEnded publishes a game ended event
func (p *Producer) PublishGameEnded(ctx context.Context, gameID uuid.UUID, winner, result string, duration, totalMoves int) {
	p.Publish(ctx, GameEvent{
		Type: EventGameEnded,
		Data: map[string]interface{}{
			"gameId":     gameID.String(),
			"winner":     winner,
			"result":     result,
			"duration":   duration,
			"totalMoves": totalMoves,
		},
	})
}

// PublishPlayerConnected publishes a player connected event
func (p *Producer) PublishPlayerConnected(ctx context.Context, username string) {
	p.Publish(ctx, GameEvent{
		Type: EventPlayerConnected,
		Data: map[string]interface{}{
			"username": username,
		},
	})
}

// PublishPlayerDisconnected publishes a player disconnected event
func (p *Producer) PublishPlayerDisconnected(ctx context.Context, username string, gameID *uuid.UUID) {
	data := map[string]interface{}{
		"username": username,
	}
	if gameID != nil {
		data["gameId"] = gameID.String()
	}
	p.Publish(ctx, GameEvent{
		Type: EventPlayerDisconnected,
		Data: data,
	})
}

// PublishMatchmakingTimeout publishes a matchmaking timeout event
func (p *Producer) PublishMatchmakingTimeout(ctx context.Context, username string, waitDuration time.Duration) {
	p.Publish(ctx, GameEvent{
		Type: EventMatchmakingTimeout,
		Data: map[string]interface{}{
			"username":     username,
			"waitDuration": waitDuration.Seconds(),
		},
	})
}
