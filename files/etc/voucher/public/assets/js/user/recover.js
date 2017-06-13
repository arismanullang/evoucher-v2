$( window ).ready(function() {
  var password = findGetParameter("key");
  if( password == null ){
        window.location = "/user/login";
  }

});

$('#user-recover').submit(function(e) {
     e.preventDefault();
     e.returnValue = false;
});

function recover(){
  var key = findGetParameter("key");
  if($("#password1").val() != $("#password2").val()){
    alert("Passwords are not match.");
    return;
  }

  var user = {
      password: $("#password1").val()
    };

  $.ajax({
      url: '/v1/ui/user/forgot/password?key='+key,
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(user),
      success: function (data){
        alert("Password Updated");
        window.location = "/user/login";
      },
      error: function (data){
        alert("Failed");
      }
  });
}

function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    location.search
    .substr(1)
        .split("&")
        .forEach(function (item) {
        tmp = item.split("=");
        if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
    });
    return result;
}
