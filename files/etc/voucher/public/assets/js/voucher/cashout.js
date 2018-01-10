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

	$('#transaction').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			bankAccount: {
				required: true,
				minlength: 10,
				maxlength: 20
			}
		}
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
			var voucher = 0;
			$("#listTransaction").html(" ");
			for(var i = 0; i < result.length; i++){
				if(result[i].vouchers[0].state == 'used')
					arrData.push(result[i]);
			}

			for(var i = 0; i < arrData.length; i++){
				var date = new Date(arrData[i].issued);
				for(var j = 0; j < arrData[i].vouchers.length; j++){
					var body = "<td class='col-lg-1 checkbox c-checkbox'><label>"
						+ "<input type='checkbox' name='transaction' class='transaction' value='"+arrData[i].transaction_id+";"+arrData[i].voucher_value+";"+arrData[i].vouchers[j].id+"'><span class='ion-checkmark-round'></span>"
						+ "</label></td>"
						+ "<td class='text-ellipsis'>"+arrData[i].transaction_code+"</td>"
						+ "<td class='text-ellipsis'>"+arrData[i].vouchers[j].voucher_code+"</td>"
						+ "<td class='text-ellipsis'> Rp. "+addDecimalPoints(arrData[i].voucher_value)+",00</td>"
						+ "<td class='text-ellipsis'>"+date.toDateString() + ", " + date.getHours() + ":" + date.getMinutes()+"</td>"
					var li = $("<tr></tr>");
					li.html(body);
					li.appendTo('#listTransaction');
					voucher++;
				}
			}

			if(voucher > 5){
				$("#tableTransaction").attr("style","overflow:scroll; max-height: 300px;");
			}else{
				$("#tableTransaction").removeAttribute("style");
			}
		},
		error: function(){
			$("#listTransaction").html(" ");
		}
	});
}

function cashout(){

	if(!$("#transaction").valid()) {
		return
	}

	var listVoucher = [];
	var listVoucherValue = [];
	var listTransaction = [];
	var li = $("input[class=transaction]:checked");
	var total = 0;

	for (i = 0; i < li.length; i++) {
		console.log(li[i].value);
		if (li[i].value != "on") {
			var tempValue = li[i].value.split(";");

			listTransaction[i] = tempValue[0];
			listVoucherValue[i] = tempValue[1];
			listVoucher[i] = tempValue[2];
		}
	}

	for (i = 0; i < listVoucher.length; i++) {
		total += parseInt(listVoucherValue[i]);
	}

	var transaction = {
		partner_id : $("#partnerList").find(":selected").val(),
		bank_account : $("#bankAccount").val(),
		total_cashout : total,
		payment_method : "transfer",
		transactions : listTransaction,
		vouchers : listVoucher
	};

	$.ajax({
		url: '/v1/ui/cashout?token='+token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(transaction),
		success: function (data) {
			console.log(data);
			$("#success-page").attr("style","display:block");
			$("#cashoutId").val(data.data);
			$("#transaction").attr("style","display:none");
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function print(){
	window.location = "/voucher/print?id="+$("#cashoutId").val();
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
