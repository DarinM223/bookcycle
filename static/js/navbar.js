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
//Popup Click
$("#notificationContainer").click(function() {
	return false
});


$.ajax({
  type: 'GET',
  url: '/unread_messages'
}).success(function(data, textStatus, jqXHR) {
  console.log(data);
}).error(function(jqXHR, textStatus, err) {
  console.log(err);
});

});
