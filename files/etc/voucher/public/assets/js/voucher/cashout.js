$( document ).ready(function() {
  var total = 0;

  var last30Date = new Date();
  last30Date.setDate(last30Date.getDate() - 30);

  $('#trx-from').val(dateFormat(last30Date, "mm/dd/yyyy"));
  $('#trx-to').val(dateFormat(new Date(), "mm/dd/yyyy"));

	$('#transaction').submit(function(e) {
		e.preventDefault();
		addElem();
		return false;
  });

  $('#datepicker-trx-to').change(function () {
    var startVal = $('#trx-from').val();
    var endVal =  $('#trx-to').val();

    if(startVal.length > 0 && endVal.length > 0){
      var startDate = new Date(startVal);
      var endDate = new Date(endVal);
      endDate.setHours(23);
      endDate.setMinutes(59);
      endDate.setSeconds(59);

      getTransactionByDate(dateFormat(startDate, 'isoUtcDateTime'), dateFormat(endDate, 'isoUtcDateTime'));
    }

});

  $('.datepicker-trx-from').datepicker({
    container: '#datepicker-trx-from',
    autoclose: true
  });

	$('.datepicker-trx-to').datepicker({
		container: '#datepicker-trx-to',
		autoclose: true
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

function getTransactionByDate(dateFrom, dateTo) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/date?token=' + token + '&state=used&start_date=' + dateFrom + '&end_date=' + dateTo,
		type: 'get',
		success: function (data) {
			var result = data.data;
      var partner = {};
      var transactionId = [];
			var transaction = {};
			var transactionValue = {};

			var dataSet = [];
			if(result != null){
				for(var i = 0; i < result.length; i++){
					if(transaction[result[i].partner_id] == null){
            partner[result[i].partner_id] = result[i].partner_name;
            transactionId.push(result[i].transaction_id);
						transaction[result[i].partner_id] = 1;
						transactionValue[result[i].partner_id] = result[i].voucher_value;
					}else{
            if(!transactionId.includes(result[i].transaction_id)){
              console.log(transactionId + "-" + result[i].transaction_id);
              transactionId.push(result[i].transaction_id);
              transaction[result[i].partner_id]++ ;
            }
						transactionValue[result[i].partner_id] = transactionValue[result[i].partner_id] + result[i].voucher_value;
					}
				}

				var keys = Object.keys(transaction);

				for(var i = 0; i < keys.length; i++){
					var button = "<button type='button' onclick='detail(\"" + keys[i] + "\",\""+dateFrom+"\",\""+dateTo+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>";

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

function detail(url, startDate, endDate){
	window.location = "/voucher/cashout-detail?partner="+url+"&start_date=" + startDate+"&end_date=" + endDate;
}
