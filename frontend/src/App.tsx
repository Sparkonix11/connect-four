import { useState, useCallback } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { Lobby } from './components/Lobby';
import { Board } from './components/Board';
import { GameStatus } from './components/GameStatus';
import { Leaderboard } from './components/Leaderboard';

function App() {
    const [username, setUsername] = useState<string | null>(null);
    const { connected, gameState, queuePosition, error, gameOver, joinQueue, makeMove, leaveGame } = useWebSocket(username);

    const handleJoin = useCallback((name: string) => {
        setUsername(name);
        // Wait a moment for connection, then join queue
        setTimeout(() => joinQueue(), 100);
    }, [joinQueue]);

    const handlePlayAgain = useCallback(() => {
        joinQueue();
    }, [joinQueue]);

    const handleLeave = useCallback(() => {
        leaveGame();
        setUsername(null);
    }, [leaveGame]);

    // Show lobby if no username or not in game
    if (!username || (!gameState && queuePosition === null && !gameOver)) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex">
                <div className="flex-1 flex items-center justify-center">
                    <Lobby onJoin={handleJoin} queuePosition={queuePosition} connected={connected || !username} />
                </div>
                <div className="w-80 p-6 flex items-center">
                    <Leaderboard />
                </div>
            </div>
        );
    }

    // Show game
    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800 flex flex-col items-center justify-center gap-6 p-8">
            <GameStatus
                gameState={gameState}
                gameOver={gameOver}
                error={error}
                onPlayAgain={handlePlayAgain}
                onLeave={handleLeave}
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

            <div className="fixed bottom-4 right-4">
                <Leaderboard />
            </div>
        </div>
    );
}

export default App;
