$(window).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	$('#createEmailUser').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			name: {
				required: true
			},
			email: {
				required: true,
				email: true
			}
		}
	});
});

function send() {
	if(!$('#createEmailUser').valid()){
		$(".error").focus();
		return;
	}

	var userReq = {
		name: $("#name").val(),
		email: $("#email").val()
	};

	$.ajax({
		url: '/v1/ui/user/email/create?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(userReq),
		success: function () {
			swal({
					title: 'Success',
					text: 'Register Success',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/user/email/search";
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			if (a.errors.detail == "Duplicate Entry.") {
				swal("Mail already registered.");
			}
		}
	});
}
