import { useState, useEffect } from 'react';
import type { LeaderboardEntry } from '../types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export function Leaderboard() {
    const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetchLeaderboard();

        // Auto-refresh leaderboard every 30 seconds
        const interval = setInterval(() => {
            fetchLeaderboard();
        }, 30000);

        return () => clearInterval(interval);
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
        <div className="w-full">
            {loading && (
                <p className="text-stone-400 text-center font-mono py-6">LOADING DATA...</p>
            )}

            {error && (
                <p className="text-red-500 text-center font-mono py-6">{error}</p>
            )}

            {!loading && !error && entries.length === 0 && (
                <p className="text-stone-400 text-center font-mono py-6">NO RECORDS YET</p>
            )}

            {!loading && !error && entries.length > 0 && (
                <table className="w-full border-collapse">
                    <thead>
                        <tr className="border-b border-black">
                            <th className="text-left py-4 font-mono text-xs text-stone-500 uppercase tracking-widest">Rank</th>
                            <th className="text-left py-4 font-mono text-xs text-stone-500 uppercase tracking-widest">Player</th>
                            <th className="text-right py-4 font-mono text-xs text-stone-500 uppercase tracking-widest">Wins</th>
                            <th className="text-right py-4 font-mono text-xs text-stone-500 uppercase tracking-widest">Games</th>
                        </tr>
                    </thead>
                    <tbody>
                        {entries.map((entry) => (
                            <tr
                                key={entry.username}
                                className="border-b border-stone-200 hover:bg-stone-100 transition-colors"
                            >
                                <td className="py-4 font-mono text-sm">
                                    {entry.rank.toString().padStart(2, '0')}
                                </td>
                                <td className={`py-4 ${entry.rank <= 3 ? 'font-bold' : ''}`}>
                                    {entry.username}
                                </td>
                                <td className="py-4 text-right font-mono text-stone-600">{entry.wins}</td>
                                <td className="py-4 text-right font-mono text-stone-400">{entry.games}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
}
