$( window ).ready(function() {
	$('#user-login').submit(function(e) {
		e.preventDefault();
		e.returnValue = false;
	});

	var token = localStorage.getItem("token");
	if(token != null){
		$.ajax({
			url: '/v1/token/check?token='+token+'&url='+window.location.pathname,
			type: 'get',
			success: function (data) {
				if(data.data == true){
					var role = data.data.role;
					if(role[0].id != 'Mn78I1wc'){
						window.location = "/program/index";
					} else{
						window.location = "/sa/search";
					}
				}
			}
		});
	}
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

        if(role[0].id != 'Mn78I1wc'){
		window.location = "/program/index";
	} else{
        	window.location = "/sa/search";
	}
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
