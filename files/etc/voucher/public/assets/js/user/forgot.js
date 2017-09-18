$(window).ready(function () {
	$("#username").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});
	getAccount();
});

function recover() {

	var user = $("#username").val();
	var accountId = $("#account-id").val();
	if (user == null) {
		alert("Username cant be empty.");
	}

	$.ajax({
		url: '/v1/ui/user/forgot/mail?username=' + user + '&accountId=' + accountId,
		type: 'get',
		dataType: 'json',
		success: function (data) {
			swal({
					title: 'Success',
					text: 'Email Sent',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/user/mail-send";
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error " + a.errors.detail);
		}
	});
}

function getAccount() {
	$.ajax({
		url: '/v1/ui/account/all',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			for (i = 0; i < arrData.length; i++) {
				// <option value="bulk">Email Blast</option>
				var li = $("<option value='" + arrData[i].id + "'></option>");
				li.html(arrData[i].name);
				li.appendTo('#account-id');
			}

			$('.select2').select2();
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

		var $form = $('#user-recover');
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
