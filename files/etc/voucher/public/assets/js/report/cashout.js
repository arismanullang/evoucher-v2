$(document).ready(function () {
  getCashout();
  $(document).ajaxStart(function(){
    // Show image container
    $(".cssload-loader").show();
   });
   $(document).ajaxComplete(function(){
    // Hide image container
    $(".cssload-loader").hide();
   });
});

function getCashout() {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/cashout?token=' + token,
		type: 'get',
		success: function (data) {
			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var dataSet = [];
			var arrData = data.data;
			var i;

			for (i = 0; i < arrData.length; i++) {
        if(arrData[i].status == "voided"){
          var button = "<button type='button' onclick='detail(\"" + arrData[i].id + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button><span class='label label-danger'>voided</span>";
        }else {
          var button = "<button type='button' onclick='detail(\"" + arrData[i].id + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>";
        }

				var tempArray = [
					arrData[i].partner_id.toUpperCase()
					, arrData[i].cashout_code
					, arrData[i].bank_account
					, arrData[i].bank_account_ref_number
					, "Rp. " + addDecimalPoints(arrData[i].total_cashout)
					, dateFormat(new Date(arrData[i].created_at), 'isoDateTime')
					, button
				];

				dataSet.push(tempArray);
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'lBrtip',
				buttons: [
					'copy', 'csv', 'excel', 'pdf', 'print'
				],
				"order": [[5, "desc"]],
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
		}
	});
}

function detail(id){
	window.location = "/report/cashout-detail?id="+id;
}
