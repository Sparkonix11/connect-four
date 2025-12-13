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
            <div className="w-full flex flex-col items-center justify-center min-h-[50vh]">
                <div className="text-center animate-fade-in">
                    {/* Animated matching icon */}
                    <div className="mb-8 relative">
                        <div className="flex gap-3 justify-center items-center">
                            {/* Left disc */}
                            <div className="w-16 h-16 rounded-full bg-gradient-to-br from-red-500 to-red-700 shadow-lg shadow-red-500/30 animate-pulse" />

                            {/* Connecting dots */}
                            <div className="flex gap-1">
                                <div className="w-2 h-2 rounded-full bg-stone-400 animate-bounce" style={{ animationDelay: '0ms' }} />
                                <div className="w-2 h-2 rounded-full bg-stone-400 animate-bounce" style={{ animationDelay: '150ms' }} />
                                <div className="w-2 h-2 rounded-full bg-stone-400 animate-bounce" style={{ animationDelay: '300ms' }} />
                            </div>

                            {/* Right disc */}
                            <div className="w-16 h-16 rounded-full bg-gradient-to-br from-yellow-400 to-amber-500 shadow-lg shadow-yellow-500/30 animate-pulse" style={{ animationDelay: '500ms' }} />
                        </div>
                    </div>

                    <h2 className="text-3xl font-bold mb-2 tracking-tight">Finding Opponent...</h2>
                    <p className="text-stone-500 font-mono text-sm mb-6">
                        Queue Position: {queuePosition === 0 ? 'Searching...' : `#${queuePosition}`}
                    </p>

                    {/* Loading bar */}
                    <div className="w-48 h-1 bg-stone-200 rounded-full overflow-hidden mx-auto">
                        <div className="h-full w-1/3 bg-gradient-to-r from-red-500 to-yellow-500 rounded-full animate-[loading_1.5s_ease-in-out_infinite]" />
                    </div>

                    <p className="text-stone-400 text-xs mt-4 font-mono">
                        {queuePosition === 0 ? 'Bot match in ~10s if no player found' : 'Waiting for players...'}
                    </p>
                </div>
            </div>
        );
    }

    return (
        <div className="w-full flex items-center justify-center min-h-[50vh]">
            <div className="max-w-md w-full p-8">
                <h1 className="text-6xl font-black mb-8 tracking-tighter leading-[0.85]">
                    READY<br />TO<br />PLAY?
                </h1>
                <p className="text-stone-500 mb-12 font-mono text-xs uppercase tracking-widest">
                    Enter your moniker to begin
                </p>

                <form onSubmit={handleSubmit} className="flex flex-col gap-8">
                    <input
                        type="text"
                        placeholder="USERNAME"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        maxLength={50}
                        autoFocus
                        className="w-full py-4 text-2xl font-bold bg-transparent border-b-2 border-stone-200 
                       text-black placeholder:text-stone-300 focus:border-black outline-none transition-colors rounded-none"
                    />
                    <button
                        type="submit"
                        disabled={!connected || !username.trim() || isJoining}
                        className="w-full py-4 text-sm font-bold tracking-widest uppercase bg-black text-white
                       hover:bg-stone-800 disabled:opacity-30 disabled:cursor-not-allowed
                       transition-all active:scale-[0.98]"
                    >
                        {!connected ? 'Connecting...' : isJoining ? 'Joining...' : 'Start Game'}
                    </button>
                </form>

                {!connected && (
                    <p className="text-red-500 mt-4 text-xs font-mono">CONNECTION LOST</p>
                )}
            </div>
        </div>
    );
}
