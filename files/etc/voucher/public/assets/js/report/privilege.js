$( document ).ready(function() {
	var total = 0;
	$('#transaction').submit(function(e) {
		e.preventDefault();
		addElem();
		return false;
	});

	$('#datepicker-privilege-to').change(function () {
      var startVal = $('#privilege-from').val();
      var endVal =  $('#privilege-to').val();

      if(startVal.length > 0 && endVal.length > 0){
        var startDate = new Date(startVal);
        var endDate = new Date(endVal);
        endDate.setHours(23);
        endDate.setMinutes(59);
        endDate.setSeconds(59);

        getPrivilege(dateFormat(startDate, 'isoUtcDateTime'), dateFormat(endDate, 'isoUtcDateTime'));
      }

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

function getPrivilege(dateFrom, dateTo) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/privilege?token=' + token + '&dateFrom=' + dateFrom + '&dateTo='+dateTo,
		type: 'get',
		success: function (data) {
			var result = data.data;
			var partner = {};
			var transaction = {};

			var dataSet = [];
			if(result != null){
				for(var i = 0; i < result.length; i++){
					if(transaction[result[i].partner_id] == null){
						partner[result[i].partner_id] = result[i].partner_name;
						transaction[result[i].partner_id] = 1;
					}else{
						transaction[result[i].partner_id]++;
					}
				}

				var keys = Object.keys(transaction);

				for(var i = 0; i < keys.length; i++){
					var button = "<button type='button' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>";

					dataSet[i] = [
						partner[keys[i]]
						, transaction[keys[i]]
					];
				}
			}

			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'rt',
				"order": [[1, "desc"]],
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
      var inputSearchClass = 'datatable_input_col_search';
			var columnInputs = $('thead .' + inputSearchClass);
			for (i = 0; i < columnInputs.length; i++) {
				if (columnInputs.get(i).tagName.toLowerCase() == "select") {
					columnInputs[i].onchange = function () {
						table.fnFilter(this.value, columnInputs.index(this));
					};
				} else {
					columnInputs[i].onkeyup = function () {
						table.fnFilter(this.value, columnInputs.index(this));
					};
				}
			}
		},
		error: function(){
			$("#listTransaction").html(" ");
		}
	});
}

function detail(url, date){
	// window.location = "/voucher/cashout-detail?partner="+url+"&date=" + date;
}


(function () {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$('.datepicker-privilege-from').datepicker({
      container: '#datepicker-privilege-from',
			autoclose: true
		}).on('changeDate', function (selected) {
      var minDate = new Date(selected.date.valueOf());
			$('.datepicker-privilege-to').datepicker('setStartDate', minDate);
		});
		$('.datepicker-privilege-to').datepicker({
			container: '#datepicker-privilege-to',
			autoclose: true
		});

	}

})();
