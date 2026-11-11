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

let gameOver = false;
let isWaitingForBot = false;

function canPlay() {
    return !gameOver && !isWaitingForBot && isViewingLatest();
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
    gameOver = false;
    isWaitingForBot = false;
    initHistory(startBoard, r2json.openingMove, boardSetup);
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
            recordMove(playerBoard, playerFrom, playerTo, playerColor, data.Status === "you won" ? "#" : "");

            if (data.BotMove && BigInt(data.BotMove[0]) !== 0n) {
                const botFrom = BigInt(data.BotMove[0]);
                const botTo = BigInt(data.BotMove[1]);
                const botColor = playerColor === WHITE ? BLACK : WHITE;
                recordMove(data.Board, botFrom, botTo, botColor, data.Status === "you lost" ? "#" : "");
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

            goToLatest();

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
            refreshCurrentMove();
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

initBoard();
