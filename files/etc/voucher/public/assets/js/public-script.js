$( document ).ready(function() {
	var transactioncode = findGetParameter('transactioncode');
	$('#transactioncode').html(transactioncode);
	if(transactioncode != null){
		getVoucherCode(transactioncode);
	}


	var x = findGetParameter('x')+"=";
	console.log(x);
	if(x != 'null='){
		getProfile(x);
	}

	$('#formsubmit').submit(function(e) {
		e.preventDefault();
		return false;
	});

	$("#tenant").change(function() {

		$.ajax({
			url: '/v1/public/challenge',
			type: 'get',
			success: function (data) {
				console.log(data);
				$('#challange-code').html(data.data.challenge);
			}
		});
	});
});

function getProfile(x){
	$.ajax({
		url: '/v1/public/redeem/profile?x='+x,
		type: 'get',
		success: function (data) {
			console.log(data.data.Vouchers[0].voucher_id);
			$("#holdername").html(data.data.holder);
			$("#program-id").val(data.data.program_id);
			$("#discount-value").val(data.data.discount_value);
			$("#voucher").val(data.data.Vouchers[0].voucher_id);
			$("#tnc").html(data.data.program_tnc);
		}
	});
}

function getVoucherCode(x){
	$.ajax({
		url: '/v1/public/transaction/'+encodeURIComponent(x),
		type: 'get',
		success: function (data) {
			var result = data.data;
			console.log(result.vouchers[0].VoucherCode);
			for(i = 0; i < result.vouchers.length; i++){
				var li = $("<strong></strong>");
				li.html(result.vouchers[i].VoucherCode);
				li.appendTo('#vouchers');
			}
		}
	});
}

function findGetParameter(parameterName) {
	var result = null,
		tmp = [];
	location.search
		.substr(1)
		.split("&")
		.forEach(function (item) {
			tmp = item.split("=");
			if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
		});
	return result;
}

function send(){
	var program = {
		program_id:$("#program-id").val(),
		redeem_method:"token",
		partner:$("#tenant").val(),
		challenge:$("#challange-code").html(),
		response:$("#response-code").val(),
		discount_value:parseInt($("#discount-value").val()),
		vouchers:[$("#voucher").val()]
	};
	console.log(program);
	$.ajax({
		url: '/v1/public/transaction',
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(program),
		success: function (data) {
			console.log(data);
			window.location = '/public/success?transactioncode='+data.data.transaction_code;
		},
		error: function (data) {
			console.log(data);
			window.location = '/public/fail';
		}
	});
}


(function() {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$('.select2').select2();
	}
})();
