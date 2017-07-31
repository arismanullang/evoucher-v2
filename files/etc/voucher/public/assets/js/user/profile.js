$( document ).ready(function() {
  getUserDetails();
  //getUser();
  getVariant();

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
          var result = data.data;
          var limit = result.role.length;
          var desc = "Act as ";
          for ( i = 0; i < limit; i++){
            desc += result.role[i].role_detail;
            if( i != limit-1){
              desc += ", ";
            }
          }
          desc += ".";
	  var date = new Date(result.created_at);

	  $("#user-accounts").html(result.account.name);
          $("#user-desc").html(desc);
          $("#user-name").html(result.username);
          $("#user-email").html(result.email);
          $("#user-phone").html(result.phone);
          $("#user-date").html(date.toDateString() + ", " + toTwoDigit(date.getHours()) + ":" + toTwoDigit(date.getMinutes()));
        },
        error: function (data) {
          alert("User Not Found.");
        }
    });
}

function getVariant() {
    console.log("Get Account Data");

    $.ajax({
        url: '/v1/ui/program/all?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;
          $("#user-program").html(result.length);
        },
        error: function (data) {
          alert("Account Not Found.");
        }
    });
}

function updateUser(){
  window.location = "/user/update?token="+token;
}
