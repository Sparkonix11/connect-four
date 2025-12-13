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
            <div className="flex gap-2 mb-2">
                {[0, 1, 2, 3, 4, 5, 6].map(col => (
                    <button
                        key={col}
                        className="w-16 h-8 flex items-end justify-center pb-1 text-stone-300 hover:text-black transition-colors disabled:opacity-0"
                        onClick={() => onColumnClick(col)}
                        disabled={disabled}
                    >
                        <span className="text-2xl leading-none">↓</span>
                    </button>
                ))}
            </div>

            {/* Game grid */}
            <div className="p-4 bg-stone-200 rounded-lg">
                <div className="flex flex-col gap-2">
                    {board.map((row, rowIdx) => (
                        <div key={rowIdx} className="flex gap-2">
                            {row.map((cell, colIdx) => (
                                <div
                                    key={colIdx}
                                    className="w-16 h-16 bg-white rounded-full flex items-center justify-center cursor-pointer relative"
                                    onClick={() => !disabled && onColumnClick(colIdx)}
                                >
                                    {cell !== 0 && (
                                        <div
                                            className={`w-14 h-14 rounded-full transition-all duration-300 animate-fade-in
                                            ${cell === 1 ? 'bg-red-500' : ''}
                                            ${cell === 2 ? 'bg-yellow-400' : ''}
                                        `}
                                        />
                                    )}
                                </div>
                            ))}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
