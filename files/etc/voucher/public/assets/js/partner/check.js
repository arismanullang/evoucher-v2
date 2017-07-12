$( document ).ready(function() {
	var id = findGetParameter('id');
	getPartner(id);
	getPerformance(id);
});

function getPartner(id) {
	console.log("Get Partner Data");

	$.ajax({
		url: '/v1/ui/partner?id='+id+"&token="+token,
		type: 'get',
		success: function (data) {
			console.log(data.data);
			var arrData = data.data[0];
			$("#initial").html(arrData.partner_name.charAt(0).toUpperCase());
			$("#partner-title").html(arrData.partner_name.toUpperCase());
			$("#partner-name").html(arrData.partner_name);
			$("#serial-number").html(arrData.serial_number.String);
			$("#tag").html(arrData.tag.String);
			$("#desciption").html(arrData.description.String);
		}
	});
}

function getPerformance(id) {
	console.log("Get Voucher Data");

	var arrData = [];
	$.ajax({
		url: '/v1/ui/partner/performance?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			console.log(data.data);
			var i;
			var arrData = data.data;
			$("#total-program").html(arrData.program);
			$("#total-transaction").html(arrData.transaction_code);
			$("#total-voucher-value").html("Rp. " + addDecimalPoints(arrData.transaction_value) + ",00");
			$("#total-valid-voucher").html(arrData.voucher_generated);
			$("#total-used-voucher").html(arrData.voucher_used);
			$("#total-customer").html(arrData.customer);
		}
	});
}

function getVariant() {
	console.log("Get Account Data");

	$.ajax({
		url: '/v1/ui/variant/all?token='+token,
		type: 'get',
		success: function (data) {
			console.log(data.data);
			var result = data.data;
			$("#user-variant").html(result.length);
		},
		error: function (data) {
			alert("Account Not Found.");
		}
	});
}

function updateUser(){
	window.location = "/user/update?token="+token;
}
