const BLACK = 0, WHITE = 1;
const PIECE_NAMES = ["rook", "knight", "bishop", "queen", "king"];

async function initBoard() {
    const response = await fetch("getBoard");
    boardSetup = await response.json();
    boardSize = parseInt(boardSetup.length);

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
                piece.className = "piece" + pieceName;
                square.appendChild(piece);
                pieceData.set(piece, { x: j, y: i, dx: 0, dy: 0, color: boardSetup[i][j].Color });
            }

            square.className = "square";
            square.className += (i + j) % 2 === 0 ? " white" : " black";
            row.appendChild(square);
        }
        board.appendChild(row);
    }
}

function getPiece(boardSetup, row, col) {
    if (boardSetup[row][col] === null) {
        return null;
    }
    let color = boardSetup[row][col].Color;
    if (color == BLACK) {
        if (boardSetup[row][col].PieceType == undefined) {
            return " black-pawn";
        } else {
            return " black-" + PIECE_NAMES[boardSetup[row][col].PieceType];
        }
    } else {
        if (boardSetup[row][col].PieceType == undefined) {
            return " white-pawn";
        } else {
            return " white-" + PIECE_NAMES[boardSetup[row][col].PieceType];
        }
    }
}

interact(".piece").draggable({
    listeners: {
        start(event) {
            removeTarget()
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
                `translate(${data.dx}px, ${data.dy}px)`;
        },
        end(event) {
            removeTarget()
            if (!event.dropzone) {
                event.target.style.transform = "none";
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
            event.target.appendChild(event.relatedTarget);
            event.relatedTarget.style.transform = "none";
            let pos = pieceData.get(event.relatedTarget);
            pos.dx = 0;
            pos.dy = 0;
        }
    })

interact(".piece").on("tap", addTarget);

function addTarget(event) {
    let data = pieceData.get(event.target);
    for (let i = 0; i < moveList.length; i++) {
        let m = moveList[i]
        if (data.x == m.P.Pos.X && data.y == m.P.Pos.Y) {
            for (let j = 0; j < m.Moves.length; j++) {
                let coords = m.Moves[j].X + m.Moves[j].Y * 8
                console.log(coords)
                board.getElementsByClassName("square")[coords].classList.add("target");
            }
        }
    }
}

function removeTarget() {
    const targets = board.getElementsByClassName("target");
    Array.from(targets).forEach(el => el.classList.remove("target"));
}

let pieceData = new Map();
let boardSetup;
let boardSize;
let playerColor;
let moveList
let board = document.getElementById("board");
initBoard();