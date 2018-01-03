var qr = "";

$( document ).ready(function() {
	var x = findGetParameter('x')+"=";
	if(x != 'null='){
		getProfile(x);
	}

	var transactioncode = localStorage.getItem('public_transaction_code');
	$('#transactioncode').html(transactioncode);
	if(transactioncode != null){
		getVoucherCode(transactioncode);
	}

	var error = localStorage.getItem('public_error_message');
	$('#error').html(error);

	$('#formsubmit').submit(function(e) {
		e.preventDefault();
		return false;
	});

	// $("#tenant").change(function() {
	// 	$.ajax({
	// 		url: '/v1/public/challenge',
	// 		type: 'get',
	// 		success: function (data) {
	// 			console.log(data);
	// 			$('#challange-code').html(data.data.challenge);
	// 		}
	// 	});
	// });
});

function handleFiles(f){
	var o=[];
	for(var i =0;i<f.length;i++){
		var reader = new FileReader();
		reader.onload = (function(theFile) {
			return function(e){
				qrcode.callback = read;
				qrcode.decode(e.target.result);
			};
		})(f[i]);

		// Read in the image file as a data URL.
		reader.readAsDataURL(f[i]);	}
}

function read(a){
	qr = a;
	var message = "Scan QR-code success.";
	if(a.includes("error")){
		alert(a);
		message = "Scan QR-code error.";
	}

	document.getElementById("message").innerHTML = message;
	console.log(qr);
}

function picChange(evt){
	//get files captured through input
	var fileInput = evt.target.files;
	if(fileInput.length>0){
		handleFiles(fileInput);
	}
}

function getProfile(x){
	$.ajax({
		url: '/v1/public/redeem/profile?x='+x,
		type: 'get',
		success: function (data) {
			console.log(data.data.vouchers[0].voucher_id);
			$("#holdername").html(data.data.holder);
			$("#programId").val(data.data.program_id);
			$("#discountValue").val(data.data.voucher_value);
			$("#voucher").val(data.data.vouchers[0].voucher_id);
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
			console.log(result.vouchers[0].voucher_code);
			for(i = 0; i < result.vouchers.length; i++){
				var li = $("<strong></strong>");
				li.html(result.vouchers[i].voucher_code);
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
		program_id:$("#programId").val(),
		redeem_method:"qr",
		partner: qr,
		discount_value:$("#discountValue").val(),
		vouchers:[$("#voucher").val()]
	};

	$.ajax({
		url: '/v1/public/transaction',
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(program),
		success: function (data) {
			console.log(data);
			localStorage.setItem("public_transaction_code", "");
			localStorage.setItem("public_transaction_code", data.data.transaction_code);
			window.location = '/public/success';
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			console.log(a.errors.detail);
			localStorage.setItem("public_error_message", "");
			localStorage.setItem("public_error_message", a.errors.detail);

//			window.location = '/public/fail';
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
