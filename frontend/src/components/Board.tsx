interface BoardProps {
    board: number[][];
    onColumnClick: (column: number) => void;
    disabled: boolean;
    yourColor: 1 | 2;
}

export function Board({ board, onColumnClick, disabled }: BoardProps) {
    return (
        <div className="flex flex-col items-center gap-2">
            {/* Column indicators */}
            <div className="flex gap-1">
                {[0, 1, 2, 3, 4, 5, 6].map(col => (
                    <button
                        key={col}
                        className="w-14 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg text-white font-bold
                       hover:scale-110 hover:shadow-lg hover:shadow-indigo-500/40 transition-all duration-200
                       disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
                        onClick={() => onColumnClick(col)}
                        disabled={disabled}
                    >
                        ▼
                    </button>
                ))}
            </div>

            {/* Game grid */}
            <div className="bg-gradient-to-br from-blue-800 to-blue-500 p-3 rounded-2xl shadow-2xl shadow-blue-900/50">
                {board.map((row, rowIdx) => (
                    <div key={rowIdx} className="flex gap-1 mb-1 last:mb-0">
                        {row.map((cell, colIdx) => (
                            <div
                                key={colIdx}
                                className="w-14 h-14 bg-white/10 rounded-full flex items-center justify-center cursor-pointer
                           hover:bg-white/20 transition-all duration-200"
                                onClick={() => !disabled && onColumnClick(colIdx)}
                            >
                                <div
                                    className={`w-12 h-12 rounded-full transition-all duration-300
                    ${cell === 1 ? 'bg-gradient-to-br from-red-500 to-red-700 shadow-lg shadow-red-500/40' : ''}
                    ${cell === 2 ? 'bg-gradient-to-br from-yellow-400 to-amber-500 shadow-lg shadow-yellow-500/40' : ''}
                  `}
                                />
                            </div>
                        ))}
                    </div>
                ))}
            </div>
        </div>
    );
}
