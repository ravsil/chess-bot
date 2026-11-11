const BLACK = 1, WHITE = 0;
// BigInt cannot be converted to json normally
BigInt.prototype.toJSON = function () {
    return JSON.rawJSON(this.toString());
};

let pieceData = new Map();
let boardSetup;
const boardSize = 8;
let playerColor;
let moveList;
let board = document.getElementById("board");

let positionHistory = [];
let viewIndex = 0;
let gameOver = false;
let isWaitingForBot = false;

function canPlay() {
    return !gameOver && !isWaitingForBot && viewIndex === positionHistory.length - 1;
}

function domToBoardCoords(i, j) {
    if (playerColor == BLACK) {
        return { rank: i, file: boardSize - 1 - j };
    }
    return { rank: boardSize - 1 - i, file: j };
}

function boardToDomIndex(rank, file) {
    if (playerColor == BLACK) {
        const i = rank;
        const j = boardSize - 1 - file;
        return j + i * boardSize;
    }
    const i = boardSize - 1 - rank;
    return file + i * boardSize;
}

async function parseBigIntJson(response) {
    const text = await response.text();
    const safeText = text.replace(/:\s*(\d{15,})(?=[,\}\]])/g, ':"$1"');
    return JSON.parse(safeText);
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

async function initBoard() {
    const response = await fetch("getBoard");
    boardSetup = await parseBigIntJson(response);

    const response2nd = await fetch("getMoves");
    let r2json = await parseBigIntJson(response2nd);
    playerColor = r2json.playerColor;
    document.getElementById("name").innerText = (playerColor == WHITE) ? "You are playing White" : "You are playing Black";
    moveList = r2json.validMoves;

    board.innerHTML = "";
    for (let i = 0; i < boardSize; i++) {
        let row = document.createElement("div");
        row.className = "row";
        for (let j = 0; j < boardSize; j++) {
            const { rank, file } = domToBoardCoords(i, j);
            let square = document.createElement("div");
            let pieceName = getPiece(boardSetup, rank, file);
            if (pieceName != null) {
                let piece = document.createElement("div");
                piece.classList.add("piece");
                piece.classList.add(pieceName);
                square.appendChild(piece);
                pieceData.set(piece, { x: file, y: rank, dx: 0, dy: 0, color: pieceName[0] == "w" ? WHITE : BLACK });
            }

            square.className = "square";
            square.className += (i + j) % 2 === 0 ? " white" : " black";
            square.x = file;
            square.y = rank;
            row.appendChild(square);
        }
        board.appendChild(row);
    }

    const startBoard = r2json.startBoard || boardSetup;
    positionHistory = [{
        board: startBoard,
        from: null,
        to: null,
        notation: "Start",
        color: null
    }];
    if (
        playerColor == BLACK &&
        r2json.openingMove &&
        BigInt(r2json.openingMove[0]) !== 0n &&
        BigInt(r2json.openingMove[1]) !== 0n
    ) {
        const fromBit = BigInt(r2json.openingMove[0]);
        const toBit = BigInt(r2json.openingMove[1]);
        positionHistory.push({
            board: boardSetup,
            from: fromBit,
            to: toBit,
            notation: formatMove(startBoard, fromBit, toBit),
            color: WHITE
        });
    }
    viewIndex = positionHistory.length - 1;
    gameOver = false;
    isWaitingForBot = false;
    goToMove(viewIndex);
}

function getPiece(b, rank, file) {
    for (let key in b) {
        let bitboard = BigInt(b[key]);
        let pos = BigInt(file + rank * boardSize);
        if ((bitboard >> pos) & 1n) {
            let name = key.toLowerCase().replace(/s$/, "");
            return name;
        }
    }
    return null;
}

interact(".piece").draggable({
    listeners: {
        start(event) {
            removeTarget();
            let data = pieceData.get(event.target);
            if (!canPlay() || !data || data.color != playerColor) {
                return;
            }
            event.target.classList.add("dragged");
            addTarget(event);
        },
        move(event) {
            let data = pieceData.get(event.target);
            if (!canPlay() || !data || data.color != playerColor) {
                return;
            }
            data.dx += event.dx;
            data.dy += event.dy;
            event.target.style.transform =
                `translate(${data.dx}px, ${data.dy}px)`;
        },
        end(event) {
            removeTarget();
            event.target.classList.remove("dragged");
            if (!event.dropzone) {
                event.target.style.transform = "";
                let data = pieceData.get(event.target);
                if (data) {
                    data.dx = 0;
                    data.dy = 0;
                }
            }
        }
    }
});

interact(".target")
    .dropzone({
        accept: ".piece",
        overlap: 0.6,
        ondrop: function (event) {
            if (!canPlay()) {
                return;
            }
            if (event.target.firstChild) {
                event.target.removeChild(event.target.firstChild);
            }
            event.target.appendChild(event.relatedTarget);
            event.relatedTarget.style.transform = "";
            let piece = pieceData.get(event.relatedTarget);
            sendMoveToServer(piece, { X: event.target.x, Y: event.target.y });
            piece.dx = 0;
            piece.dy = 0;
            piece.x = event.target.x;
            piece.y = event.target.y;
        }
    });

