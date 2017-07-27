$( window ).ready(function() {

});

$('#user-login').submit(function(e) {
     e.preventDefault();
     e.returnValue = false;
});

function login(){
  $.ajax({
      url: '/v1/ui/user/login',
      type: 'get',
      dataType: 'json',
      contentType: "application/json",
      beforeSend: function (xhr) {
          xhr.setRequestHeader ("Authorization", "Basic " + btoa($("#username").val() + ":" + $("#password").val()));
      },
      success: function (data){

        var token = data.data.token.token;
        var role = data.data.role;
        console.log(token);
        if (typeof(Storage) !== "undefined") {
          localStorage.setItem("token", token);
          tempRole = "";
          for(var i = 0; i < role.length; i++){
          	tempRole += role[i].id+",";
	  }
          localStorage.setItem("r", tempRole);
        }

        window.location = "/program/index?token="+token;
      },
      error: function (data){
        alert("Invalid Username or Password");
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
