const appDiv: HTMLElement = document.getElementById('app')!;
appDiv.innerHTML = `<h1>Tic Tac Toe</h1>`;

const conn = new WebSocket("ws://" + document.location.host + "/ws");

type GameState = string;

enum CellType {
  Empty = 'E',
  Cross = 'X',
  Circle = 'O',
}

let boardState: Array<string> | null = null;
let gameId = null;
let isPlayerMove = false;
let playerChar = null;

const button = document.createElement('button');
button.innerHTML = 'Search';
button.addEventListener('click', (ev) => {
  console.log('searching', conn, ev);
  conn.send('search');
})

const search = document.createElement('p')
search.innerHTML = 'searching...'

const parseWebSocketMessage = (msg: string) => {
  const split_msg = msg.split('/');
  if (split_msg.length === 6) {
    const gameState = msg.split('/');
    gameId = gameState[0];
    const row1 = gameState[1];
    const row2 = gameState[2];
    const row3 = gameState[3];
    boardState = [row1, row2, row3];
    isPlayerMove = gameState[4] === "true";
    playerChar = gameState[5];
    constructGrid(boardState);
    appDiv.removeChild(search);
  } else if (msg == "searching") {
    appDiv.appendChild(search)
    appDiv.removeChild(button);
  } 
};

const constructGameState = (boardState: string[]) => {
  return boardState.reduce((acc, curr, index) => {
    if (index === 0) {
      return curr;
    }
    return acc + "/" + curr;
  }, "");
}

const constructGrid = (boardState: string[]) => {
  const gridDiv: HTMLElement = document.createElement('table');
  gridDiv.className = 'grid-table';
  if (boardState) {
    for (let i = 0; i < boardState.length; ++i) {
      const rowDiv = document.createElement('tr');
      rowDiv.className = 'grid-row';
      for (let j = 0; j < boardState[i].length; ++j) {
        const cell = document.createElement('td');
        const cellValue = boardState[i][j];
        cell.className = 'grid-cell';
        if (cellValue === '.') {
          cell.innerHTML = CellType.Empty;
        } else if (cellValue == '*') {
          cell.innerHTML = CellType.Cross;
        } else {
          cell.innerHTML = CellType.Circle;
        }
        cell.addEventListener('click', (ev) => {
          if (cell.innerHTML == CellType.Empty && isPlayerMove) {
            conn.send(gameId + "/" + constructGameState(boardState))
          }
        });
        rowDiv.appendChild(cell);
      }
      gridDiv.appendChild(rowDiv);
    }
  } else {
    return;
  }
  appDiv.appendChild(gridDiv);
};

conn.onmessage = (ev) => parseWebSocketMessage(ev.data);
appDiv.appendChild(button);