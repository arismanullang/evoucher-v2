$(document).ready(function () {
	var id = findGetParameter('id');
	getPartner(id);
	getPerformance(id);
	makeCode(id);
	getVoucher(id);
});

function getVoucher(id) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/cashout/partner?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;
			var limit = arrData.length;

			var dataSet = [];
			for (i = 0; i < limit; i++) {
				var dateValid = new Date(arrData[i].updated_at);
				dataSet[i] = [
					arrData[i].program_name.toUpperCase()
					, arrData[i].voucher_code
					, arrData[i].holder_description.toUpperCase()
					, dateValid.toDateString().toUpperCase()
					, arrData[i].state.toUpperCase()
				];
			}

			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'lBrtip',
				buttons: [
					'copy', 'csv', 'excel', 'pdf', 'print'
				],
				"order": [[3, "desc"]],
				columns: [
					{title: "PROGRAM"},
					{title: "VOUCHER"},
					{title: "HOLDER"},
					{title: "LAST UPDATE"},
					{title: "STATUS"}
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

function makeCode(id) {
	var qrcode = new QRCode("qrcode", {
		text: id,
		width: 150,
		height: 150,
		colorDark: "#000000",
		colorLight: "#ffffff",
		correctLevel: QRCode.CorrectLevel.H
	});
	qrcode.makeCode(id);
}

function getPartner(id) {
	$.ajax({
		url: '/v1/ui/partner?id=' + id + "&token=" + token,
		type: 'get',
		success: function (data) {
			var arrData = data.data[0];
			$("#initial").html(arrData.name.charAt(0).toUpperCase());
			$("#partner-title").html(arrData.name.toUpperCase());
			$("#partner-name").html(arrData.name);
			$("#serial-number").html(arrData.serial_number.String);
			$("#tag").html(arrData.tag.String);
			$("#desciption").html(arrData.description.String);
		}
	});
}

function getPerformance(id) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/partner/performance?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;
			$("#total-program").html(arrData.program);
			$("#total-transaction").html(arrData.transaction_code);
			$("#total-voucher-value").html("Rp. " + addDecimalPoints(arrData.transaction_value) + ",00");
			$("#total-valid-voucher").html(arrData.voucher_generated);
			$("#total-used-voucher").html(arrData.voucher_used);
			$("#total-customer").html(arrData.customer);
		}
	});
}

function edit() {
	window.location = "/partner/update?id=" + findGetParameter('id');
}
