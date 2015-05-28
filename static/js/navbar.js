$(document).ready(function() {
	$(".messages").click(function() {
		$("#notificationContainer").fadeToggle(300);
		$("#notification_count").fadeOut("slow");
		return false;
	});

//Document Click
$(document).click(function() {
	$("#notificationContainer").hide();
});

$.ajax({
  type: 'GET',
  url: '/unread_messages'
}).success(function(data, textStatus, jqXHR) {
	var senderName;
	var message = [];
	var senderIdList = [];
	var messageNum = 0;
	for(var i = 0; i < data.length; i++) {
		// console.log(($.inArray(data[i]['senderId'], senderIdList)));
		if (($.inArray(data[i]['senderId'], senderIdList)) == -1) {
			if((data[i]['message']).length > 60) {
				message.push((data[i]['message']).substring(0,60) + "...");
			}
			else {
				message.push(data[i]['message']);
			}
			senderIdList.push(data[i]['senderId']);
			(function(messageNum) {
				$.ajax({
					type: 'GET',
					url: '/users/' + data[i]['senderId'] + '/json'
				}).success(function(data, textStatus, jqXHR) {
					senderName = data['first_name'] + " " + data['last_name'];
					// console.log(senderName);
					$('#notificationsBody').append('<a href="/message/' + senderIdList[messageNum] +'"><div id="msg"><span id="sendername">' + senderName + '</span><br><br><span id="msgpreview">' + message[messageNum] +'</span></div></a>');
				});
			})(messageNum);
			messageNum++;
		}
		
	}
	if (message.length != 0) {
		$('.messages').append('<span id="notification_count">' + message.length + '</span>');
	}
	
}).error(function(jqXHR, textStatus, err) {
  console.log(err);
});

});
