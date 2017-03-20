$( window ).ready(function() {
   getAccount();
});

function getAccount(){
  $.ajax({
    url: 'http://voucher.apps.id:8889/v1/api/get/accountId',
    type: 'get',
    success: function (data) {
      var i = 0;
      for( i = 0; i < data.data.length; i ++){
        var option;
        if( i == 0){
          $('#select2-account-container').html(data.data[i].AccountName);
          option =$("<option selected='selected' value="+data.data[i].Id+">"+data.data[i].AccountName+"</li>")
        }
        else{
          option =$("<option value="+data.data[i].Id+">"+data.data[i].AccountName+"</li>")
        }

        option.appendTo('#account');
      }
    }
  });
}

$('#user-login').submit(function(e) {
     e.preventDefault();
     e.returnValue = false;
});

function login(){
  var request = {
      username: $("#username").val(),
      password: $("#password").val(),
      account_id: $("#account").find(":selected").val()
  };

  $.ajax({
      url: 'http://voucher.apps.id:8889/v1/api/login',
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(request),
      success: function (data){

        var token = data.data.Token;

        if (typeof(Storage) !== "undefined") {
          localStorage.setItem("token", token);
        }

        window.location = "http://voucher.apps.id:8889/variant/";
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

(function() {
    'use strict';

    $(userLogin);

    function userLogin() {

        var $form = $('#user-login');
        $form.validate({
            errorPlacement: errorPlacementInput
        });
    }

    // Necessary to place dyncamic error messages
    // without breaking the expected markup for custom input
    function errorPlacementInput(error, element) {
        if( element.parent().is('.mda-form-control') ) {
            error.insertAfter(element.parent()); // insert after .mda-form-control
        }
        else if ( element.is(':radio') || element.is(':checkbox')) {
            error.insertAfter(element.parent());
        }
        else {
            error.insertAfter(element);
        }
    }

})();
