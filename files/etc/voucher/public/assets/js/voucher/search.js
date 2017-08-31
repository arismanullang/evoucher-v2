$(document).ready(function () {
	var holder = findGetParameter("holder");
	getVoucher(holder);
});

function getVoucher(holder) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/voucher?holder=' + holder + '&token=' + token,
		type: 'get',
		success: function (data) {
			arrData = data.data;
			var i;
			var dataSet = [];
			for (i = 0; i < arrData.length; i++) {
				var date1 = arrData[i].valid_at.substring(0, 10).split("-");
				var date2 = arrData[i].expired_at.substring(0, 10).split("-");

				var dateStart = new Date(date1[0], date1[1] - 1, date1[2]);
				var dateEnd = new Date(date2[0], date2[1] - 1, date2[2]);
				var dateNow_ms = Date.now();

				var one_day = 1000 * 60 * 60 * 24;
				var dateStart_ms = dateStart.getTime();
				var dateEnd_ms = dateEnd.getTime();
				// var dateNow_ms = dateNow.getTime();

				var diffNow = Math.round((dateEnd_ms - dateStart_ms) / one_day);
				var persen = 100;

				if (dateStart_ms < dateNow_ms) {
					diffNow = Math.round((dateEnd_ms - dateNow_ms) / one_day);
					var diffTotal = Math.round((dateEnd_ms - dateStart_ms) / one_day);
					persen = diffNow / diffTotal * 100;
				}

				if (dateNow_ms > dateEnd_ms) {
					persen = -1;
				}

				diffNow = diffNow + " hari";

				if (persen < 0) {
					diffNow = "Expired";
				}

				var status = "";
				switch (arrData[i].state) {
					case "paid":
						status = "PAID";
						break;
					case "used":
						status = "REDEEM";
						break;
					default:
						status = "ISSUED";
				}

				var tempArray = [
					arrData[i].holder_description.toUpperCase()
					, arrData[i].program_name.toUpperCase()
					, arrData[i].voucher_code.toUpperCase()
					, arrData[i].reference_no.toUpperCase()
					, status.toUpperCase()
					, "<button type='button' onclick='detail(\"" + arrData[i].id + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"
				];

				dataSet.push(tempArray);
			}
			console.log(dataSet);

			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'lBrtip',
				buttons: [
					'copy', 'csv', 'excel', 'pdf', 'print'
				],
				"order": [[1, "asc"]],
				columns: [
					{title: "Holder"},
					{title: "Program"},
					{title: "Voucher Code"},
					{title: "Reference Code"},
					{title: "Status"},
					{title: "Action"}
				],
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

function detail(url) {
	window.location = "/voucher/check?id=" + url;
}
