$( document ).ready(function() {
	$('#transaction').submit(function(e) {
		e.preventDefault();
		addElem();
		return false;
	});
});

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

	console.log(textArea.value);
	console.log(transaction);
	$.ajax({
		url: '/v1/ui/transaction/cashout/update?token='+token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(transaction),
		success: function () {
			localStorage.setItem("transaction_cashout", "");
			localStorage.setItem("transaction_cashout", transactions);
			window.location = "/voucher/print?token="+token;
		}
	});
}

function addElem(){
	var id = $('#transaction-code').val();
	$.ajax({
		url: '/v1/ui/transaction?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			console.log("Render Data");
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

				var body = "<td name='list-transaction-code' class='text-ellipsis'>"+result.transaction_code+"</td>"
					+ "<td name='list-transaction-partner' class='text-ellipsis'>"+result.partner_name+"</td>"
					+ "<td name='list-transaction-value' class='text-ellipsis'>"+result.discount_value+"</td>"
					+ "<td name='list-transaction-date' class='text-ellipsis'>"+date.toDateString() + ", " + date.getHours() + ":" + date.getMinutes()+"</td>"
					+ "<td name='list-transaction-action'><button type='button' onclick='removeElem(this)' class='btn btn-flat btn-sm btn-info pull-right'><em class='ion-close-circled'></em></button></td>";
				var li = $("<tr class='msg-display clickable'></tr>");
				li.html(body);
				li.appendTo('#list-transaction');
				$("#error").html('');

			}
		},
		error: function (data) {
			$('#transaction-code').val('');
			$("#error").html(data.responseJSON.errors.detail);
		}
	});
}

function removeElem(elem){
	console.log("remove");
	$(elem).parent().closest('tr').remove();
}
