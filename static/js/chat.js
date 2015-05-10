function Message(sender_id, receiver_id) {
  console.log(sender_id);
  console.log(receiver_id);

  $(function() {
    var conn;
    var msg = $("#msg");
    var log = $("#log");
  
    function appendLog(msg) {
      var d = log[0]
      var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
      msg.appendTo(log)
      if (doScroll) {
        d.scrollTop = d.scrollHeight - d.clientHeight;
      }
    }
  
    $("#form").submit(function() {
      if (!conn) {
        return false;
      }
      if (!msg.val()) {
        return false;
      }
      conn.send(JSON.stringify({
        senderId: sender_id,
        receiverId: receiver_id,
        message: msg.val()
      }));
      msg.val("")
      return false
    });
  
    if (window["WebSocket"]) {
      conn = new WebSocket("ws://localhost:8080/ws");
      conn.onclose = function(evt) {
        appendLog($("<div><b>Connection closed.</b></div>"))
      };
      conn.onmessage = function(evt) {
        console.log(evt.data);
        var parsedMessage = JSON.parse(evt.data);
        appendLog($("<div/>").text("Id: " + parsedMessage.receiverId + " Message: " + parsedMessage.message));
      };
    } else {
      appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
  });
}
