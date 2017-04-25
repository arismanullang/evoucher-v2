$( document ).ready(function() {
  getUserDetails();
  getAccount();
  getUser();
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
        url: '/v1/api/get/userDetails?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.RoleId.length;
          var desc = "Act as ";
          for ( i = 0; i < limit; i++){
            desc += arrData.RoleId[i];
            if( i != limit-1){
              desc += ", ";
            }
          }
          desc += ".";
          $("#user-desc").html(desc);
          $("#user-name").html(arrData.Username);
          $("#user-email").html(arrData.Email);
          $("#user-phone").html(arrData.Phone);
          $("#user-date").html(arrData.CreatedAt.substr(0,10));
        },
        error: function (data) {
          alert("User Not Found.");
        }
    });
}

function getUser() {
    console.log("Get Voucher Data");

    var arrData = [];
    $.ajax({
        url: '/v1/api/get/users?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.length;

          for ( i = 0; i < limit; i++){
            var html = "<img src='/assets/img/user/04.jpg' alt='List user' class='mda-list-item-img'>"
              + "<div class='mda-list-item-text mda-2-line'>"
              +    "<h3>"+arrData[i].Username+"</h3>"
              + "</div>";
            var li = $("<div class='mda-list-item'></div>").html(html);
            li.appendTo('#listUser');
          }
        },
        error: function (data) {
          alert("Teammates Not Found.");
        }
    });
}

function getAccount() {
    console.log("Get Account Data");

    $.ajax({
        url: '/v1/api/get/accountsDetail?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;

          var limit = result.length;
          var desc = "";
          for ( i = 0; i < limit; i++){
            desc += result[i].AccountName;
            if( i != limit-1){
              desc += ", ";
            }
          }

          $("#user-accounts").html(desc.toUpperCase());
        },
        error: function (data) {
          alert("Account Not Found.");
        }
    });
}

function getVariant() {
    console.log("Get Account Data");

    $.ajax({
        url: '/v1/api/get/totalVariant?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;
          $("#user-variant").html(result);
        },
        error: function (data) {
          alert("Account Not Found.");
        }
    });
}

function updateUser(){
  window.location = "/user/update";
}
