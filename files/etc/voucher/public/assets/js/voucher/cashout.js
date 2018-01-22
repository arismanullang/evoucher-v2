$( document ).ready(function() {
	var total = 0;
	$('#transaction').submit(function(e) {
		e.preventDefault();
		addElem();
		return false;
	});

	$('.datepicker-transaction').datepicker({
		container: '#datepicker-transaction',
		autoclose: true
	});

	$('#transaction-date').change(function () {
			getTransactionByDate($('#transaction-date').val());
	}
	);
});

function getTransactionByDate(date) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/date?token=' + token + '&date=' + date,
		type: 'get',
		success: function (data) {
			var result = data.data;

		},
		error: function(){
			$("#listTransaction").html(" ");
		}
	});
}
