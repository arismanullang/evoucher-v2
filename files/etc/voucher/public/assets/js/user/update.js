$( document ).ready(function() {
  getUserDetails();

  $('#profileForm').submit(function(e) {
       e.preventDefault();
       e.returnValue = false;
  });
});

function getUserDetails() {
    console.log("Get User Data");

    var arrData = [];
    $.ajax({
        url: '/v1/ui/user?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;

          $("#username").val(arrData.Username);
          $("#email").val(arrData.Email);
          $("#phone").val(arrData.Phone);
        },
        error: function (data) {
          alert("User Not Found.");
        }
    });
}

function send() {
  var i;

  var error = false;
  $('input[check="true"]').each(function() {
    if($(this).val() == ""){
      $(this).addClass("error");
      $(this).parent().closest('div').addClass("input-error");
      error = true;
    }
  });

  if(error){
    alert("Please check your input.");
    return
  }

  var userReq = {
      username: $("#username").val(),
      email: $("#email").val(),
      phone: $("#phone").val(),
    };

    console.log(userReq);
    $.ajax({
       url: '/v1/ui/user/update?type=detail&token='+token,
       type: 'post',
       dataType: 'json',
       contentType: "application/json",
       data: JSON.stringify(userReq),
       success: function () {
           alert("User Updated.");
           window.location = "/user/profile?token="+token;
       }
   });
}
