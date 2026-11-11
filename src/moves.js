let positionHistory = [];
let viewIndex = 0;

function isViewingLatest() {
    return positionHistory.length > 0 && viewIndex === positionHistory.length - 1;
}

function bitToIndex(bit) {
    let idx = 0;
    let v = BigInt(bit);
    if (v <= 0n) return -1;
    while (v > 1n) {
        v >>= 1n;
        idx++;
    }
    return idx;
}

function squareName(idx) {
    if (idx < 0 || idx > 63) return "";
    const file = "abcdefgh"[idx % 8];
    const rank = Math.floor(idx / 8) + 1;
    return file + rank;
}

function pieceLetterAt(b, bit) {
    const target = BigInt(bit);
    const map = [
        ["WhitePawns", ""], ["BlackPawns", ""],
        ["WhiteKnights", "N"], ["BlackKnights", "N"],
        ["WhiteBishops", "B"], ["BlackBishops", "B"],
        ["WhiteRooks", "R"], ["BlackRooks", "R"],
        ["WhiteQueens", "Q"], ["BlackQueens", "Q"],
        ["WhiteKing", "K"], ["BlackKing", "K"],
    ];
    for (const [key, letter] of map) {
        if (b[key] !== undefined && (BigInt(b[key]) & target) !== 0n) {
            return letter;
        }
    }
    return "";
}

function isSquareOccupied(b, bit) {
    const target = BigInt(bit);
    const keys = [
        "WhitePawns", "BlackPawns", "WhiteKnights", "BlackKnights",
        "WhiteBishops", "BlackBishops", "WhiteRooks", "BlackRooks",
        "WhiteQueens", "BlackQueens", "WhiteKing", "BlackKing"
    ];
    for (const key of keys) {
        if (b[key] !== undefined && (BigInt(b[key]) & target) !== 0n) {
            return true;
        }
    }
    return false;
}

function formatMove(prevBoard, fromBit, toBit) {
    const fromIdx = bitToIndex(fromBit);
    const toIdx = bitToIndex(toBit);
    const piece = pieceLetterAt(prevBoard, fromBit);
    const fromFile = fromIdx % 8;
    const toFile = toIdx % 8;
    const toRank = Math.floor(toIdx / 8);

    if (piece === "K" && Math.abs(fromFile - toFile) === 2) {
        return toFile === 6 ? "O-O" : "O-O-O";
    }

    const isEnPassant = piece === "" && fromFile !== toFile && !isSquareOccupied(prevBoard, toBit);
    const isCapture = isSquareOccupied(prevBoard, toBit) || isEnPassant;
    const isPromotion = piece === "" && (toRank === 7 || toRank === 0);

    let notation = "";
    if (piece === "") {
        if (isCapture) {
            notation += squareName(fromIdx)[0] + "x";
        }
        notation += squareName(toIdx);
        if (isPromotion) {
            notation += "=Q";
        }
    } else {
        notation += piece;
        if (isCapture) {
            notation += "x";
        }
        notation += squareName(toIdx);
    }
    return notation;
}

function initHistory(startBoard, openingMove, currentBoard) {
    positionHistory = [{
        board: startBoard,
        from: null,
        to: null,
        notation: "Start",
        color: null
    }];

    if (
        playerColor == BLACK &&
        openingMove &&
        BigInt(openingMove[0]) !== 0n &&
        BigInt(openingMove[1]) !== 0n
    ) {
        const fromBit = BigInt(openingMove[0]);
        const toBit = BigInt(openingMove[1]);
        positionHistory.push({
            board: currentBoard,
            from: fromBit,
            to: toBit,
            notation: formatMove(startBoard, fromBit, toBit),
            color: WHITE
        });
    }

    viewIndex = positionHistory.length - 1;
    goToMove(viewIndex);
}

function recordMove(boardState, fromBit, toBit, color, suffix = "") {
    const prevBoard = positionHistory[positionHistory.length - 1].board;
    positionHistory.push({
        board: boardState,
        from: fromBit,
        to: toBit,
        notation: formatMove(prevBoard, fromBit, toBit) + suffix,
        color
    });
}

