(function () {
    const cells = 8;
    let html = '';
    contain = document.querySelector('#board-container');
    console.log(contain);
    for(let i=0; i < cells; i++) {
        html += '<tr>'
        for(let j=0; j<cells; j++) {
            let colorClass = ((i+j) % 2 == 0) ? 'black' : 'white';
            html += '<td class="' + colorClass + '"></td>';
        }
        html += '</tr>'
    }
    contain.innerHTML = html;
})();

document.getElementById('board-container').addEventListener('click', pieceSelect);

function pieceSelect(e) {
    console.log(e.target.cellIndex);
    console.log(e.target.parentNode.rowIndex);
}


