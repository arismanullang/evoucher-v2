$( window ).load(function() {
   getAccount();
});

function getAccount(){
  $.ajax({
    url: 'http://evoucher.elys.id:8889/get/accountId',
    type: 'get',
    success: function (data) {
      var i = 0;
      for( i = 0; i < data.data.Data.length; i ++){
        var option;
        if( i == 0){
          $('#select2-account-container').html(data.data.Data[i].AccountName);
          option =$("<option selected='selected' value="+data.data.Data[i].Id+">"+data.data.Data[i].AccountName+"</li>")
        }
        else{
          option =$("<option value="+data.data.Data[i].Id+">"+data.data.Data[i].AccountName+"</li>")
        }

        option.appendTo('#account');
      }
    }
  });
}

$('#user-login').submit(function(e) {

     // this code prevents form from actually being submitted
     e.preventDefault();
     e.returnValue = false;

     login();
});

function login(){
  var request = {
      username: $("#username").val(),
      password: $("#password").val(),
      account_id: $("#account").find(":selected").val()
  };

  $.ajax({
      url: 'http://evoucher.elys.id:8889/login',
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(request),
      success: function (data){
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
