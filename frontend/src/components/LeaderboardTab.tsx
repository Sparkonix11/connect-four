import { Leaderboard } from './Leaderboard';

export function LeaderboardTab() {
    return (
        <div className="flex-1 w-full max-w-4xl mx-auto p-4 animate-fade-in flex flex-col items-center">
            <h1 className="text-4xl font-black text-black mb-12 tracking-tighter">
                GLOBAL RANKING
            </h1>
            <div className="w-full">
                <Leaderboard />
            </div>
        </div>
    );
}
