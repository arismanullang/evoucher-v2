$( document ).ready(function() {
	var partner = findGetParameter("partner");
	console.log(encodeURIComponent(partner));
	cashout(partner);
	var date = new Date();
	$("#date").html(date.toDateString() + ", " + date.getHours() + ":" + date.getMinutes());
});

function cashout(partner){
	var transactionCode = [];
	var elem = $("*[name='list-transaction-code']");
	var i = 0;
	for ( i = 0; i < elem.length; i++){
		transactionCode[i] = elem[i].innerHTML;
	}

	$.ajax({
		url: '/v1/print/transaction/cashout?partner='+partner.replace("&","*")+'&token='+token,
		type: 'get',
		success: function (data) {
			console.log(data.data);
			var result = data.data;
			for ( i = 0; i < result.length; i++){
				var date = new Date(result[i].redeemed)
				var body = "<td>"+result[i].partner_name+"</td>"
					+ "<td>"+result[i].transaction_id+"</td>"
					+ "<td>"+date.toDateString() + ", " + date.getHours() + ":" + date.getMinutes()+"</td>"
				var li = $("<tr class='text-center'></tr>");
				li.html(body);
				li.appendTo('#list-transaction');

			}
		}
	});
}
