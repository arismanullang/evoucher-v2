$( document ).ready(function() {
	var id = findGetParameter('id');
	cashout(id);
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
				for ( y = 0; y < result.transactions[i].vouchers.length; y++){
					var date = new Date(result.created_at);
					var body = "<td>"+result.transactions[i].transaction_code+"</td>"
						+ "<td>"+result.transactions[i].vouchers[y].voucher_code+"</td>"
						+ "<td>Rp. "+addDecimalPoints(result.transactions[i].vouchers[y].voucher_value)+",00</td>"
						+ "<td>"+date.toDateString() + ", " +toTwoDigit(date.getHours()) + ":" + toTwoDigit(date.getMinutes())+"</td>"
						+ "<td>"+result.transactions[i].vouchers[y].state.toUpperCase()+"</td>";
					var li = $("<tr class='text-center'></tr>");
					li.html(body);
					li.appendTo('#list-transaction');
				}
			}

			getPartnerName(result.partner_id, result.bank_account);
			$("#cashoutCode").html(result.cashout_code);
			$("#total").html("Rp. " + addDecimalPoints(result.total_cashout) + ",00");
		}
	});
}

function getPartnerName(id, cashoutCode){
	$.ajax({
		url: '/v1/ui/partner?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;
			$("#tenant").html(result[0].name + " ("+ cashoutCode +")");
		}
	});
}
