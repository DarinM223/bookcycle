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
	$('.messages').append('<span id="notification_count">' + data.length + '</span>');
	for(var i = 0; i < data.length; i++) {
		message.push(data[i]['message']);
		senderIdList.push(data[i]['senderId']);
		// console.log(senderIdList);
		(function(i) {
			$.ajax({
				type: 'GET',
				url: '/users/' + data[i]['senderId'] + '/json'
			}).success(function(data, textStatus, jqXHR) {
				senderName = data['first_name'] + " " + data['last_name'];
				console.log(senderIdList[i]);
				$('#notificationsBody').append('<a href="/message/' + senderIdList[i] +'"><div id="msg"><strong>' + senderName + '</strong><br><br>' + message[i] +'</div></a>');
			});
		})(i);
	}
}).error(function(jqXHR, textStatus, err) {
  console.log(err);
});

});
