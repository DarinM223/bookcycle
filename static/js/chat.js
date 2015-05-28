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
    var recentLocation = null;


    function appendLog(msg) {
      msg.appendTo(log);
      //window.scrollTo(0, document.body.scrollHeight);
    }

    function addMessage(msg) {
      console.log(msg);
      if (typeof(msg.latitude) !== 'undefined' && msg.latitude !== 0 && 
          typeof(msg.longitude) !== 'undefined' && msg.longitude !== 0) { // if location change

        recentLocation = msg;
        var locationP = document.createElement('p');
        var _messageDiv = document.createTextNode(msg.message);
        locationP.appendChild(_messageDiv);
        appendLog($(locationP.innerHTML));
        // TODO: change location on map
      } else { // if chat message
        var wrapper = document.createElement('div');
        var wrapDiv = document.createElement('div');
        wrapDiv.className = 'chat-messages-wrapper';

        var messageDiv = document.createElement('div');
        var messageTextNode = document.createTextNode(msg.message);
        messageDiv.className = 'chat-message ' + (msg.senderId === sender_id ? 'to' : 'from');
        messageDiv.appendChild(messageTextNode);

        wrapDiv.appendChild(messageDiv);
        wrapper.appendChild(wrapDiv);
        appendLog($(wrapper.innerHTML));
      }
    }

    $.ajax({
      type: 'GET',
      url: '/past_messages/' + receiver_id
    }).success(function(data, textStatus, jqXHR) {
      var parsedResults = data;
      if (parsedResults !== null) {
        for (var i = parsedResults.length-1; i >= 0; i--) {
          addMessage(parsedResults[i]);
        }
      }
      //window.scrollTo(0, document.body.scrollHeight);
      $("#log").scrollTop($("#log")[0].scrollHeight);
    }).error(function(jqXHR, textStatus, err) {
      console.log(err);
    });

    $.ajax({
      type: 'GET',
      url: '/users/' + receiver_id + '/json'
    }).success(function(data) {
      $('#title').text('Messaging with ' + data.first_name + ' ' + data.last_name);
    }).error(function(j, t, err) {
      console.log(err);
    });

    $("#send").click(function(e) {
      e.preventDefault();
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
        $("#log").scrollTop($("#log")[0].scrollHeight);
      };
    } else {
      appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"));
    }
  });
}
