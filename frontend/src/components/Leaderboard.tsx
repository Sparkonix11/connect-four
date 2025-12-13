import { useState, useEffect } from 'react';
import type { LeaderboardEntry } from '../types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export function Leaderboard() {
    const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetchLeaderboard();
    }, []);

    const fetchLeaderboard = async () => {
        try {
            const res = await fetch(`${API_URL}/api/leaderboard?limit=10`);
            if (!res.ok) throw new Error('Failed to fetch');
            const data = await res.json();
            setEntries(data || []);
        } catch (err) {
            setError('Unable to load leaderboard');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="bg-white/5 backdrop-blur-lg border border-white/10 rounded-2xl p-6 min-w-[300px]">
            <h2 className="text-white text-xl font-bold mb-4">🏆 Leaderboard</h2>

            {loading && (
                <p className="text-slate-400 text-center py-6">Loading...</p>
            )}

            {error && (
                <p className="text-red-400 text-center py-6">{error}</p>
            )}

            {!loading && !error && entries.length === 0 && (
                <p className="text-slate-400 text-center py-6">No games played yet. Be the first!</p>
            )}

            {!loading && !error && entries.length > 0 && (
                <table className="w-full">
                    <thead>
                        <tr>
                            <th className="text-left p-2 text-slate-400 font-medium text-sm border-b border-white/10">Rank</th>
                            <th className="text-left p-2 text-slate-400 font-medium text-sm border-b border-white/10">Player</th>
                            <th className="text-left p-2 text-slate-400 font-medium text-sm border-b border-white/10">Wins</th>
                            <th className="text-left p-2 text-slate-400 font-medium text-sm border-b border-white/10">Games</th>
                        </tr>
                    </thead>
                    <tbody>
                        {entries.map((entry) => (
                            <tr
                                key={entry.username}
                                className={`
                  ${entry.rank === 1 ? 'bg-yellow-500/10' : ''}
                  ${entry.rank === 2 ? 'bg-slate-400/10' : ''}
                  ${entry.rank === 3 ? 'bg-amber-700/10' : ''}
                `}
                            >
                                <td className="p-3 text-white text-xl">
                                    {entry.rank === 1 && '🥇'}
                                    {entry.rank === 2 && '🥈'}
                                    {entry.rank === 3 && '🥉'}
                                    {entry.rank > 3 && `#${entry.rank}`}
                                </td>
                                <td className="p-3 text-white font-semibold">{entry.username}</td>
                                <td className="p-3 text-white">{entry.wins}</td>
                                <td className="p-3 text-white">{entry.games}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
}
