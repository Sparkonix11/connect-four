import { useState, useCallback } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { PlayTab } from './components/PlayTab';
import { LeaderboardTab } from './components/LeaderboardTab';

function App() {
    const [activeTab, setActiveTab] = useState<'play' | 'leaderboard'>('play');
    const [username, setUsername] = useState<string | null>(null);
    const { connected, gameState, queuePosition, error, gameOver, joinQueue, makeMove, leaveGame, resetGame } = useWebSocket(username);

    const handleJoin = useCallback((name: string) => {
        setUsername(name);
        // Wait a moment for connection, then join queue
        setTimeout(() => joinQueue(), 100);
    }, [joinQueue]);

    const handlePlayAgain = useCallback(() => {
        // Reset game state without sending leave message (game is already over)
        resetGame();
        setTimeout(() => joinQueue(), 100);
    }, [joinQueue, resetGame]);

    const handleLeave = useCallback(() => {
        leaveGame();
        setUsername(null);
    }, [leaveGame]);

    return (
        <div className="min-h-screen bg-[#F3F3F1] flex flex-col font-sans text-stone-900 selection:bg-stone-300">
            {/* Minimal Navigation */}
            <nav className="flex-none p-8 flex justify-between items-center max-w-7xl mx-auto w-full">
                <div className="text-xl font-bold tracking-tight">
                    CONNECT / FOUR
                </div>
                <div className="flex gap-8">
                    <button
                        onClick={() => setActiveTab('play')}
                        className={`text-sm tracking-widest uppercase transition-all duration-300 ${activeTab === 'play'
                            ? 'font-bold border-b-2 border-black'
                            : 'text-stone-500 hover:text-black'
                            }`}
                    >
                        Play
                    </button>
                    <button
                        onClick={() => setActiveTab('leaderboard')}
                        className={`text-sm tracking-widest uppercase transition-all duration-300 ${activeTab === 'leaderboard'
                            ? 'font-bold border-b-2 border-black'
                            : 'text-stone-500 hover:text-black'
                            }`}
                    >
                        Leaderboard
                    </button>
                </div>
                {/* Placeholder for symmetry or menu */}
                <div className="w-[100px] text-right text-xs text-stone-400 hidden md:block">
                    V 1.0
                </div>
            </nav>

            {/* Content Area */}
            <main className="flex-1 flex flex-col relative overflow-hidden max-w-7xl mx-auto w-full">
                {activeTab === 'play' ? (
                    <PlayTab
                        username={username}
                        connected={connected}
                        gameState={gameState}
                        queuePosition={queuePosition}
                        gameOver={gameOver}
                        error={error}
                        onJoin={handleJoin}
                        makeMove={makeMove}
                        onPlayAgain={handlePlayAgain}
                        onLeave={handleLeave}
                    />
                ) : (
                    <LeaderboardTab />
                )}
            </main>
        </div>
    );
}

export default App;
