/**
 * Sets up a chatroom with two people
 * @param {integer} sender_id the userid of the currently logged in user
 * @param {integer} receiver_id the userid of the user to message
 */
function Message(sender_id, receiver_id) {
  $(document).ready(function() {
    var conn;
    var msg = $("#msg");
    var log = $("#log");


    function appendLog(msg) {
      var d = log[0];
      var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
      msg.appendTo(log);
      if (doScroll) {
        $(d).animate({scrollTop: d.scrollHeight - d.clientHeight}, "slow");
      }
    }

    function addMessage(msg) {
      var wrapper = document.createElement('div');
      var clearDiv = document.createElement('div');
      clearDiv.className = 'clear';

      var messageDiv = document.createElement('div');
      var messageTextNode = document.createTextNode(msg.message);
      messageDiv.className = (msg.senderId === sender_id ? 'from-me' : 'from-them');
      messageDiv.appendChild(messageTextNode);

      clearDiv.appendChild(messageDiv);
      wrapper.appendChild(clearDiv);
      appendLog($(wrapper.innerHTML));
    }

    $(document).ready(function() {
      $.ajax({
        type: 'GET',
        url: '/past_messages/' + receiver_id
      }).success(function(data, textStatus, jqXHR) {
        var parsedResults = JSON.parse(data);
        if (parsedResults !== null) {
          for (var i = parsedResults.length-1; i >= 0; i--) {
            addMessage(parsedResults[i]);
          }
        }
        if (log.length > 0) {
          var d = log[0];
          d.scrollTop = d.scrollHeight;
        }
      }).error(function(jqXHR, textStatus, err) {
        console.log(err);
      });
    });
  
    $("#form").submit(function() {
      if (!conn) {
        return false;
      }
      if (!msg.val()) {
        return false;
      }

      var parsedMessage = {
        senderId: sender_id,
        receiverId: receiver_id,
        message: msg.val()
      };

      conn.send(JSON.stringify(parsedMessage));
      addMessage(parsedMessage);
      $("#log").scrollTop($("#log")[0].scrollHeight);
      msg.val("");
      return false;
    });
  
    if (window.WebSocket) {
      conn = new WebSocket("ws://localhost:8080/ws");

      conn.onclose = function(evt) {
        appendLog($("<div><b>Connection closed.</b></div>"));
      };

      conn.onmessage = function(evt) {
        var parsedMessage = JSON.parse(evt.data);
        addMessage(parsedMessage);
      };
    } else {
      appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"));
    }
  });
}
