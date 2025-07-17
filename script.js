let board = document.getElementById("board");

const boardSetup = [
    "rook", "knight", "bishop", "queen", "king", "bishop", "knight", "rook",
    "pawn", "pawn", "pawn", "pawn", "pawn", "pawn", "pawn", "pawn"]

let boardSize = 8;
for (let i = 0; i < boardSize; i++) {
    let row = document.createElement("div");
    row.className = "row";
    for (let j = 0; j < boardSize; j++) {
        let square = document.createElement("div");
        let pieceName = getPiece(boardSetup, i, j);
        if (pieceName != "") {
            let piece = document.createElement("div");
            piece.className = "piece" + pieceName;
            square.className = "square";
            square.appendChild(piece);
        }

        square.className += (i + j) % 2 === 0 ? " white" : " black";
        row.appendChild(square);
    }
    board.appendChild(row);
}

function getPiece(boardSetup, row, col) {
    let len = boardSetup.length;
    console.log(row, 4 - (len / 8), 4 + (len / 8))
    if (row < 8 - (len / 8)) {
        return " black-" + boardSetup[col + row * 8];
    } else if (row >= (len / 8)) {
        return " white-" + boardSetup[col + (7 - row) * 8];
    }
    return "";
}