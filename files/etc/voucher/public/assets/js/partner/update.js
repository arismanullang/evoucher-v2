$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	$('.select2').select2({
		tags: true
	});
	getTag();
	var id = findGetParameter("id");
	getPartner(id);

	$('#update-partner').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			'partner-name': {
				required: true
			},
			'serial-number':{
				digits: true
			},
			'email': {
				required: true
			}
		}
	});
});

function getTag() {
	$.ajax({
		url: '/v1/ui/tag/all',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<option></option>").html(arrData[i]);
				li.appendTo('#tags');
			}
		}
	});
}

function getPartner(id) {
	$.ajax({
		url: '/v1/ui/partner?id=' + id + "&token=" + token,
		type: 'get',
		success: function (data) {
			var arrData = data.data[0];
			$("#partner-name").val(arrData.name);
			$("#serial-number").val(arrData.serial_number.String);
			$("#email").val(arrData.email);
			$("#description").val(arrData.description.String);
			$("#building").val(arrData.building);
			$("#city").val(arrData.city);
			$("#province").val(arrData.province);
			$("#zip-code").val(arrData.zip_code);
			$("#address").val(arrData.address);
			$("#bank-name").val(arrData.bank_name);
			$("#bank-branch").val(arrData.bank_branch);
			$("#bank-account-holder").val(arrData.bank_account_holder);
			$("#bank-account-number").val(arrData.bank_account_number);
			$("#company-email").val(arrData.company_email);
			$("#company-telp").val(arrData.company_telp);
			$("#company-name").val(arrData.company_name);
			$("#company-pic").val(arrData.company_pic);
			var tags = arrData.tag.String.split('#');
			$("#tags").select2().val(tags);
			$("#tags").select2().trigger('change');
		}
	});
}

function update() {
	if(!$("#update-partner").valid()){
		$(".error").focus();
		return;
	}

	var id = findGetParameter("id");

	var listTag = "";
	var li = $("ul.select2-selection__rendered").find("li");
	if (li.length > 0) {
		for (i = 0; i < li.length - 1; i++) {
			var text = li[i].getAttribute("title");

			listTag = listTag + "#" + text;
		}
	}

	var partner = {
		name: $("#partner-name").val(),
		serial_number: $("#serial-number").val(),
		email: $("#email").val(),
		tag: listTag,
		description: $("#description").val(),
		bank_account: $("#bank-accounts").find(":selected").val(),
		address: $("#address").val(),
		city: $("#city").val(),
		province: $("#province").val(),
		building: $("#building").val(),
		zip_code: $("#zip-code").val(),
		bank_name: $("#bank-name").val(),
		bank_branch: $("#bank-branch").val(),
		bank_account_holder: $("#bank-account-holder").val(),
		bank_account_number: $("#bank-account-number").val(),
		company_email: $("#company-email").val(),
		company_telp: $("#company-telp").val(),
		company_name: $("#company-name").val(),
		company_pic: $("#company-pic").val()
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
					window.location = "/partner/check?id="+id;
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
