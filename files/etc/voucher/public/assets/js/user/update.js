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
						if (tempElem.getAttribute("text") == arrData.role[y].role_detail) {
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
		url: '/v1/ui/role/all',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<div class='col-sm-4'></div>");
				var html = "<label class='checkbox-inline c-checkbox'>"
					+ "<input type='checkbox' value='" + arrData[i].id + "' text='" + arrData[i].detail + "'>"
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
	var i;

	var listRole = [];
	var li = $("input[type=checkbox]:checked");

	if (li.length == 0 || parseInt($("#length").val()) < 8) {
		error = true;
	}

	for (i = 0; i < li.length; i++) {
		listRole[i] = li[i].value;
	}

	var error = false;
	$('input[check="true"]').each(function () {
		if ($(this).val() == "") {
			$(this).addClass("error");
			$(this).parent().closest('div').addClass("input-error");
			error = true;
		}
	});

	if (error) {
		swal("Please check your input.");
		return
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
			swal("User Updated.");
			window.location = "/user/search";
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
