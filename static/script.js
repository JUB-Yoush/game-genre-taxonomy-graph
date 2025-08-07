var toggleVisibility = function (elem) {
  //elem.nextElementSibling.classList.toggle("closed");
}

let recalculateLines = () => {
  let puzzle = document.getElementById("puzzle");
  let sokoban = document.getElementById("sokoban");
  let falling_block_puzzle = document.getElementById("falling_block_puzzle");
  new LeaderLine(puzzle, sokoban, { color: 'red', size: 4 })
  new LeaderLine(puzzle, falling_block_puzzle, { color: 'red', size: 4 })

}

recalculateLines()
