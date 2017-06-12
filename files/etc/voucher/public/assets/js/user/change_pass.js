$( window ).ready(function() {
});

$('#user-recover').submit(function(e) {
     e.preventDefault();
     e.returnValue = false;
});

function recover(){
  if($("#newpass").val() != $("#confpass").val()){
    alert("Passwords are not match.");
    return;
  }

  var user = {
      old_password: $("#oldpass").val(),
      new_password: $("#newpass").val()
    };

  $.ajax({
      url: '/v1/ui/user/update?type=password&token='+token,
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(user),
      success: function (data){
        alert("Password Updated");
        window.location = "/user/profile?token="+token;
      },
      error: function (data){
        alert("Failed");
      }
  });
}
