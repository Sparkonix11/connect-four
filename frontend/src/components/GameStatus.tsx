import type { GameState, GameOverState } from '../types';

interface GameStatusProps {
    gameState: GameState | null;
    gameOver: GameOverState | null;
    error: string | null;
    onPlayAgain: () => void;
    onLeave: () => void;
}

export function GameStatus({ gameState, gameOver, error, onPlayAgain, onLeave }: GameStatusProps) {
    if (gameOver) {
        return (
            <div className="bg-white/5 backdrop-blur-lg border border-white/10 rounded-2xl p-8 text-center">
                <div className="text-3xl mb-6">
                    {gameOver.result === 'win' && <span className="text-green-400">🎉 You Won!</span>}
                    {gameOver.result === 'loss' && <span className="text-red-400">😔 You Lost</span>}
                    {gameOver.result === 'draw' && <span className="text-yellow-400">🤝 Draw!</span>}
                    {gameOver.result === 'forfeit' && <span className="text-green-400">🏆 {gameOver.winner} wins by forfeit!</span>}
                </div>
                <div className="flex gap-4 justify-center">
                    <button
                        onClick={onPlayAgain}
                        className="px-6 py-3 font-semibold rounded-xl bg-gradient-to-r from-indigo-500 to-purple-600
                       text-white hover:-translate-y-0.5 hover:shadow-lg hover:shadow-indigo-500/40
                       transition-all duration-200"
                    >
                        Play Again
                    </button>
                    <button
                        onClick={onLeave}
                        className="px-6 py-3 font-semibold rounded-xl bg-white/10 text-white border border-white/20
                       hover:bg-white/15 transition-all duration-200"
                    >
                        Leave
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white/5 backdrop-blur-lg border border-white/10 rounded-2xl px-8 py-5
                    flex items-center justify-between gap-8 min-w-[500px]">
            <div className="flex flex-col gap-1">
                <span className="text-slate-400 text-sm">vs {gameState?.opponent || 'Unknown'}</span>
                <span className={`text-xl font-semibold ${gameState?.yourTurn ? 'text-green-400' : 'text-white'}`}>
                    {gameState?.yourTurn ? '🟢 Your Turn' : "⏳ Opponent's Turn"}
                </span>
            </div>

            <div className="flex items-center gap-3 text-slate-400">
                <span>You are:</span>
                <div
                    className={`w-8 h-8 rounded-full shadow-lg
            ${gameState?.yourColor === 1 ? 'bg-gradient-to-br from-red-500 to-red-700 shadow-red-500/40' : ''}
            ${gameState?.yourColor === 2 ? 'bg-gradient-to-br from-yellow-400 to-amber-500 shadow-yellow-500/40' : ''}
          `}
                />
            </div>

            {error && <div className="text-red-400 text-sm">{error}</div>}
        </div>
    );
}
