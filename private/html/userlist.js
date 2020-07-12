const ws = new WebSocket("wss://chesslife.herokuapp.com/public/api/userwebsocket");

ws.addEventListener("open", () => {
    console.log('We are connected');
});
