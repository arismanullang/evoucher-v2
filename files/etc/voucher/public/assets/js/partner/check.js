$(document).ready(function () {
	var id = findGetParameter('id');
	getPartner(id);
	getPerformance(id);
	makeCode(id);
});


function makeCode(id) {
	var qrcode = new QRCode("qrcode", {
		text: id,
		width: 150,
		height: 150,
		colorDark: "#000000",
		colorLight: "#ffffff",
		correctLevel: QRCode.CorrectLevel.H
	});
	qrcode.makeCode(id);
}

function getPartner(id) {
	$.ajax({
		url: '/v1/ui/partner?id=' + id + "&token=" + token,
		type: 'get',
		success: function (data) {
			var arrData = data.data[0];
			$("#initial").html(arrData.name.charAt(0).toUpperCase());
			$("#partner-title").html(arrData.name.toUpperCase());
			$("#partner-name").html(arrData.name);
			$("#serial-number").html(arrData.serial_number.String);
			$("#tag").html(arrData.tag.String);
			$("#desciption").html(arrData.description.String);
		}
	});
}

function getPerformance(id) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/partner/performance?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
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

function edit() {
	window.location = "/partner/update?id=" + findGetParameter('id');
}
