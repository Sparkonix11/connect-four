import { Lobby } from './Lobby';
import { Board } from './Board';
import { GameStatus } from './GameStatus';
import { GameState, GameOverState } from '../types';

interface PlayTabProps {
    username: string | null;
    connected: boolean;
    gameState: GameState | null;
    queuePosition: number | null;
    gameOver: GameOverState | null;
    error: string | null;
    onJoin: (username: string) => void;
    makeMove: (col: number) => void;
    onPlayAgain: () => void;
    onLeave: () => void;
}

export function PlayTab({
    username,
    connected,
    gameState,
    queuePosition,
    gameOver,
    error,
    onJoin,
    makeMove,
    onPlayAgain,
    onLeave
}: PlayTabProps) {
    // Show lobby if no username or not in game
    if (!username || (!gameState && queuePosition === null && !gameOver)) {
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
