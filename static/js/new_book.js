$(document).ready(function() {
  var departmentSelector = '#department',
      courseIDSelector = '#courseID',
      professorSelector = '#professor',
      isbnSelector = "#isbn_enter",
      titleSelector = "#title_enter",
      isbnSearchSelector = '#search_button',
      submitFormSelector = '#submit_form';

  function SearchReplace(wildcard, url, urlEncodedQuery) {
    var department = $(departmentSelector).val();
    var courseID = $(courseIDSelector).val();
    var professor = $(professorSelector).val();

    var result;
    url = url.replace(wildcard, encodeURIComponent(urlEncodedQuery));
    if (wildcard == '%DEPARTMENT') {
      result = url.replace('%COURSEID', encodeURIComponent(courseID))
                  .replace('%PROFESSOR', encodeURIComponent(professor));
    } else if (wildcard == '%COURSEID') {
      result = url.replace('%DEPARTMENT', encodeURIComponent(department))
                .replace('%PROFESSOR', encodeURIComponent(professor));
    } else if (wildcard == '%PROFESSOR') {
      result = url.replace('%COURSEID', encodeURIComponent(courseID))
                .replace('%DEPARTMENT', encodeURIComponent(department));
    } else {
      throw new Error('Wildcard error');
    }
    console.log(result);
    return result;
  }

  var departmentSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('department'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      cache: false,
      url: '/course_search.json?type=department&department=%DEPARTMENT&course_id=%COURSEID&professor=%PROFESSOR',
      wildcard: '%DEPARTMENT',
      replace: SearchReplace.bind(null, '%DEPARTMENT')
    }
  });

  var courseIDSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('course_id'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      cache: false,
      url: '/course_search.json?type=course&department=%DEPARTMENT&course_id=%COURSEID&professor=%PROFESSOR',
      wildcard: '%COURSEID',
      replace: SearchReplace.bind(null, '%COURSEID')
    }
  });

  var professorSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('professor'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      cache: false,
      url: '/course_search.json?type=professor&department=%DEPARTMENT&course_id=%COURSEID&professor=%PROFESSOR',
      wildcard: '%PROFESSOR',
      replace: SearchReplace.bind(null, '%PROFESSOR')
    }
  });

  departmentSuggestion.initialize();
  courseIDSuggestion.initialize();
  professorSuggestion.initialize();

  $(departmentSelector).typeahead({
    highlight: true
  }, {
    display: 'department',
    source: departmentSuggestion.ttAdapter()
  });

  $(courseIDSelector).typeahead({
    highlight: true
  }, {
    display: 'course_id',
    source: courseIDSuggestion.ttAdapter()
  });

  $(professorSelector).typeahead({
    highlight: true
  }, {
    source: professorSuggestion.ttAdapter(),
    display: 'professor'
  });

  function getCourseID(department, courseID, professor, callback) {
    var url = '/course_search.json?type=professor&department=' + 
        encodeURIComponent(department) + '&course_id=' +
        encodeURIComponent(courseID) + '&professor=' +
        encodeURIComponent(professor);

    $.ajax({
      type: 'GET',
      url: url
    }).success(function(data) {
      if (data === null || data.length === 0) {
        callback(new Error('Course is invalid'), null);
      } else if (data.length > 1) {
        callback(new Error('Course is ambigous, did you fill out all of the course fields?'), null);
      } else {
        callback(null, data[0]);
      }
    });
  }

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

  $(isbnSearchSelector).click(function() {
    var isbn = $(isbnSelector).val();
    getISBN(isbn, function(err, data) {
      if (err) {
        alert(err.message);
        return;
      }

      if ($('#dropdown_isbn').hasClass('hidden')) {
        $('.book_title').text('Title: ' + data.volumeInfo.title);
        $('.book_authors').text('Authors: ' + data.volumeInfo.authors.join(', '));
        $('.book_date').text('Published Date: ' + data.volumeInfo.publishedDate);
        $('.book_image').attr('src', data.volumeInfo.imageLinks.thumbnail);

        $('#dropdown_isbn').removeClass('hidden');
      } else {
        $('#dropdown_isbn').addClass('hidden');
      }
    });
  });
  var canSubmit = false;

  $('#post-edit').submit(function(e) {
    if (!canSubmit) {
      e.preventDefault();
      if ($('#price').val().trim().length === 0) {
        alert('Price cannot be empty');
        return;
      }
      if ($('#condition').val().trim().length === 0) {
        alert('Condition cannot be empty');
        return;
      }
      if (isNaN(parseFloat($('#price').val()))) {
        alert('Price has to be a decimal');
        return;
      }
      if (isNaN(parseInt($('#condition').val(), 10))) {
        alert('Condition has to be an integer');
        return;
      }

      getCourseID($(departmentSelector).val(), $(courseIDSelector).val(), $(professorSelector).val(), function(err, course) {
        if (err) {
          alert(err.message);
          return;
        }

        $('#course_id').val(course.id);
        var isbn = $(isbnSelector).val();
        getISBN(isbn, function(err, book) {
          if (err) {
            alert(err.message);
            return;
          }

          $('#isbn').val(isbn);
          $('#title').val(book.volumeInfo.title);
          canSubmit = true;
          $('#post-edit').submit();
        });
      });
    }
  });
});
