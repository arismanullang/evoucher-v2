$( window ).load(function() {
   getAccount();
});
function getAccount(){
  $.ajax({
    url: 'http://evoucher.elys.id:8889/get/accountId',
    type: 'get',
    success: function (data) {
      alert(data.data[0].Id);
    }
  });
}

function login(){
  var request = {
      username: $("#username").val(),
      password: $("#password").val()
  };

  $.ajax({
      url: 'http://evoucher.elys.id:8889/login',
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(request),
      success: function (data){
        alert(data.data);
        window.location = "http://evoucher.elys.id:8889/variant/create?token="+data.data;
      }
  });
}

(function() {
    'use strict';

    $(formAdvanced);

    function formAdvanced() {
        $('.select2').select2();
    }

})();
