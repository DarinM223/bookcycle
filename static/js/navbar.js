/* global $, jQuery */

$(document).ready(function () {
  $('.messages').click(function () {
    $('#notificationContainer').fadeToggle(300)
    $('#notification_count').fadeOut('slow')
    return false
  })

  // Document Click
  $(document).click(function () {
    $('#notificationContainer').hide()
  })

  jQuery.fn.insertAt = function (index, element) {
    var lastIndex = this.children().size()
    if (index < 0) {
      index = Math.max(0, lastIndex + 1 + index)
    }
    this.append(element)
    if (index < lastIndex) {
      this.children().eq(index).before(this.children().last())
    }
    return this
  }

  setInterval(function () {
    $.ajax({
      type: 'GET',
      url: '/messages'
    }).success(function (data, textStatus, jqXHR) {
      $('.msg').remove()
      var messageNum = 0
      var msgCounter = 0
      var senderIdList = []
      var message = []
      var senderName
      var read = false

      for (var i = 0; i < data.length; i++) {
        read = data[i]['read']
        if (!read) {
          msgCounter++
        }

        if (($.inArray(data[i]['senderId'], senderIdList)) === -1) {
          if ((data[i]['message']).length > 60) {
            message.push((data[i]['message']).substring(0, 60) + '...')
          } else {
            message.push(data[i]['message'])
          }
          senderIdList.push(data[i]['senderId'])
          ;(function (messageNum, read) {
            $.ajax({
              type: 'GET',
              url: '/users/' + data[i]['senderId'] + '/json'
            }).success(function (data, textStatus, jqXHR) {
              senderName = data['first_name'] + ' ' + data['last_name']
              if (read) {
                $('#notificationsBody').insertAt(
                  messageNum,
                  '<a class="msg" href="/message/' +
                    senderIdList[messageNum] +
                    '"><div class="readmsg"><span id="sendername">' +
                    senderName + '</span><br><br><span id="msgpreview">' +
                    message[messageNum] +
                    '</span></div></a>'
                )
              } else {
                $('#notificationsBody').insertAt(
                  messageNum,
                  '<a class="msg" href="/message/' +
                    senderIdList[messageNum] +
                    '"><div class="unreadmsg"><span id="sendername">' +
                    senderName +
                    '</span><br><br><span id="msgpreview">' +
                    message[messageNum] +
                    '</span></div></a>'
                )
              }
            })
          })(messageNum, read)
          messageNum++
        }
      }

      if (msgCounter !== 0) {
        $('.messages').append('<span id="notification_count">' + msgCounter + '</span>')
      }
    }).error(function (jqXHR, textStatus, err) {
      console.log(err)
    })
  }, 1000)
})
