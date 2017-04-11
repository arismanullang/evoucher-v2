$( window ).ready(function() {

});

$('#user-login').submit(function(e) {
     e.preventDefault();
     e.returnValue = false;
});

function login(){
  $.ajax({
      url: '/v1/token',
      type: 'get',
      dataType: 'json',
      contentType: "application/json",
      beforeSend: function (xhr) {
          xhr.setRequestHeader ("Authorization", "Basic " + btoa($("#username").val() + ":" + $("#password").val()));
      },
      success: function (data){

        var token = data.data.token;
        console.log(token);
        if (typeof(Storage) !== "undefined") {
          localStorage.setItem("token", token);
        }

        window.location = "/variant/";
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
