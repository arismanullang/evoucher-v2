$(window).ready(function () {
	var token = localStorage.getItem("token");
	if (token != null) {
		$.ajax({
			url: '/v1/ui/token/check?token=' + token + '&url=' + window.location.pathname,
			type: 'get',
			success: function (data) {
				var result = data.data;
				window.location = result.destination;
			}
		});
	}

	$("#password").on('keyup', function (e) {
		if (e.keyCode == 13) {
			login();
		}
	});
});

function login() {
	$.ajax({
		url: '/v1/ui/user/login',
		type: 'get',
		dataType: 'json',
		contentType: "application/json",
		beforeSend: function (xhr) {
			xhr.setRequestHeader("Authorization", "Basic " + btoa($("#username").val() + ":" + $("#password").val()));
		},
		success: function (data) {
			var result = data.data;
			var token = result.token.token;
			var role = result.role;
			if (typeof(Storage) !== "undefined") {
				localStorage.setItem("token", token);
				tempRole = "";
				for (var i = 0; i < role.length; i++) {
					tempRole += role[i].id + ",";
				}
				localStorage.setItem("r", tempRole);
				localStorage.setItem("ui", result.ui);
			}

			window.location = result.destination;
		},
		error: function (data) {
			alert("Invalid Username or Password");
		}
	});
}

(function () {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$('.select2').select2();
	}

})();

(function () {
	'use strict';

	$(userLogin);

	function userLogin() {
		var $form = $('#userLogin');
		$form.validate({
			errorPlacement: errorPlacementInput
		});
	}

	// Necessary to place dyncamic error messages
	// without breaking the expected markup for custom input
	function errorPlacementInput(error, element) {
		if (element.parent().is('.mda-form-control')) {
			error.insertAfter(element.parent()); // insert after .mda-form-control
		}
		else if (element.is(':radio') || element.is(':checkbox')) {
			error.insertAfter(element.parent());
		}
		else {
			error.insertAfter(element);
		}
	}

})();
