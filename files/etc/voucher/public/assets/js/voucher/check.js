$(document).ready(function () {
	getVoucher();
});

function getVoucher() {
	var id = findGetParameter('id');
	$.ajax({
		url: '/v1/ui/voucher/' + id + '?token=' + token,
		type: 'get',
		success: function (data) {
			var data = data.data;

			var date1 = data.valid_at.substring(0, 19).replace("T", " ");
			var date2 = data.expired_at.substring(0, 19).replace("T", " ");

			$("#program-name").html(data.program_name);
			$("#voucher-code").html(data.voucher_code);
			$("#voucher-type").html(data.voucher_type);
			$("#voucher-value").html("Rp " + addDecimalPoints(data.voucher_value) + ",00");
			$("#reference-no").html(data.reference_no);
			$("#period").html(date1 + "</br></br>To</br></br>" + date2)

			var email = data.holder_email;
			if (data.holder_email == "") {
				email = "Unknown";
			}
			var phone = data.holder_phone;
			if (data.holder_phone == "") {
				phone = "Unknown";
			}

			$("#holder-name").html(toTitleCase(data.holder_description));
			$("#holder-email").html(email);
			$("#holder-phone").html(phone);

			var dateCreated = data.created_at.substring(0, 19).replace("T", " ");
			$("#issued-state").html(dateCreated);
		}
	});
}
