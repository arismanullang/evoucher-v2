$( document ).ready(function() {
	var x = findGetParameter('x')+"=";
	console.log(x);
	getProfile(x);

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
			console.log(data);
			$("#holdername").html(data.data.holder);
			$("#variant-id").html(data.data.variant_id);
			$("#discount-value").val(data.data.discount_value);
			$("#voucher").val(data.data.Vouchers[0]);
		}
	});
}

function send(){
	var variant = {
		variant_id:$("#variant-id").val(),
		redeem_method:"token",
		partner:$("#partner-id").val(),
		challenge:$("#challange-code").html(),
		response:$("#response-code").html(),
		discount_value:$("#discount-value").val(),
		vouchers:[$("voucher").val()]
	};
	$.ajax({
		url: '/v1/public/transaction',
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(variant),
		success: function (data) {
			console.log(data);
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
