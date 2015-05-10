$.ajax({
  type: 'GET',
  url: '/unread_messages'
}).success(function(data, textStatus, jqXHR) {
  console.log(data);
}).error(function(jqXHR, textStatus, err) {
  console.log(err);
});
