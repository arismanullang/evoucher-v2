$( document ).ready(function() {
	var id = findGetParameter('id');
	cashout(id);
	getAccount();
	var date = new Date();
	$("#date").html(date.toDateString() + ", " + toTwoDigit(date.getHours()) + ":" + toTwoDigit(date.getMinutes()));
});

function cashout(id){
	var transcations = id;
	$.ajax({
		url: '/v1/ui/cashout/print?transcation_code='+encodeURIComponent(transcations)+'&token='+token,
		type: 'get',
		success: function (data) {
			console.log(data.data);
			var result = data.data;
			var total = 0;
			for ( var i = 0; i < result.transactions.length; i++){
				var date = new Date(result.transactions[i].created_at);
				var body = "<td>"+result.transactions[i].transaction_id+"</td>"
					+ "<td>"+result.transactions[i].voucher_id+"</td>"
					+ "<td>Rp. "+addDecimalPoints(parseInt(result.transactions[i].voucher_value))+",00</td>"
					+ "<td>"+date.toDateString() + ", " +toTwoDigit(date.getHours()) + ":" + toTwoDigit(date.getMinutes())+"</td>";
				var li = $("<tr class='text-center'></tr>");
				li.html(body);
				li.appendTo('#list-transaction');
			}

			getPartnerName(result.partner_id);
			$("#total").html("Rp. " + addDecimalPoints(result.total_cashout));
			$("#bank-account").html(result.bank_account + ", Acc. No. " + result.bank_account_number);
			$("#bank-account-ref-number").html(result.bank_account_ref_number);
			$("#company-name").html(result.bank_account_company);
			$("#cashout-code").html(result.cashout_code);
		}
	});
}

function getPartnerName(id){
	$.ajax({
		url: '/v1/ui/partner?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {

			// bank_account	BCA
			// bank_account_company	Shigeru
			// bank_account_number	1231231234
			// bank_account_ref_number	11111111111111

			var result = data.data;
			$("#tenant").html(result[0].name);
		}
	});
}

function getAccount(){
	$.ajax({
		url: '/v1/ui/account?token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;
			$("#account-name").html(result.name);
			$("#account-bulding").html(result.building);
			$("#account-address").html(result.address);
			$("#account-province").html(result.city + ", " +result.province + ", " + result.zip_code);
		}
	});
}
