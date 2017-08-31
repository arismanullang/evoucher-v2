$( document ).ready(function() {
	cashout();
	var date = new Date();
	$("#date").html(date.toDateString() + ", " + toTwoDigit(date.getHours()) + ":" + toTwoDigit(date.getMinutes()));
});

function cashout(){
	var transcations = localStorage.getItem("transaction_cashout");
	$.ajax({
		url: '/v1/ui/transaction/cashout/print?transcation_code='+encodeURIComponent(transcations)+'&token='+token,
		type: 'get',
		success: function (data) {
			console.log(data.data);
			var result = data.data;
			var total = 0;
			for ( i = 0; i < result.length; i++){
				for ( y = 0; y < result[i].vouchers.length; y++){
					var date = new Date(result[i].redeemed);
					var body = "<td>"+result[i].transaction_code+"</td>"
						+ "<td>"+result[i].vouchers[y].voucher_code+"</td>"
						+ "<td>"+date.toDateString() + ", " +toTwoDigit(date.getHours()) + ":" + toTwoDigit(date.getMinutes())+"</td>"
						+ "<td>"+result[i].state.toUpperCase()+"</td>";
					var li = $("<tr class='text-center'></tr>");
					li.html(body);
					li.appendTo('#list-transaction');
					total += result[i].voucher_value;
				}
			}

			$("#tenant").html(result[0].partner_name);
			$("#total").html("Rp. " + addDecimalPoints(total.toString()) + ",00");
		}
	});
}
