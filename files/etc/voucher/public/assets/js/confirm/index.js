$(document).ready(function(){
  function setExpiryYear() {
    var currYear = (new Date).getFullYear();
    for(var i = 0; i < 6; i++) {
      var val = ((currYear+i)+'').substring(2);
      $('.expiry_year').append('<option value="' + val + '">' + val + '</option>');
    }
  }

  function enablePlugins() {
    $('[data-toggle="tooltip"]').tooltip();
  }

  setExpiryYear();
  enablePlugins();
});
