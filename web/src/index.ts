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
let isPlayerMove = false;

const parseWebSocketMessage = (msg: string) => {
  const split_msg = msg.split('/');
  if (split_msg.length > 0) {
    const gameState = msg.split('/');
    const row1 = gameState[1];
    const row2 = gameState[2];
    const row3 = gameState[3];
    boardState = [row1, row2, row3];
    constructGrid(boardState);
  }
};

const constructGrid = (boardState) => {
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
          console.log('cell-val', cellValue);
          if (cell.innerHTML == CellType.Empty && isPlayerMove) {
            cell.innerHTML = CellType.Cross;
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

const searchButton = () => {
  const button = document.createElement('button');
  button.innerHTML = 'Search';
  button.addEventListener('click', (ev) => {
    console.log('searching', conn, ev);
    conn.send('search');
  })
  return button
}
conn.onmessage = (ev) => parseWebSocketMessage(ev.data);

appDiv.appendChild(searchButton());
