$( document ).ready(function() {
	var total = 0;
	$('#transaction').submit(function(e) {
		e.preventDefault();
		addElem();
		return false;
	});
	$('.select2').select2();

	var partner = findGetParameter('partner');
	getPartner(partner);
	getTransactionByPartner(partner);

	$('#partner-list').change(function () {
		getTransactionByPartner(this.value);
	});
	$("#transactionAll").change(function () {
		var _this = $(this);
		$('#list-transaction').find("input.transaction").prop('checked', _this.prop('checked'));

		updateTotal();
	});
});

function updateTotal() {
	var li = $("input[class=transaction]:checked");
	total = 0;
	for( var i = 0; i < li.length; i++){
		total += parseInt(li[i].value.split(";")[1]);
	}

	$('#total-transaction').html("Rp. "+addDecimalPoints(total));
}

function getPartner(id) {
	$.ajax({
		url: '/v1/ui/partner?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data[0];

			$('#partner').html(result.name);
			$('#bank-account-company').html(result.company_name);
			$('#bank-account').html(result.bank_name);
			$('#bank-account-number').html(result.bank_account_number);
		},
		error: function (data) {
		}
	});
}

function getTransactionByPartner(partnerId) {
	var date = findGetParameter('date');
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/cashout/partner?token=' + token + '&partner='+partnerId + '&date='+date,
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
						+ "<td class='text-ellipsis'>Rp. "+addDecimalPoints(arrData[i].voucher_value)+"</td>"
						+ "<td class='text-ellipsis'>"+date.toDateString() + ", " + date.getHours() + ":" + date.getMinutes()+"</td>"
					var li = $("<tr></tr>");
					li.html(body);
					li.appendTo('#listTransaction');
					voucher++;
				}
			}

			$('.transaction').change(function () {
				updateTotal();
			});

			if(voucher > 5){
				$("#tableTransaction").attr("style","overflow:scroll; max-height: 300px;");
			}else{
				$("#tableTransaction").removeAttr("style");
			}
		},
		error: function(){
			$("#listTransaction").html(" ");
		}
	});
}

function cashout(){
	var partner = findGetParameter('partner');

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
	if(total == 0) {
		return
	}

	if($('#bank-account-ref-number').val() == '') {
		return
	}
	var transaction = {
		partner_id: partner,
		bank_account: $('#bank-account').html(),
		bank_account_company: $('#bank-account-company').html(),
		bank_account_number: $('#bank-account-number').html(),
		bank_account_ref_number: $('#bank-account-ref-number').val(),
		total_cashout: total,
		payment_method: "transfer",
		transactions: listTransaction,
		vouchers: listVoucher
	};

	$.ajax({
		url: '/v1/ui/cashout?token='+token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(transaction),
		success: function (data) {
			window.location = "/voucher/cashout-success?id="+data.data;
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
