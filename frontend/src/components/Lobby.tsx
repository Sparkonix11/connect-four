import { useState, FormEvent } from 'react';

interface LobbyProps {
    onJoin: (username: string) => void;
    queuePosition: number | null;
    connected: boolean;
}

export function Lobby({ onJoin, queuePosition, connected }: LobbyProps) {
    const [username, setUsername] = useState('');
    const [isJoining, setIsJoining] = useState(false);

    const handleSubmit = (e: FormEvent) => {
        e.preventDefault();
        if (username.trim()) {
            setIsJoining(true);
            onJoin(username.trim());
        }
    };

    if (queuePosition !== null) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-800">
                <div className="bg-white/5 backdrop-blur-lg border border-white/10 rounded-3xl p-12 text-center animate-fade-in">
                    <div className="w-12 h-12 border-4 border-white/10 border-t-indigo-500 rounded-full animate-spin mx-auto" />
                    <h2 className="text-white text-2xl font-bold mt-6 mb-2">Finding Opponent...</h2>
                    <p className="text-slate-400">Queue position: #{queuePosition}</p>
                    <p className="text-slate-500 text-sm mt-4">A bot will join if no player is found within 10 seconds</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-800">
            <div className="bg-white/5 backdrop-blur-lg border border-white/10 rounded-3xl p-12 text-center max-w-md w-full">
                <h1 className="text-4xl font-bold mb-2 bg-gradient-to-r from-indigo-500 to-purple-500 bg-clip-text text-transparent">
                    🎮 Connect Four
                </h1>
                <p className="text-slate-400 mb-8">Challenge players or play against our AI</p>

                <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                    <input
                        type="text"
                        placeholder="Enter your username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        maxLength={50}
                        autoFocus
                        className="px-6 py-4 text-lg border-2 border-white/10 rounded-xl bg-white/5 text-white
                       outline-none focus:border-indigo-500 transition-colors placeholder:text-slate-500"
                    />
                    <button
                        type="submit"
                        disabled={!connected || !username.trim() || isJoining}
                        className="px-8 py-4 text-lg font-semibold rounded-xl bg-gradient-to-r from-indigo-500 to-purple-600
                       text-white hover:-translate-y-0.5 hover:shadow-xl hover:shadow-indigo-500/40
                       transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed
                       disabled:hover:translate-y-0 disabled:hover:shadow-none"
                    >
                        {!connected ? 'Connecting...' : isJoining ? 'Joining...' : 'Play Now'}
                    </button>
                </form>

                {!connected && (
                    <p className="text-red-400 mt-4">Connecting to server...</p>
                )}
            </div>
        </div>
    );
}
