$(document).ready(function () {
	var id = findGetParameter('id');

	getVoucher(id);
	getTransactionVoucher(id);
});

function getVoucher(id) {
	$.ajax({
		url: '/v1/ui/voucher/' + id + '?token=' + token,
		type: 'get',
		success: function (data) {
			var data = data.data;

			var date1 = new Date(data.valid_at).toString().substring(0, 24);
			var date2 = new Date(data.expired_at).toString().substring(0, 24);

			$("#program-name").html(data.program_name);
			$("#voucher-code").html(data.voucher_code);
			$("#voucher-value").html("Rp " + addDecimalPoints(data.voucher_value));
			$("#reference-no").html(data.reference_no);
			$("#period").html("</br>" + date1 + "</br></br>" + date2)

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

			var dateCreated = new Date(data.created_at).toString().substring(0, 24);
			$("#issued-state").html(dateCreated);
		}
	});
}

function getTransactionVoucher(voucherId) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/voucher?token=' + token + '&id=' + voucherId,
		type: 'get',
		success: function (data) {
			var result = data.data;
			console.log(result.voucher.state);
			if(result.voucher.state == 'used' || result.voucher.state == 'paid'){
				$('#redeem-icon').removeClass('ion-android-close');
				$('#redeem-icon').addClass('ion-android-done');

				var dateRedeem = new Date(result.redeemed).toString().substring(0, 24);
				$('#redeem-state').html(result.transaction_code + " || " + dateRedeem);
			}

			if(result.voucher.state == 'paid'){
				$('#cashout-icon').removeClass('ion-android-close');
				$('#cashout-icon').addClass('ion-android-done');

				var dateCashout = new Date(result.cashout.String).toString().substring(0, 24);
				$('#cashed-state').html(dateCashout);
			}
		}
	});
}
