$(document).ready(function () {
	var id = findGetParameter('id');
	$("#programId").val(id);
	$("#token").val(token);
	getVoucher(id);
	getPartner(id);
});

function chart(remaining, outstanding, pending, paid) {
	var pieData = [{
		'label': "<div class='col-md-6'>Stock : </div><div class='col-md-6'>Rp. " + addDecimalPoints(remaining) + ",00</div>",
		'color': '#e4eff7',
		'data': remaining
	}, {
		'label': "<div class='col-md-6'>Outstanding : </div><div class='col-md-6'>Rp. " + addDecimalPoints(outstanding) + ',00</div>',
		'color': '#f4eff7',
		'data': outstanding
	}, {
		'label': "<div class='col-md-6'>Pending : </div><div class='col-md-6'>Rp. " + addDecimalPoints(pending) + ',00</div>',
		'color': '#AFEFEF',
		'data': pending
	}, {
		'label': "<div class='col-md-6'>Paid : </div><div class='col-md-6'>Rp. " + addDecimalPoints(paid) + ',00</div>',
		'color': '#69CDCD',
		'data': paid
	}];
	var pieOptions = {
		series: {
			pie: {
				show: true
			}
		}
	};
	$('#pie').plot(pieData, pieOptions);
}

function getVoucher(id) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/voucher?program_id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;
			var limit = arrData.length;
			var used = 0;
			var paid = 0;

			var dataSet = [];
			for (i = 0; i < limit; i++) {
				var voucherState = "redeemed";
				if (arrData[i].state == "used") {
					used++;
					voucherState = "used";
				}

				if (arrData[i].state == "paid") {
					used++;
					paid++;
					voucherState = "paid";
				}

				var dateValid = new Date(arrData[i].valid_at);
				var dateExpired = new Date(arrData[i].expired_at);

				var button = "<button type='button' onclick='detail(\"" + arrData[i].id + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"
				dataSet[i] = [
					arrData[i].voucher_code
					, arrData[i].holder_description.toUpperCase()
					, dateValid.toDateString().toUpperCase()
					, dateExpired.toDateString().toUpperCase()
					, voucherState.toUpperCase()
					, button
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
				"order": [[4, "desc"]],
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

			getProgram(id, arrData.length, used, paid);
		},
		error: function (data) {
			getProgram(id, 0, 0, 0);
			$('#table-voucher').attr('style', 'display:none');
		}
	});
}

function getPartner(id) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/program/partner?program_id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			if ($.fn.DataTable.isDataTable("#datatablePartner")) {
				$('#datatablePartner').DataTable().clear().destroy();
			}

			var arrData = data.data;
			var i;
			var dataSet = [];
			for (i = 0; i < arrData.length; i++) {
				var tempArray = [
					arrData[i].name.toUpperCase()
					, arrData[i].transactions
					, arrData[i].vouchers
					, "Rp. "+addDecimalPoints(arrData[i].transaction_values)+",00"
				];

				dataSet.push(tempArray);
			}

			var table = $('#datatablePartner').dataTable({
				data: dataSet,
				dom: 'lBrtip',
				buttons: [
					'copy', 'csv', 'excel', 'pdf', 'print'
				],
				"order": [[2, "desc"]],
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
			var inputSearchClass = 'partner_datatable_input_col_search';
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
		error: function (data) {
			$("<div class='card-body text-center'>No Partner Found</div>").appendTo('#tablePartner');
		}
	});
}

function getProgram(id, voucher, used, paid) {
	var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
	var arrData = [];
	$.ajax({
		url: '/v1/ui/program/detail?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var result = data.data[0];

			var date1 = result.start_date.substring(0, 10).split("-");
			var date2 = result.end_date.substring(0, 10).split("-");
			var startDate = date1[2] + " " + months[parseInt(date1[1]) - 1] + " " + date1[0];
			var endDate = date2[2] + " " + months[parseInt(date2[1]) - 1] + " " + date2[0];

			var period = startDate + " - " + endDate;
			var programType = "Email Blast"
			var remainingVoucher = result.max_quantity_voucher;
			if (voucher != null) {
				remainingVoucher = result.max_quantity_voucher - voucher;
			}

			switch(result.type){
				case 'on-demand':
					programType = "Mobile App"
					$('#button-campaign').attr('style', 'display:none');
					$('#button-voucher').attr('style', 'display:none');
					break;
				case 'gift':
					programType = "Gift Voucher"
					$('#button-campaign').attr('style', 'display:none');
					$('#button-voucher').attr('style', 'display:none');
					break;
			}

			$("#visibility").val(result.visibility);
			if(result.visibility == true){
				$("#label-visibility").html("Show");
			}else{
				$("#label-visibility").html("Hidden");
			}

			// Program
			$('#programName').html(result.name);
			$('#programNames').html(result.name);
			$('#programDescription').html(result.description);
			$('#programType').html(programType);
			$('#conversionRate').html(result.voucher_price + ' Point');
			$('#voucherValue').html('Rp. ' + addDecimalPoints(result.voucher_value) + ',00');
			$('#period').html(period);
			$('#programTnc').html(result.tnc);

			// Voucher
			$('#maxQuantityVoucher').html(result.max_quantity_voucher);
			$('#remainingVoucher').html(remainingVoucher);
			$('#createdVoucher').html(voucher);
			$('#usedVoucher').html(used);
			$('#paidVoucher').html(paid);

			// Cashflow
			totalMon = parseInt(result.voucher_value) * parseInt(result.max_quantity_voucher);
			$('#totalMoney').html('Rp. ' + addDecimalPoints(totalMon) + ',00');

			remainingMon = (parseInt(result.max_quantity_voucher) - parseInt(voucher)) * parseInt(result.voucher_value);
			$('#remainingMoney').html('Rp. ' + addDecimalPoints(remainingMon) + ',00');

			outstandingMon = (parseInt(voucher) - parseInt(used)) * parseInt(result.voucher_value);
			$('#remainingMoney').html('Rp. ' + addDecimalPoints(remainingMon) + ',00');

			pendingPay = parseInt(result.voucher_value) * (parseInt(used)-parseInt(paid));
			$('#pendingPayment').html('Rp. ' + addDecimalPoints(pendingPay) + ',00');

			paidPay = parseInt(result.voucher_value) * parseInt(paid);
			$('#paidPayment').html('Rp. ' + addDecimalPoints(paidPay) + ',00');

			remainingPercent = remainingMon / totalMon * 100;

			pendingPercent = pendingPay / totalMon * 100;

			paidPercent = paidPay / totalMon * 100;

			chart(remainingMon, outstandingMon, pendingPay, paidPay);
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function changeVisibility() {
	var status = "";
	var id = $("#programId").val();
	$.ajax({
		url: '/v1/ui/program/visibility?id=' + id + '&visible='+$('#visibility').val()+'&token=' + token,
		type: 'get',
		success: function (data) {
			swal({
					title: 'Success',
					text: 'Change Visibility Success',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: true
				},
				function() {
					window.location.reload();
				});
		},
		error: function (xhr, ajaxOptions, thrownError) {
			swal('Error!', xhr.responseJSON.errors.detail);
		}
	});
	return status;
}

function programCampaign(){
	var id = findGetParameter("id");
	window.location = "/program/campaign?id=" + id;
}

function editProgram() {
	var id = findGetParameter("id");
	window.location = "/program/update?id=" + id;
}

function detail(id) {
	window.location = "/voucher/check?id=" + id;
}

function sendEmail() {
	var id = findGetParameter("id");
	window.location = "/user/list/send?id=" + id;
}
