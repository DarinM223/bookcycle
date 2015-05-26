$(document).ready(function() {
  var departmentSelector = '#department',
      courseIDSelector = '#courseID',
      professorSelector = '#professor';

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
      url: '/course_search.json?department=%DEPARTMENT&course_id=%COURSEID&professor=%PROFESSOR',
      wildcard: '%DEPARTMENT',
      replace: SearchReplace.bind(null, '%DEPARTMENT')
    }
  });
  var courseIDSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('course_id'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      url: '/course_search.json?department=%DEPARTMENT&course_id=%COURSEID&professor=%PROFESSOR',
      wildcard: '%COURSEID',
      replace: SearchReplace.bind(null, '%COURSEID')
    }
  });

  var professorSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('professor_last_name', 'professor_first_name'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      url: '/course_search.json?department=%DEPARTMENT&course_id=%COURSEID&professor=%PROFESSOR',
      templates: {
        suggestion: function(data) {
          console.log(data);
          return '<div>' + data.professor_first_name + ' ' + data.professor_last_name + '</div>';
        }
      },
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

  $(courseIDSelector).typeahead(null, {
    display: 'course_id',
    source: courseIDSuggestion.ttAdapter(),
  });
  $(professorSelector).typeahead(null, {
    display: 'professor',
    source: professorSuggestion.ttAdapter(),
  });
});
