$( document ).ready(function() {
	$('#transaction').submit(function(e) {
		e.preventDefault();
		addElem();
		return false;
	});
});

function addElem(){
	var id = $('#transaction-code').val();
	$.ajax({
		url: '/v1/public/transaction/'+id,
		type: 'get',
		success: function (data) {
			console.log("Render Data");
			var result = data.data;
			$('#transaction-code').val('');

			var date = new Date(result.created_at);
			// var body = "";
			// for( i = 0; i < result.vouchers.length; i++){
			// 	body += result.vouchers[i].VoucherCode + "<br>";
			// }

			$("#label-transaction-code").html(result.transaction_code);
			// $("#voucher-code").html(body);
			$("#voucher-value").html("Rp. "+toDigit(result.discount_value.toString())+",00");
			$("#transaction-date").html(date.toDateString() + ", " +toTwoDigit(date.getHours()) + ":" + toTwoDigit(date.getMinutes()));
			$("#partner-name").html(result.partner_name);
			$("#member-name").html(result.vouchers[0].holder);
			$("#voucher-status").html(result.state);

			$("#error").html('');
		},
		error: function (data) {
			$('#transaction-code').val('');
			$("#voucher-code").html('');
			$("#voucher-value").html('');
			$("#transaction-date").html('');
			$("#partner-name").html('');
			$("#voucher-status").html('');
			$("#error").html(data.responseJSON.errors.detail);
		}
	});
}

function toDigit(param) {
	var result = "";
	var index = 0;
	var i = 0;
	for( i = param.length; i >=0; i--){
		result = param.charAt(i) + result;
		if(index % 3 == 0 && i != param.length && i != 0){
			result = '.' + result;
		}
		index++;
	}
	return result;
}

function toTwoDigit(val){
	if (val < 10){
		return '0'+val;
	}
	else {
		return val;
	}
}
