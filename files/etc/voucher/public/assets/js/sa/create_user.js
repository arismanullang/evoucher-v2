$(window).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	getAccount();
	$("#account").change(function (){
		getRole(this.value);
	});

	$('#createUser').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			username: {
				required: true
			},
			password: {
				required: true
			},
			email: {
				required: true,
				email: true
			},
			phone: {
				required: true
			},
			'role[]': {
				required: true
			}
		}
	});
});

function send() {
	if(!$("#createUser").valid()){
		return;
	}

	var i;

	var listRole = [];
	var li = $("input[type=checkbox]:checked");

	for (i = 0; i < li.length; i++) {
		listRole[i] = li[i].value;
	}

	var userReq = {
		account_id: $("#account").val(),
		username: $("#username").val(),
		password: $("#password").val(),
		email: $("#email").val(),
		phone: $("#phone").val(),
		role_id: listRole
	};

	$.ajax({
		url: '/v1/ui/sa/create?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(userReq),
		success: function () {
			swal({
					title: 'Success',
					text: 'User Created',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/sa/search";
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			if (a.errors.detail == "Duplicate Entry.") {
				swal("Username already used.");
			}
		}
	});
}

function getRole(account) {
	$('#role').html("");
	$.ajax({
		url: '/v1/ui/role/account?token='+token+'&id='+account,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<div class='col-sm-4'></div>");
				var html = "<label class='checkbox-inline c-checkbox'>"
					+ "<input type='checkbox' name='role[]' value='" + arrData[i].id + "' text='" + arrData[i].detail + "'>"
					+ "<span class='ion-checkmark-round'></span>" + arrData[i].detail
					+ "</label>";
				li.html(html);
				li.appendTo('#role');
			}
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
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
				li.appendTo('#account');
			}

			getRole(arrData[0].id);
			$('.select2').select2();
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