function highlightLastMove(entry) {
    const squares = board.getElementsByClassName("square");
    for (const sq of squares) {
        sq.classList.remove("last-move");
    }
    if (!entry || entry.from === null || entry.to === null) {
        return;
    }
    for (const bit of [entry.from, entry.to]) {
        const idx = bitToIndex(bit);
        if (idx >= 0) {
            const file = idx % 8;
            const rank = Math.floor(idx / 8);
            const coords = boardToDomIndex(rank, file);
            if (squares[coords]) {
                squares[coords].classList.add("last-move");
            }
        }
    }
}

function renderHistory() {
    const moveListEl = document.getElementById("move-list");
    moveListEl.innerHTML = "";

    for (let idx = 1; idx < positionHistory.length; idx += 2) {
        const moveNum = Math.floor((idx - 1) / 2) + 1;
        const row = document.createElement("div");
        row.className = "move-row";

        const numEl = document.createElement("span");
        numEl.className = "move-num";
        numEl.innerText = moveNum + ".";
        row.appendChild(numEl);

        const whiteEntry = positionHistory[idx];
        const whiteEl = document.createElement("span");
        whiteEl.className = "move-item" + (viewIndex === idx ? " active" : "");
        whiteEl.innerText = whiteEntry.notation;
        whiteEl.onclick = () => goToMove(idx);
        row.appendChild(whiteEl);

        const blackEl = document.createElement("span");
        if (idx + 1 < positionHistory.length) {
            const blackIdx = idx + 1;
            const blackEntry = positionHistory[blackIdx];
            blackEl.className = "move-item" + (viewIndex === blackIdx ? " active" : "");
            blackEl.innerText = blackEntry.notation;
            blackEl.onclick = () => goToMove(blackIdx);
        }
        row.appendChild(blackEl);

        moveListEl.appendChild(row);
    }

    const activeEl = moveListEl.querySelector(".move-item.active");
    if (activeEl) {
        activeEl.scrollIntoView({ block: "nearest" });
    }
}

function updateViewState() {
    const isLatest = isViewingLatest();
    const statusEl = document.getElementById("history-status");
    statusEl.classList.remove("rewound", "ended");

    if (isWaitingForBot) {
        statusEl.innerText = "Bot thinking...";
    } else if (!isLatest) {
        statusEl.innerText = `Viewing move ${viewIndex} / ${positionHistory.length - 1} (rewind)`;
        statusEl.classList.add("rewound");
    } else if (gameOver) {
        statusEl.innerText = "Game over";
        statusEl.classList.add("ended");
    } else {
        statusEl.innerText = "Live";
    }

    if (canPlay()) {
        board.classList.remove("read-only");
        interact(".piece").draggable({ enabled: true });
    } else {
        board.classList.add("read-only");
        interact(".piece").draggable({ enabled: false });
    }

    document.getElementById("btn-first").disabled = viewIndex === 0;
    document.getElementById("btn-prev").disabled = viewIndex === 0;
    document.getElementById("btn-next").disabled = isLatest;
    document.getElementById("btn-last").disabled = isLatest;
}

function goToMove(idx) {
    if (positionHistory.length === 0) return;
    viewIndex = Math.max(0, Math.min(idx, positionHistory.length - 1));
    removeTarget();
    const entry = positionHistory[viewIndex];
    updatePieces(entry.board);
    highlightLastMove(entry);
    renderHistory();
    updateViewState();
}

function goToLatest() {
    goToMove(positionHistory.length - 1);
}

function refreshCurrentMove() {
    goToMove(viewIndex);
}

document.getElementById("btn-first").addEventListener("click", () => goToMove(0));
document.getElementById("btn-prev").addEventListener("click", () => goToMove(viewIndex - 1));
document.getElementById("btn-next").addEventListener("click", () => goToMove(viewIndex + 1));
document.getElementById("btn-last").addEventListener("click", () => goToLatest());

document.addEventListener("keydown", (e) => {
    if (e.key === "ArrowLeft") {
        e.preventDefault();
        goToMove(viewIndex - 1);
    } else if (e.key === "ArrowRight") {
        e.preventDefault();
        goToMove(viewIndex + 1);
    } else if (e.key === "Home") {
        e.preventDefault();
        goToMove(0);
    } else if (e.key === "End") {
        e.preventDefault();
        goToLatest();
    }
});
