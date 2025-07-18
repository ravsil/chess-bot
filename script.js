let board = document.getElementById("board");

const boardSetup = [
    "rook", "knight", "bishop", "queen", "king", "bishop", "knight", "rook",
    "pawn", "pawn", "pawn", "pawn", "pawn", "pawn", "pawn", "pawn"]

let boardSize = 8;
let pieceData = new Map();

for (let i = 0; i < boardSize; i++) {
    let row = document.createElement("div");
    row.className = "row";
    for (let j = 0; j < boardSize; j++) {
        let square = document.createElement("div");
        let pieceName = getPiece(boardSetup, i, j);
        if (!pieceName.includes("undefined")) {
            let piece = document.createElement("div");
            piece.className = "piece" + pieceName;
            square.appendChild(piece);
            pieceData.set(piece, { dx: 0, dy: 0 });
        }

        square.className = "square";
        square.className += (i + j) % 2 === 0 ? " white" : " black";
        row.appendChild(square);
    }
    board.appendChild(row);
}

function getPiece(boardSetup, row, col) {
    let len = boardSetup.length;
    if (row < 8 - (len / 8)) {
        return " black-" + boardSetup[col + row * 8];
    } else if (row >= (len / 8)) {
        return " white-" + boardSetup[col + (7 - row) * 8];
    }
    return "";
}

interact(".piece").draggable({
    listeners: {
        move(event) {
            let pos = pieceData.get(event.target);
            pos.dx += event.dx;
            pos.dy += event.dy;
            event.target.style.transform =
                `translate(${pos.dx}px, ${pos.dy}px)`;
        },
        end(event) {
            if (!event.dropzone) {
                event.target.style.transform = "none";
                let pos = pieceData.get(event.target);
                pos.dx = 0;
                pos.dy = 0;
            }
        }
    }
});

interact(".square")
    .dropzone({
        accept: ".piece",
        overlap: 0.6,
        ondrop: function (event) {
            event.target.appendChild(event.relatedTarget);
            event.relatedTarget.style.transform = "none";
            let pos = pieceData.get(event.relatedTarget);
            pos.dx = 0;
            pos.dy = 0;
            pos.x
        }
    })