$(document).ready(function () {
	var id = findGetParameter('id');
	cashout(id);
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
			$("#ref-number").html(result.bank_account + ", Acc. No. " + result.bank_account_number + " - " + result.bank_account_ref_number);
			$("#total-transaction").html(result.transactions.length);
			$("#total-reimburse").html(result.total_cashout);
			$("#company-name").html(result.bank_account_company);
			$("#reimburse-code").html(result.cashout_code);
			$("#reimburse-date").html(new Date(result.created_at));
		}
	});
}

function getPartnerName(id){
	$.ajax({
		url: '/v1/ui/partner?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;
			$("#partner-name").html(result[0].name);
		}
	});
}
