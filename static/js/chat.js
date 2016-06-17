/* global $, WebSocket */

/**
 * Sets up a chatroom with two people
 * @param {integer} senderId the userid of the currently logged in user
 * @param {integer} receiverId the userid of the user to message
 */
function Message (senderId, receiverId) {
  $(document).ready(function () {
    var conn
    var msg = $('#msg')
    var log = $('#log')
    var wrapper, wrapDiv, messageDiv, messageTextNode

    function appendLog (msg) {
      msg.appendTo(log)
    }

    function addMessage (msg) {
      if (typeof (msg.latitude) !== 'undefined' && msg.latitude !== 0 &&
          typeof (msg.longitude) !== 'undefined' && msg.longitude !== 0) { // if location change
        wrapper = document.createElement('div')
        wrapDiv = document.createElement('div')
        wrapDiv.className = 'chat-messages-wrapper'

        messageDiv = document.createElement('div')
        var strongMessageNode = document.createElement('strong')
        messageTextNode = document.createTextNode(msg.message)
        messageDiv.style['text-align'] = 'center'
        messageDiv.style['margin-left'] = 'auto'
        messageDiv.style['margin-right'] = 'auto'
        strongMessageNode.appendChild(messageTextNode)
        messageDiv.appendChild(strongMessageNode)

        wrapDiv.appendChild(messageDiv)
        wrapper.appendChild(wrapDiv)
        appendLog($(wrapper.innerHTML))
      } else { // if chat message
        wrapper = document.createElement('div')
        wrapDiv = document.createElement('div')
        wrapDiv.className = 'chat-messages-wrapper'

        messageDiv = document.createElement('div')
        messageTextNode = document.createTextNode(msg.message)
        messageDiv.className = 'chat-message ' + (msg.senderId === senderId ? 'to' : 'from')
        messageDiv.appendChild(messageTextNode)

        wrapDiv.appendChild(messageDiv)
        wrapper.appendChild(wrapDiv)
        appendLog($(wrapper.innerHTML))
      }
    }

    $.ajax({
      type: 'GET',
      url: '/past_messages/' + receiverId
    }).success(function (data, textStatus, jqXHR) {
      var parsedResults = data
      if (parsedResults !== null) {
        for (var i = parsedResults.length - 1; i >= 0; i--) {
          addMessage(parsedResults[i])
        }
      }
      $('#log').scrollTop($('#log')[0].scrollHeight)
    }).error(function (jqXHR, textStatus, err) {
      console.log(err)
    })

    $.ajax({
      type: 'GET',
      url: '/users/' + receiverId + '/json'
    }).success(function (data) {
      $('#title').text('Messaging with ' + data.first_name + ' ' + data.last_name)
    }).error(function (j, t, err) {
      console.log(err)
    })

    $('#send').click(function (e) {
      e.preventDefault()
      if (!conn) {
        return false
      }
      if (!msg.val()) {
        return false
      }

      var parsedMessage = {
        senderId: senderId,
        receiverId: receiverId,
        message: msg.val()
      }

      conn.send(JSON.stringify(parsedMessage))
      addMessage(parsedMessage)
      $('#log').scrollTop($('#log')[0].scrollHeight)
      msg.val('')
      return false
    })

    if (window.WebSocket) {
      conn = new WebSocket('ws://' + window.location.host + '/ws')

      conn.onclose = function (evt) {
        appendLog($('<div><b>Connection closed.</b></div>'))
      }

      conn.onmessage = function (evt) {
        var parsedMessage = JSON.parse(evt.data)
        addMessage(parsedMessage)
        $('#log').scrollTop($('#log')[0].scrollHeight)
      }
    } else {
      appendLog($('<div><b>Your browser does not support WebSockets.</b></div>'))
    }
  })
}
