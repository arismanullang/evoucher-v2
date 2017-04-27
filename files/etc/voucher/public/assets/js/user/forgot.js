$( window ).ready(function() {
});

$('#user-recover').submit(function(e) {
     e.preventDefault();
     e.returnValue = false;
});

function recover(){

  var user = $("#username").val();
  if( user == null ){
        alert("Username cant be empty.");
  }

  $.ajax({
      url: '/v1/api/mail?username='+user,
      type: 'get',
      dataType: 'json',
      success: function (data){
        window.location = "/user/mail-send";
      }
  });
}
