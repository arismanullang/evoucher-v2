$(document).ready(function () {
	var id = findGetParameter('id');
	getPartner(id);
	getPerformance(id, 'all');
	getVoucher(id, 'all');
	getProgram(id);

	$("#performanceType").change(function () {
		var type = this.value;
		getPerformance(id, type);
		getVoucher(id, type);
	});

	$(".select2").select2();
});

function getVoucher(id, type) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/voucher/partner?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;
			var limit = arrData.length;

			var dataSet = [];
			for (i = 0; i < limit; i++) {
				var dateValid = new Date(arrData[i].updated_at.Time);
				if(arrData[i].updated_at.Time == "0001-01-01T00:00:00Z"){
					dateValid = new Date(arrData[i].created_at);
				}
				var voucherState = '';
				if(arrData[i].state == 'used'){
					voucherState = 'used';
				}else if(arrData[i].state == 'created'){
					voucherState = 'redeemed';
				}else{
					voucherState = 'paid';
				}
				var button = "<button type='button' onclick='detail(\"" + arrData[i].id + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"
				dataSet[i] = [
					arrData[i].program_name.toUpperCase()
					, arrData[i].voucher_code
					, arrData[i].holder_description.String.toUpperCase()
					, dateValid.toDateString().toUpperCase()
					, voucherState.toUpperCase()
					, button
				];
			}

			$("#totalValidVoucher").html(limit);
			initGraphReport(dataSet);
		},
		error: function(data){
			$('#tableVoucher').attr("style", "display:none");
		}
	});
}

function getPerformance(id, type) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/partner/performance?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;
			$("#totalProgram").html(arrData.program);
			$("#totalTransaction").html(arrData.transaction_code);
			$("#totalVoucherValue").html("Rp. " + addDecimalPoints(arrData.transaction_value));
			$("#totalUsedVoucher").html(arrData.voucher_used);
			$("#totalCustomer").html(arrData.customer);
		}
	});
}

function makeQRCode(id, name) {
	var qrcode = new QRCode("qrcode", {
		text: name,
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
			$("#partnerTitle").html(arrData.name.toUpperCase());
			$("#partner-name").html(arrData.name);
			$("#serial-number").html(arrData.serial_number.String);
			$("#tag").html(arrData.tag.String);
			$("#desciption").html(arrData.description.String);
			$('#bank-account').html(arrData.company_name + ", " + arrData.bank_name + " - " + arrData.bank_account_number);

			makeQRCode(id, arrData.name);
		}
	});
}

function getProgram(id) {
	$.ajax({
		url: '/v1/ui/partner/programs?id=' + id + "&token=" + token,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;
			var limit = arrData.length;

			var dataSet = [];
			for (i = 0; i < limit; i++) {
				var type = 'mobile application';
				if(arrData[i].type == 'bulk'){
					type = 'email blast';
				} else if(arrData[i].type == 'gift'){
					type = 'gift voucher';
				}

				dataSet[i] = [
					arrData[i].name.toUpperCase()
					, type.toUpperCase()
					, new Date(arrData[i].start_date).toDateString().toUpperCase()
					, new Date(arrData[i].end_date).toDateString().toUpperCase()
					, arrData[i].voucher
				];
			}
			initProgramReport(dataSet);
		}
	});
}

function detail(id) {
	window.location = "/voucher/check?id=" + id;
}

function edit() {
	window.location = "/partner/update?id=" + findGetParameter('id');
}

function initGraphReport(dataSet) {
	$('#tableVoucher').attr("style", "display:block");

	if ($.fn.DataTable.isDataTable("#datatable1")) {
		$('#datatable1').DataTable().clear().destroy();
	}

	var table = $('#datatable1').dataTable({
		data: dataSet,
		dom: 'lBrtip',
		buttons: [
			'copy', 'csv', 'excel', 'pdf', 'print'
		],
		"order": [[4, "desc"]],
		oLanguage: {
			sSearch: '<em class="ion-search"></em>',
			sLengthMenu: 'Display _MENU_ records',
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

function initProgramReport(dataSet) {
	if ($.fn.DataTable.isDataTable("#datatableProgram")) {
		$('#datatableProgram').DataTable().clear().destroy();
	}

	var table = $('#datatableProgram').dataTable({
		data: dataSet,
		dom: 'lBrtip',
		buttons: [
			'copy', 'csv', 'excel', 'pdf', 'print'
		],
		"order": [[2, "desc"]],
		oLanguage: {
			sSearch: '<em class="ion-search"></em>',
			sLengthMenu: 'Display _MENU_ records',
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
	var inputSearchClass = 'datatableProgram_input_col_search';
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