interact(".piece").on("tap", addTarget);

function addTarget(event) {
    if (!canPlay()) {
        return;
    }
    let data = pieceData.get(event.target);
    if (!data || data.color != playerColor) {
        return;
    }
    let pos = BigInt(1) << BigInt(data.x + data.y * boardSize);
    let moves;
    for (let key in moveList) {
        if (key == pos) {
            moves = BigInt(moveList[key]);
            break;
        }
    }
    if (!moves) {
        return;
    }
    const squares = board.getElementsByClassName("square");
    for (let i = 0; i < boardSize; i++) {
        for (let j = 0; j < boardSize; j++) {
            const { rank, file } = domToBoardCoords(i, j);
            let targetPos = BigInt(1) << BigInt(file + rank * boardSize);
            if ((moves & targetPos) != 0n) {
                let coords = j + i * boardSize;
                squares[coords].classList.add("target");
            }
        }
    }
}

function removeTarget() {
    const targets = board.getElementsByClassName("target");
    Array.from(targets).forEach(el => el.classList.remove("target"));
}

function sendMoveToServer(piece, target) {
    let pos = BigInt(1) << BigInt(piece.x + piece.y * boardSize);
    let destination = BigInt(1) << BigInt(target.X + target.Y * boardSize);
    const prevBoard = positionHistory[positionHistory.length - 1].board;

    isWaitingForBot = true;
    updateViewState();

    fetch("movePiece", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ piece: pos, move: destination })
    }).then(response => {
        if (!response.ok) {
            throw new Error("Invalid move");
        }
        return parseBigIntJson(response);
    })
        .then(data => {
            console.log("Move successful:", data);

            const playerBoard = data.PlayerBoard || data.Board;
            const playerFrom = data.PlayerMove ? BigInt(data.PlayerMove[0]) : pos;
            const playerTo = data.PlayerMove ? BigInt(data.PlayerMove[1]) : destination;
            let playerNotation = formatMove(prevBoard, playerFrom, playerTo);
            if (data.Status === "you won") {
                playerNotation += "#";
            }

            positionHistory.push({
                board: playerBoard,
                from: playerFrom,
                to: playerTo,
                notation: playerNotation,
                color: playerColor
            });

            if (data.BotMove && BigInt(data.BotMove[0]) !== 0n) {
                const botFrom = BigInt(data.BotMove[0]);
                const botTo = BigInt(data.BotMove[1]);
                const botColor = playerColor === WHITE ? BLACK : WHITE;
                let botNotation = formatMove(playerBoard, botFrom, botTo);
                if (data.Status === "you lost") {
                    botNotation += "#";
                }
                positionHistory.push({
                    board: data.Board,
                    from: botFrom,
                    to: botTo,
                    notation: botNotation,
                    color: botColor
                });
            }

            if (playerColor == WHITE) {
                moveList = data.ValidWhiteMoves;
            } else {
                moveList = data.ValidBlackMoves;
            }

            isWaitingForBot = false;
            if (data.Status === "you won" || data.Status === "you lost" || data.Status === "stalemate") {
                gameOver = true;
            }

            goToMove(positionHistory.length - 1);

            if (data.Status === "you won") {
                setTimeout(() => alert("Check mate! You won!"), 50);
            } else if (data.Status === "stalemate") {
                setTimeout(() => alert("Stalemate!"), 50);
            } else if (data.Status === "you lost") {
                setTimeout(() => alert("Check mate, you lost!"), 50);
            } else {
                checkMoveList(moveList);
            }
        })
        .catch(error => {
            console.error("Error:", error);
            isWaitingForBot = false;
            goToMove(viewIndex);
        });
}

function updatePieces(b) {
    document.querySelectorAll(".piece").forEach(piece => {
        piece.remove();
    });
    pieceData.clear();
    let squares = document.querySelectorAll(".square");
    for (let i = 0; i < boardSize; i++) {
        for (let j = 0; j < boardSize; j++) {
            const { rank, file } = domToBoardCoords(i, j);
            let pieceName = getPiece(b, rank, file);
            if (pieceName != null) {
                let square = squares[j + i * boardSize];
                let piece = document.createElement("div");
                piece.classList.add("piece");
                piece.classList.add(pieceName);
                square.appendChild(piece);
                pieceData.set(piece, { x: file, y: rank, dx: 0, dy: 0, color: pieceName[0] == "w" ? WHITE : BLACK });
            }
        }
    }
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
    const isLatest = viewIndex === positionHistory.length - 1;
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

function checkMoveList(ml) {
    for (let key in ml) {
        if (ml[key] != 0) {
            return;
        }
    }
    gameOver = true;
    updateViewState();
    setTimeout(() => alert("Check mate, you lost!"), 50);
}

document.getElementById("btn-first").addEventListener("click", () => goToMove(0));
document.getElementById("btn-prev").addEventListener("click", () => goToMove(viewIndex - 1));
document.getElementById("btn-next").addEventListener("click", () => goToMove(viewIndex + 1));
document.getElementById("btn-last").addEventListener("click", () => goToMove(positionHistory.length - 1));

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
        goToMove(positionHistory.length - 1);
    }
});

initBoard();
