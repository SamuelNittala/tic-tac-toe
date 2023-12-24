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
button.id = 'search-btn'
button.addEventListener('click', (ev) => {
  console.log('searching', conn, ev);
  conn.send('search');
})

const search = document.createElement('p')
search.innerHTML = 'searching...'
search.id = 'search-tag'

const gridDiv: HTMLElement = document.createElement('table');
gridDiv.id = "game-grid"

const removeSearchUI = () => {
  if (document.getElementById('search-tag')) {
    appDiv.removeChild(search)
  }
  if (document.getElementById('search-btn')) {
    appDiv.removeChild(button)
  }
}

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
    removeSearchUI();
  } else if (msg == "searching") {
    appDiv.appendChild(search);
    if (document.getElementById('search-btn')) {
      appDiv.removeChild(button)
    }
  } else {
    appDiv.removeChild(gridDiv)
    const gameMessage = document.createElement('p')
    const gameState = msg.split('/');
    gameMessage.innerHTML = gameState[0] 
    appDiv.appendChild(gameMessage)
  } 
};

const constructGameState = (row, col) => {
  const newBoardState = [];
  for (let i = 0; i < boardState.length; ++i) {
    const newRow = [];
    for (let j  = 0; j < boardState[i].length; ++j) {
      if (i === row && j === col) {
        newRow.push(playerChar);
      } else {
        newRow.push(boardState[i][j])
      }
    }
    newBoardState.push(newRow);
  }
  return newBoardState[0].join('') + "/" + newBoardState[1].join('') + "/" + newBoardState[2].join('') + "/"; 
}


const constructGrid = (boardState: string[]) => {
  gridDiv.className = isPlayerMove ? 'grid-table active' : 'grid-table wait';
  if (boardState) {
    for (let i = 0; i < boardState.length; ++i) {
      const rowId = `grid-row-${i+1}`;
      let rowDiv = document.getElementById(rowId);
      if (!rowDiv) {
        rowDiv = document.createElement('tr')
        rowDiv.className = 'grid-row';
        rowDiv.id = rowId 
      }
      for (let j = 0; j < boardState[i].length; ++j) {
        const cellId = `grid-cell-${i+1}-${j+1}`
        let cell = document.getElementById(cellId);
        if (!cell) {
            cell = document.createElement('td');
            cell.className = 'grid-cell';
            cell.id = cellId
        }
        const cellValue = boardState[i][j];
        if (cellValue === '.') {
          cell.innerHTML = CellType.Empty;
        } else if (cellValue == '*') {
          cell.innerHTML = CellType.Cross;
        } else {
          cell.innerHTML = CellType.Circle;
        }
        cell.addEventListener('click', (ev) => {
          if (cell.innerHTML == CellType.Empty && isPlayerMove) {
            conn.send(gameId + "/" + constructGameState(i, j))
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