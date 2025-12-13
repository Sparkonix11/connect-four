import { Lobby } from './Lobby';
import { Board } from './Board';
import { GameStatus } from './GameStatus';
import { GameState, GameOverState, ExistingSession } from '../types';

interface PlayTabProps {
    username: string | null;
    connected: boolean;
    gameState: GameState | null;
    queuePosition: number | null;
    gameOver: GameOverState | null;
    error: string | null;
    existingSession: ExistingSession | null;
    onJoin: (username: string) => void;
    makeMove: (col: number) => void;
    onPlayAgain: () => void;
    onLeave: () => void;
    onResumeSession: () => void;
    onAbandonSession: () => void;
}

export function PlayTab({
    username,
    connected,
    gameState,
    queuePosition,
    gameOver,
    error,
    existingSession,
    onJoin,
    makeMove,
    onPlayAgain,
    onLeave,
    onResumeSession,
    onAbandonSession
}: PlayTabProps) {
    // Show session confirmation if user has an existing session
    if (existingSession) {
        return (
            <div className="flex-1 flex items-center justify-center p-4">
                <div className="text-center animate-fade-in max-w-md">
                    <h2 className="text-3xl font-bold mb-4 tracking-tight">Active Game Found</h2>
                    <p className="text-stone-500 mb-2 font-mono text-sm">
                        You have an ongoing game against:
                    </p>
                    <p className="text-xl font-bold mb-6">
                        {existingSession.isBot ? '🤖 Bot' : existingSession.opponent}
                    </p>

                    <div className="flex flex-col gap-3">
                        <button
                            onClick={onResumeSession}
                            className="w-full py-4 text-sm font-bold tracking-widest uppercase bg-black text-white
                           hover:bg-stone-800 transition-all active:scale-[0.98]"
                        >
                            Resume Game
                        </button>
                        <button
                            onClick={onAbandonSession}
                            className="w-full py-4 text-sm font-bold tracking-widest uppercase border-2 border-red-200 text-red-600
                           hover:bg-red-50 transition-all active:scale-[0.98]"
                        >
                            Forfeit & Start New
                        </button>
                    </div>

                    <p className="text-stone-400 text-xs mt-4 font-mono">
                        Abandoning will count as a loss
                    </p>
                </div>
            </div>
        );
    }

    // Show lobby if: no username, OR in queue (queuePosition is not null), OR not in game state
    const showLobby = !username || queuePosition !== null || (!gameState && !gameOver);

    if (showLobby) {
        return (
            <div className="flex-1 flex items-center justify-center p-4">
                <Lobby onJoin={onJoin} queuePosition={queuePosition} connected={connected || !username} />
            </div>
        );
    }

    // Show game
    return (
        <div className="flex-1 flex flex-col items-center justify-center gap-6 p-4 w-full max-w-4xl mx-auto animate-fade-in">
            <GameStatus
                gameState={gameState}
                gameOver={gameOver}
                error={error}
                onPlayAgain={onPlayAgain}
                onLeave={onLeave}
            />

            {gameState && !gameOver && (
                <Board
                    board={gameState.board}
                    onColumnClick={makeMove}
                    disabled={!gameState.yourTurn}
                    yourColor={gameState.yourColor}
                />
            )}

            {gameOver?.finalBoard && (
                <Board
                    board={gameOver.finalBoard}
                    onColumnClick={() => { }}
                    disabled={true}
                    yourColor={gameState?.yourColor || 1}
                />
            )}
        </div>
    );
}
