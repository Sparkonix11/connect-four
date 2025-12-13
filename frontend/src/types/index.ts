// Re-export generated types for easier imports
// These types are auto-generated from api/openapi.yaml
// Run: npm run generate

import type { components } from './api.gen';

// Schema types
export type LeaderboardEntry = components['schemas']['LeaderboardEntry'];
export type Player = components['schemas']['Player'];
export type GameRecord = components['schemas']['GameRecord'];
export type Board = components['schemas']['Board'];

// WebSocket types
export type WSMessageType = components['schemas']['WSMessageType'];

// Generic WSMessage that accepts any payload (for incoming messages)
export interface WSMessage {
    type: WSMessageType;
    payload: unknown;
    timestamp: string;
}

// Client -> Server payloads
export type JoinQueuePayload = components['schemas']['JoinQueuePayload'];
export type MakeMovePayload = components['schemas']['MakeMovePayload'];
export type ReconnectPayload = components['schemas']['ReconnectPayload'];

// Server -> Client payloads
export type QueueJoinedPayload = components['schemas']['QueueJoinedPayload'];
export type GameStartedPayload = components['schemas']['GameStartedPayload'];
export type MoveMadePayload = components['schemas']['MoveMadePayload'];
export type InvalidMovePayload = components['schemas']['InvalidMovePayload'];
export type GameOverPayload = components['schemas']['GameOverPayload'];
export type OpponentDisconnectedPayload = components['schemas']['OpponentDisconnectedPayload'];
export type GameForfeitedPayload = components['schemas']['GameForfeitedPayload'];
export type ErrorPayload = components['schemas']['ErrorPayload'];
export type GameStatePayload = components['schemas']['GameStatePayload'];

// Frontend state types
export type GameState = components['schemas']['GameState'];
export type GameOverState = components['schemas']['GameOverState'];
