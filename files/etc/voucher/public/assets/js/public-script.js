$( document ).ready(function() {
	var transactioncode = findGetParameter('transactioncode');
	$('#transactioncode').html(transactioncode);

	var x = findGetParameter('x')+"=";
	console.log(x);
	getProfile(x);

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
			$("#variant-id").val(data.data.variant_id);
			$("#discount-value").val(data.data.discount_value);
			$("#voucher").val(data.data.Vouchers[0].voucher_id);
		}
	});
}

function send(){
	var variant = {
		variant_id:$("#variant-id").val(),
		redeem_method:"token",
		partner:$("#tenant").val(),
		challenge:$("#challange-code").html(),
		response:$("#response-code").val(),
		discount_value:parseInt($("#discount-value").val()),
		vouchers:[$("#voucher").val()]
	};
	console.log(variant);
	$.ajax({
		url: '/v1/public/transaction',
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(variant),
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
