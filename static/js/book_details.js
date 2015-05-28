  function getISBN(isbn, callback) {
    if (isbn.trim().length === 0) {
      return callback(Error('ISBN field is empty'), null);
    }
    $.ajax({
      method: 'GET',
      url: 'https://www.googleapis.com/books/v1/volumes?q=isbn:' + encodeURIComponent(isbn)
    }).success(function(data) {
      if (data.totalItems === 0) {
        callback(new Error('ISBN is invalid'), null);
      } else if (data.totalItems > 1) {
        callback(new Error('ISBN is ambiguous'), null);
      } else {
        callback(null, data.items[0]);
      }
    });
  }

function DisplayBook(isbn, courseid) {
	$(document).ready(function() {


    getISBN(isbn, function(err, data) {
      if (err) {
        alert(err.message);
        return;
      }
    $('.author').text('Author(s): ' + data.volumeInfo.authors.join(', '));

    if (typeof(data.volumeInfo.imageLinks) === 'undefined' || 
            data.volumeInfo.imageLinks === null || 
            typeof(data.volumeInfo.imageLinks.thumbnail) === 'undefined' ||
            data.volumeInfo.imageLinks.thumbnail === null) {

          $('.bookpicture').attr('src', '/images/no_image.png');
        } else {
          $('.bookpicture').attr('src', data.volumeInfo.imageLinks.thumbnail);
        }

    $.ajax({
	  type: 'GET',
	  url: '/courses/'+courseid+'/json'
	}).success(function(data, textStatus, jqXHR) {
		$('.department').text('Department: ' + data['department']);
		$('.course_id').text('Class: ' + data['course_id']);
		$('.professor').text('Professor: ' + data['professor']);
	}).error(function(jqXHR, textStatus, err) {
  		console.log(err);
	});


  	});



	});
}
