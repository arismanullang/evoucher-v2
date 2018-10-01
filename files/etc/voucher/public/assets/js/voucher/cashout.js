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
  });

  $(document).ajaxStart(function(){
    // Show image container
    $(".cssload-loader").show();
   });
   $(document).ajaxComplete(function(){
    // Hide image container
    $(".cssload-loader").hide();
   });

});

function getTransactionByDate(date) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/date?token=' + token + '&date=' + date,
		type: 'get',
		success: function (data) {
			var result = data.data;
			var partner = {};
			var transaction = {};
			var transactionValue = {};

			var dataSet = [];
			if(result != null){
				for(var i = 0; i < result.length; i++){
					if(transaction[result[i].partner_id] == null){
						partner[result[i].partner_id] = result[i].partner_name;
						transaction[result[i].partner_id] = 1;
						transactionValue[result[i].partner_id] = result[i].vouchers.length * result[i].voucher_value;
					}else{
						transaction[result[i].partner_id]++;
						transactionValue[result[i].partner_id] = transactionValue[result[i].partner_id] + (result[i].vouchers.length * result[i].voucher_value);
					}
				}

				var keys = Object.keys(transaction);

				for(var i = 0; i < keys.length; i++){
					var button = "<button type='button' onclick='detail(\"" + keys[i] + "\",\""+date+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>";

					dataSet[i] = [
						partner[keys[i]]
						, transaction[keys[i]]
						, "Rp. " + addDecimalPoints(transactionValue[keys[i]])
						, button
					];
				}
			}

			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'rt',
				"order": [[2, "desc"]],
				paging: false,
				oLanguage: {
					sSearch: '<em class="ion-search"></em>',
					sLengthMenu: '_MENU_ records per page',
					info: 'Showing page _PAGE_ of _PAGES_',
					zeroRecords: 'Nothing found - sorry',
					infoEmpty: 'No records available',
					infoFiltered: '(filtered from _MAX_ total records)',
					oPaginate: {
						sNext: '<em class="ion-ios-arrow-right"></em>',
						sPrevious: '<em class="ion-ios-arrow-left"></em>'
					}
				}
			});
		},
		error: function(){
			$("#listTransaction").html(" ");
		}
	});
}

function detail(url, date){
	window.location = "/voucher/cashout-detail?partner="+url+"&date=" + date;
}
