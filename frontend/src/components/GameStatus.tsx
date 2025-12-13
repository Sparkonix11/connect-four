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
            <div className="text-center animate-fade-in py-12">
                <div className="text-6xl font-black mb-6 tracking-tighter leading-none">
                    {gameOver.result === 'win' && 'VICTORY'}
                    {gameOver.result === 'loss' && 'DEFEAT'}
                    {gameOver.result === 'draw' && 'DRAW'}
                    {gameOver.result === 'forfeit' && 'OPPONENT FORFEIT'}
                </div>
                <div className="flex gap-4 justify-center">
                    <button
                        onClick={onPlayAgain}
                        className="px-8 py-3 text-sm font-bold tracking-widest uppercase bg-black text-white
                       hover:bg-stone-800 transition-all active:scale-[0.98]"
                    >
                        Play Again
                    </button>
                    <button
                        onClick={onLeave}
                        className="px-8 py-3 text-sm font-bold tracking-widest uppercase border-2 border-zinc-200 text-black
                       hover:bg-zinc-100 transition-all active:scale-[0.98]"
                    >
                        Leave
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="w-full max-w-2xl py-6 flex items-center justify-between border-b-2 border-stone-200">
            <div className="flex flex-col">
                <span className="text-xs font-mono text-stone-400 uppercase tracking-widest mb-1">Opponent</span>
                <span className="text-xl font-bold tracking-tight">{gameState?.opponent || 'Unknown'}</span>
            </div>

            <div className="flex flex-col items-center">
                <span className={`text-sm font-bold tracking-widest uppercase ${gameState?.yourTurn ? 'text-black' : 'text-stone-300'}`}>
                    {gameState?.yourTurn ? 'Your Turn' : "Waiting..."}
                </span>
            </div>

            <div className="flex flex-col items-end">
                <span className="text-xs font-mono text-stone-400 uppercase tracking-widest mb-1">You</span>
                <div className="flex items-center gap-2">
                    <span className="text-sm font-bold">Player {gameState?.yourColor}</span>
                    <div
                        className={`w-3 h-3 rounded-full border border-black
            ${gameState?.yourColor === 1 ? 'bg-red-500' : ''}
            ${gameState?.yourColor === 2 ? 'bg-yellow-400' : ''}
          `}
                    />
                </div>
            </div>

            {error && <div className="absolute top-full left-0 right-0 mt-2 text-red-500 text-xs font-mono text-center">{error}</div>}
        </div>
    );
}
