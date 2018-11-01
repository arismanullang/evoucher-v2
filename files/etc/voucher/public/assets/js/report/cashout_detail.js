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
      $("#cashout-void").val(result.id)
      $("#reimburse-status-row").hide();
      if(result.status != "created"){
        $("#reimburse-status-row").show();
        $("#button-void-row").hide();
        $("#reimburse-status").html("<span class='label label-danger'>"+ result.status +"</span>");
      }

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


function voidCashout(id) {
  var param = {
		cashout_id: id,
		description: "voided from cms",
  };

	var status = "";
	$.ajax({
		url: '/v1/cashout/void?token=' + token,
    type: 'post',
    dataType: 'json',
    data: JSON.stringify(param),
		success: function (data) {
			cashout(id);
			swal('Void Success!');
		},
		error: function (xhr, ajaxOptions, thrownError) {
			swal('Cashout cannot be voided', xhr.responseJSON.errors.detail);
		}
	});
	return status;
}

(function () {
	'use strict';

  $(runSweetAlert);

	function runSweetAlert() {
		$(document).on('click', '.swal-void-cashout', function (e) {
			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want void this cashout?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Void',
					closeOnConfirm: false
				},
				function () {
					voidCashout(e.target.value);
				});

		});
	}

})();
