const BLACK = 1, WHITE = 0;
// BigInt cannot be converted to json normally
BigInt.prototype.toJSON = function () {
    return JSON.rawJSON(this.toString());
};


async function initBoard() {
    const response = await fetch("getBoard");
    boardSetup = await response.json();

    const response2nd = await fetch("getMoves");
    let r2json = await response2nd.json()
    playerColor = await r2json.playerColor
    document.getElementById("name").innerText = (playerColor == WHITE) ? "You are playing White" : "You are playing Black"
    moveList = await r2json.validMoves

    for (let i = 0; i < boardSize; i++) {
        let row = document.createElement("div");
        row.className = "row";
        for (let j = 0; j < boardSize; j++) {
            let square = document.createElement("div");
            let pieceName = getPiece(boardSetup, i, j);
            if (pieceName != null) {
                let piece = document.createElement("div");
                piece.classList.add("piece");
                piece.classList.add(pieceName);
                square.appendChild(piece);
                pieceData.set(piece, { x: j, y: i, dx: 0, dy: 0, color: pieceName[0] == "w" ? WHITE : BLACK });
            }

            square.className = "square";
            square.className += (i + j) % 2 === 0 ? " white" : " black";
            square.x = j;
            square.y = i;
            row.appendChild(square);
        }
        board.appendChild(row);
    }
}

function getPiece(boardSetup, row, col) {
    for (let key in boardSetup) {
        let bitboard = BigInt(boardSetup[key]);
        let pos = BigInt((col + row * boardSize));
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
            removeTarget()
            event.target.classList.add("dragged");
            addTarget(event)
        },
        move(event) {
            let data = pieceData.get(event.target);
            if (data.color != playerColor) {
                return
            }
            data.dx += event.dx;
            data.dy += event.dy;
            event.target.style.transform =
                `translate(${-data.dx}px, ${-data.dy}px) rotate(180deg)`;
        },
        end(event) {
            removeTarget()
            event.target.classList.remove("dragged");
            if (!event.dropzone) {
                event.target.style.transform = "rotate(180deg)";
                let data = pieceData.get(event.target);
                data.dx = 0;
                data.dy = 0;
            }
        }
    }
});

interact(".target")
    .dropzone({
        accept: ".piece",
        overlap: 0.6,
        ondrop: function (event) {
            if (event.target.firstChild) {
                event.target.removeChild(event.target.firstChild);
            }
            event.target.appendChild(event.relatedTarget);
            event.relatedTarget.style.transform = "rotate(180deg)";
            let piece = pieceData.get(event.relatedTarget);
            sendMoveToServer(piece, { X: event.target.x, Y: event.target.y });
            piece.dx = 0;
            piece.dy = 0;
            piece.x = event.target.x;
            piece.y = event.target.y;
        }
    })

interact(".piece").on("tap", addTarget);

function addTarget(event) {
    let data = pieceData.get(event.target);
    let pos = BigInt(1) << BigInt(data.x + data.y * boardSize)
    let moves;
    for (key in moveList) {
        if (key == pos) {
            moves = BigInt(moveList[key])
            console.log(key)
            break;
        }
    }
    if (!moves) {
        return;
    }
    for (let i = 0; i < boardSize; i++) {
        for (let j = 0; j < boardSize; j++) {
            let targetPos = BigInt(1) << BigInt(j + i * boardSize);
            if ((moves & targetPos) != 0n) {
                let coords = j + i * boardSize;
                board.getElementsByClassName("square")[coords].classList.add("target");
            }
        }
    }
}

function removeTarget() {
    const targets = board.getElementsByClassName("target");
    Array.from(targets).forEach(el => el.classList.remove("target"));
}

function sendMoveToServer(piece, target) {
    let pos = BigInt(1) << BigInt(piece.x + piece.y * boardSize)
    let destination = BigInt(1) << BigInt(target.X + target.Y * boardSize);
    fetch("movePiece", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ piece: pos, move: destination })
    }).then(response => response.json())
        .then(data => {
            console.log("Move successful:", data);
            if (playerColor == WHITE) {
                moveList = data.ValidWhiteMoves
            } else {
                moveList = data.ValidBlackMoves
            }
            updatePieces(data.Board);
        })
        .catch(error => console.error("Error:", error));
}

function updatePieces(b) {
    document.querySelectorAll(".piece").forEach(piece => {
        piece.remove();
    });
    pieceData.clear();
    let squares = document.querySelectorAll(".square");
    for (let i = 0; i < boardSize; i++) {
        for (let j = 0; j < boardSize; j++) {
            let pieceName = getPiece(b, i, j);
            if (pieceName != null) {
                let square = squares[j + i * boardSize];
                let piece = document.createElement("div");
                piece.classList.add("piece");
                piece.classList.add(pieceName);
                square.appendChild(piece);
                pieceData.set(piece, { x: j, y: i, dx: 0, dy: 0, color: pieceName[0] == "w" ? WHITE : BLACK });
            }
        }
    }
}

let pieceData = new Map();
let boardSetup;
const boardSize = 8;
let playerColor;
let moveList;
let board = document.getElementById("board");
initBoard();