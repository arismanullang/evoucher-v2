$( document ).ready(function() {
	$('#transaction').submit(function(e) {
		e.preventDefault();
		addElem();
		return false;
	});
	$('.select2').select2();
	getPartner();
	$('#partnerList').change(function () {
		getTransactionByPartner(this.value);
	});
});

function getPartner() {
	$.ajax({
		url: '/v1/ui/partner/all?token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;
			for(var i = 0; i < result.length; i++){
				var li = $("<option value='"+result[i].id+"'>"+result[i].name+"</td>");
				li.appendTo('#partnerList');
			}
		},
		error: function (data) {
		}
	});
}

function getTransactionByPartner(partnerId) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/partner?token=' + token + '&partner=' + partnerId,
		type: 'get',
		success: function (data) {
			var result = data.data;
			for(var i = 0; i < result.length; i++){
				if(result[i].state == 'used')
					arrData[i] = result[i];
			}

			for(var i = 0; i < result.length; i++){
				var date = new Date(result[i].issued);
				var body = "<td class='col-lg-1 checkbox c-checkbox'><label>"
					+ "<input type='checkbox' name='transaction-code' value='"+result[i].transaction_code+"'><span class='ion-checkmark-round'></span>"
					+ "</label></td>"
					+ "<td class='text-ellipsis'>"+result[i].transaction_code+"</td>"
					+ "<td class='text-ellipsis'>"+result[i].voucher_value+"</td>"
					+ "<td class='text-ellipsis'>"+date.toDateString() + ", " + date.getHours() + ":" + date.getMinutes()+"</td>"
				var li = $("<tr></tr>");
				li.html(body);
				li.appendTo('#listTransaction');

			}
		}
	});
}

function cashout(){
	var transactionCode = [];
	var transactions = "";
	var elem = $("*[name='list-transaction-code']");
	var i = 0;
	for ( i = 0; i < elem.length; i++){
		transactionCode[i] = elem[i].innerHTML;
		transactions += elem[i].innerHTML + ';';
	}

	var transaction = {
		transaction_code: transactionCode
	};

	var decoded = $("*[name='list-transaction-partner']")[0];

	var textArea = document.createElement('textarea');
	textArea.innerHTML = decoded.innerHTML;

	$.ajax({
		url: '/v1/ui/transaction/cashout/update?token='+token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(transaction),
		success: function () {
			localStorage.setItem("transaction_cashout", "");
			localStorage.setItem("transaction_cashout", transactions);
			$("#success-page").attr("style","display:block");
			$("#transaction").attr("style","display:none");
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function addElem(){
	var id = $('#transaction-code').val();
	$.ajax({
		url: '/v1/ui/transaction?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;
			$('#transaction-code').val('');
			var elem = $("*[name='list-transaction-code']");
			var i = 0;
			for ( i = 0; i < elem.length; i++){
				if(id == elem[i].innerHTML){
					$("#error").html('Voucher Already Added');
					return
				}
			}

			if(result.state == "paid"){
				$("#error").html('Voucher Already Used');
			}else{
				var date = new Date(result.created_at);
				for(var i = 0; i < result.vouchers.length; i++){
					var body = "<td name='list-transaction-partner' class='text-ellipsis'>"+result.partner_name+"</td>"
						+ "<td name='list-transaction-code' class='text-ellipsis'>"+result.transaction_code+"</td>"
						+ "<td name='list-voucher-code' class='text-ellipsis'>"+result.vouchers[i].voucher_code+"</td>"
						+ "<td name='list-transaction-value' class='text-ellipsis'>"+result.discount_value+"</td>"
						+ "<td name='list-transaction-date' class='text-ellipsis'>"+date.toDateString() + ", " + date.getHours() + ":" + date.getMinutes()+"</td>"
						+ "<td name='list-transaction-action'><button type='button' onclick='removeElem(this)' class='btn btn-flat btn-sm btn-info pull-right'><em class='ion-close-circled'></em></button></td>";
					var li = $("<tr class='msg-display clickable'></tr>");
					li.html(body);
					li.appendTo('#list-transaction');
				}
				$("#error").html('');

			}
		},
		error: function (data) {
			$('#transaction-code').val('');
			$("#error").html(data.responseJSON.errors.title);
		}
	});
}

function removeElem(elem){
	$(elem).parent().closest('tr').remove();
}

function print(){
	window.location = "/voucher/print";
}

function next(){
	swal({
			title: 'Are you already print the invoice?',
			text: 'You will not be able to recover the last details',
			type: 'warning',
			showCancelButton: true,
			confirmButtonColor: '#4CAF50',
			confirmButtonText: 'Yes',
			closeOnConfirm: true
		},
		function() {
			window.location.reload();
		});
}
