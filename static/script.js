var toggleVisibility = function (elem) {
  elem.nextElementSibling.classList.toggle("closed")
}
let genreQueue = []

// bfs style row by row rendering
// render a button
// add children to queue
// render children, connect child to parent (if parent isn't toggeled connect to the button)
// zooming in and out?
