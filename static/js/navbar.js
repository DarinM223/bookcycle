$(document).ready(function() {
	// console.log(senderIdList);
	$(".messages").click(function() {
		$("#notificationContainer").fadeToggle(300);
		$("#notification_count").fadeOut("slow");
		return false;
	});

//Document Click
$(document).click(function() {
	$("#notificationContainer").hide();
});
// console.log(senderIdList);

function isInArray(value, array) {
  return array.indexOf(value) > -1;
}
setInterval(function() {
	// senderIdList = [];
  // console.log('Sending AJAX');
$.ajax({
  type: 'GET',
  url: '/messages'
}).success(function(data, textStatus, jqXHR) {
  $('.msg').remove();
  var messageNum = 0;
  var msgCounter = 0;
  var senderIdList = [];
  var message = [];
  var senderName;
  var read = false;
  
  for(var i = 0; i < data.length; i++) {
    // console.log(($.inArray(data[i]['senderId'], senderIdList)));
    if (($.inArray(data[i]['senderId'], senderIdList)) == -1) {
    	// console.log("here");
      if((data[i]['message']).length > 60) {
        message.push((data[i]['message']).substring(0,60) + "...");
      }
      else {
        message.push(data[i]['message']);
      }
      read = data[i]['read'];
      senderIdList.push(data[i]['senderId']);
      (function(messageNum, read) {
        $.ajax({
          type: 'GET',
          url: '/users/' + data[i]['senderId'] + '/json'
        }).success(function(data, textStatus, jqXHR) {
          senderName = data['first_name'] + " " + data['last_name'];
          // console.log(senderName);
          if(read) {
            $('#notificationsBody').append('<a class="msg" href="/message/' + senderIdList[messageNum] +'"><div class="readmsg"><span id="sendername">' + senderName + '</span><br><br><span id="msgpreview">' + message[messageNum] +'</span></div></a>');
          }
          else {
            $('#notificationsBody').append('<a class="msg" href="/message/' + senderIdList[messageNum] +'"><div class="unreadmsg"><span id="sendername">' + senderName + '</span><br><br><span id="msgpreview">' + message[messageNum] +'</span></div></a>');
          }
        });
      })(messageNum, read);
      messageNum++;
      if(!read) {
        msgCounter++;
      }
    }
  }
  // console.log(msgCounter);
  if (msgCounter != 0) {
    $('.messages').append('<span id="notification_count">' + msgCounter + '</span>');
  }
}).error(function(jqXHR, textStatus, err) {
  console.log(err);
});

}, 1000);

});
