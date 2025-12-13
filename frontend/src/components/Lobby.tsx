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
                    <h2 className="text-4xl font-bold mb-4 tracking-tighter">WAITING</h2>
                    <div className="flex flex-col gap-2 text-stone-500 font-mono text-sm">
                        <p>QUEUE POSITION: {queuePosition.toString().padStart(2, '0')}</p>
                        <p>EST. TIME: 00:10</p>
                    </div>
                    <div className="mt-8 w-16 h-0.5 bg-black animate-pulse mx-auto opacity-20" />
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
