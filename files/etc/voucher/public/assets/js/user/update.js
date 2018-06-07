$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	var id = findGetParameter("id");
	var type = "detail";
	if (id != null) {
		type = "other";
		getRole();
	}
	$("#id").val(id);
	$("#type").val(type);
	getUserDetails(id, type);

	$('#updateUser').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
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

function getUserDetails(id, type) {
	var url = '/v1/ui/user?token=' + token;
	if (id != null) {
		url = '/v1/ui/user/other?id=' + id + '&token=' + token;
	}
	var arrData = [];
	$.ajax({
		url: url,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;

			$("#username").html(arrData.username);
			$("#email").val(arrData.email);
			$("#phone").val(arrData.phone);

			if (type == "other") {
				var li = $("input[type=checkbox]");

				for (i = 0; i < li.length; i++) {
					var tempElem = li[i];
					var limit = arrData.role.length;
					for (y = 0; y < limit; y++) {
						if (tempElem.value == arrData.role[y].id) {
							tempElem.checked = true;
						}
					}
				}
			} else {
				$("#row-role").attr("style", "display:none");
			}
		},
		error: function (data) {
			swal("User Not Found.");
		}
	});
}

function getRole() {
	$.ajax({
		url: '/v1/ui/role/all?token='+token,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<div class='col-sm-4 checkbox-add-padding'></div>");
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

function send() {
	if(!$("#updateUser").valid()){
		return;
	}

	var i;

	var listRole = [];
	var li = $("input[type=checkbox]:checked");

	for (i = 0; i < li.length; i++) {
		listRole[i] = li[i].value;
	}

	var userReq = {
		id: $("#id").val(),
		username: $("#username").html(),
		email: $("#email").val(),
		phone: $("#phone").val(),
		role_id: listRole,
	};

	var type = $("#type").val();
	var url = '/v1/ui/user/update?type=detail&token=' + token;
	if (type == "other") {
		url = '/v1/ui/user/update?type=other&token=' + token;
	}

	$.ajax({
		url: url,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(userReq),
		success: function () {
			swal({
					title: 'Success',
					text: 'User Updated',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					if(type == "other"){
						window.location = "/sa/search";
					} else{
						window.location = "/user/search";
					}
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
