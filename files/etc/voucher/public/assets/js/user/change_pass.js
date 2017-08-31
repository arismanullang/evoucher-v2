$(window).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});
});

function recover() {
	if ($("#newpass").val() != $("#confpass").val()) {
		swal("Passwords are not match.");
		return;
	}

	var user = {
		old_password: $("#oldpass").val(),
		new_password: $("#newpass").val()
	};

	$.ajax({
		url: '/v1/ui/user/update?type=password&token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			swal("Password Updated");
			window.location = "/user/profile?token=" + token;
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
