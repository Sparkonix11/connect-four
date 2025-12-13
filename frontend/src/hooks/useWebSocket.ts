import { useState, useEffect, useCallback, useRef } from 'react';
import type {
    WSMessage,
    GameState,
    GameOverState,
    QueueJoinedPayload,
    GameStartedPayload,
    MoveMadePayload,
    InvalidMovePayload,
    GameOverPayload,
    OpponentDisconnectedPayload,
    GameForfeitedPayload,
    ErrorPayload,
    GameStatePayload,
} from '../types';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

interface UseWebSocketReturn {
    connected: boolean;
    gameState: GameState | null;
    queuePosition: number | null;
    error: string | null;
    gameOver: GameOverState | null;
    joinQueue: () => void;
    makeMove: (column: number) => void;
    leaveGame: () => void;
}

export function useWebSocket(username: string | null): UseWebSocketReturn {
    const [connected, setConnected] = useState(false);
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [queuePosition, setQueuePosition] = useState<number | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [gameOver, setGameOver] = useState<GameOverState | null>(null);
    const wsRef = useRef<WebSocket | null>(null);

    const connect = useCallback(() => {
        if (!username) return;

        const ws = new WebSocket(`${WS_URL}?username=${encodeURIComponent(username)}`);
        wsRef.current = ws;

        ws.onopen = () => {
            console.log('WebSocket connected');
            setConnected(true);
            setError(null);
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected');
            setConnected(false);
        };

        ws.onerror = () => {
            console.error('WebSocket error');
            setError('Connection error');
        };

        ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data) as WSMessage;
                handleMessage(message);
            } catch (e) {
                console.error('Failed to parse message:', e);
            }
        };

        return () => {
            ws.close();
        };
    }, [username]);

    const handleMessage = (message: WSMessage) => {
        console.log('Received:', message.type, message.payload);

        switch (message.type) {
            case 'queue_joined': {
                const payload = message.payload as QueueJoinedPayload;
                setQueuePosition(payload.position);
                break;
            }

            case 'game_started': {
                const payload = message.payload as GameStartedPayload;
                setGameState({
                    gameId: payload.gameId,
                    opponent: payload.opponent,
                    yourTurn: payload.yourTurn,
                    yourColor: payload.yourColor,
                    board: Array(6).fill(null).map(() => Array(7).fill(0)),
                });
                setQueuePosition(null);
                setGameOver(null);
                break;
            }

            case 'move_made': {
                const payload = message.payload as MoveMadePayload;
                setGameState(prev => prev ? {
                    ...prev,
                    board: payload.board,
                    yourTurn: prev.yourColor !== payload.player,
                } : null);
                break;
            }

            case 'invalid_move': {
                const payload = message.payload as InvalidMovePayload;
                setError(payload.reason);
                setTimeout(() => setError(null), 3000);
                break;
            }

            case 'game_over': {
                const payload = message.payload as GameOverPayload;
                setGameOver({
                    winner: payload.winner,
                    result: payload.result,
                    finalBoard: payload.finalBoard,
                });
                break;
            }

            case 'game_state': {
                const payload = message.payload as GameStatePayload;
                setGameState({
                    gameId: payload.gameId,
                    opponent: payload.opponent,
                    yourTurn: payload.yourTurn,
                    yourColor: payload.yourColor,
                    board: payload.board,
                });
                break;
            }

            case 'opponent_disconnected': {
                const payload = message.payload as OpponentDisconnectedPayload;
                setError(`Opponent disconnected. Waiting ${payload.timeout}s for reconnect...`);
                break;
            }

            case 'opponent_reconnected':
                setError(null);
                break;

            case 'game_forfeited': {
                const payload = message.payload as GameForfeitedPayload;
                setGameOver({
                    winner: payload.winner,
                    result: 'forfeit',
                });
                break;
            }

            case 'error': {
                const payload = message.payload as ErrorPayload;
                setError(payload.message);
                break;
            }
        }
    };

    const sendMessage = useCallback((type: string, payload: unknown) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                type,
                payload,
                timestamp: new Date().toISOString(),
            }));
        }
    }, []);

    const joinQueue = useCallback(() => {
        sendMessage('join_queue', { username });
    }, [sendMessage, username]);

    const makeMove = useCallback((column: number) => {
        sendMessage('make_move', { column });
    }, [sendMessage]);

    const leaveGame = useCallback(() => {
        sendMessage('leave_game', {});
        setGameState(null);
        setGameOver(null);
    }, [sendMessage]);

    useEffect(() => {
        const cleanup = connect();
        return cleanup;
    }, [connect]);

    return {
        connected,
        gameState,
        queuePosition,
        error,
        gameOver,
        joinQueue,
        makeMove,
        leaveGame,
    };
}
