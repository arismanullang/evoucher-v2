$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	var id = findGetParameter("id");
	getPartner(id);
});

function getPartner(id) {
	$.ajax({
		url: '/v1/ui/partner?id=' + id + "&token=" + token,
		type: 'get',
		success: function (data) {
			var arrData = data.data[0];
			$("#partnerName").html(arrData.name);
			$("#serialNumber").val(arrData.serial_number.String);
			$("#description").val(arrData.description.String);
		}
	});
}

function update() {
	var id = findGetParameter("id");

	var partner = {
		serial_number: $("#serialNumber").val(),
		description: $("#description").val()
	};

	$.ajax({
		url: '/v1/ui/partner/update?id=' + id + '&token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(partner),
		success: function () {
			swal({
					title: 'Success',
					text: 'Partner Updated',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					var id = findGetParameter("id");
					window.location = "/partner/search";
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
