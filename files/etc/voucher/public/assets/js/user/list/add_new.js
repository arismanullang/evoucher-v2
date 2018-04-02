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

	var id = findGetParameter("id");
	$("#list-id").val(id);
	getUser(id);
});

function getUser(id) {
	$.ajax({
		url: '/v1/ui/user/list?token=' + token + '&id=' + id,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			$("#list-name").html(arrData.name);
		}
	});
}


function send() {
	if(!$('#createEmailUser').valid()){
		$(".error").focus();
		return;
	}

	var userReq = {
		name: $("#name").val(),
		email: $("#email").val(),
		list_id: $("#list-id").val()
	};

	$.ajax({
		url: '/v1/ui/user/list/add-new?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(userReq),
		success: function () {
			swal({
					title: 'Success',
					text: 'Submit Success',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/user/list/check?id="+$("#list-id").val();
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			if (a.errors.detail.toLowerCase() == "duplicate entry.") {
				swal("Mail already registered.");
			}
		}
	});
}
