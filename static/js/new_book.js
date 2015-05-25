$(document).ready(function() {
  function SearchURL(departmentSelector,
                     courseIDSelector,
                     professorSelector) {

    this.departmentSelector = departmentSelector;
    this.courseIDSelector = courseIDSelector;
    this.professorSelector = professorSelector;
  }

  SearchURL.prototype.getDepartmentURL = function() {
    return '/course_search.json?' + 'department=%QUERY' +
                                  '&course_id=' + $(this.courseIDSelector).val() +
                                  '&professor=' + $(this.professorSelector).val();
  };

  SearchURL.prototype.getCourseIDURL = function() {
    return '/course_search.json?' + 'department=' + $(this.departmentSelector).val() +
                                  '&course_id=%QUERY' +
                                  '&professor=' + $(this.professorSelector).val();
  };

  SearchURL.prototype.getProfessorURL = function() {
    return '/course_search.json?' + 'department=' + $(this.departmentSelector).val() +
                                  '&course_id=' + $(this.courseIDSelector).val() +
                                  '&professor=%QUERY';
  };

  var departmentSelector = '#department',
      courseIDSelector = '#courseID',
      professorSelector = '#professor';

  var searchURL = new SearchURL(departmentSelector, courseIDSelector, professorSelector);

  var departmentSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('department'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      url: searchURL.getDepartmentURL(),
      wildcard: '%QUERY'
    }
  });
  var courseIDSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('course_id'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      url: searchURL.getCourseIDURL(),
      wildcard: '%QUERY',
    }
  });

  var professorSuggestion = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('professor'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      url: searchURL.getProfessorURL(),
      wildcard: '%QUERY'
    }
  });

  departmentSuggestion.initialize();
  courseIDSuggestion.initialize();
  professorSuggestion.initialize();

  $(departmentSelector).typeahead(null, {
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
