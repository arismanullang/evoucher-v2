$(document).ready(function () {
	var id = findGetParameter('id');
	$("#program-id").val(id);
	$("#token").val(token);
	getVoucher(id);
	getPartner(id);
});

function chart(remaining, outstanding, pending, paid) {
	var pieData = [{
		'label': 'Stock : Rp. ' + addDecimalPoints(remaining) + ',00',
		'color': '#e4eff7',
		'data': remaining
	}, {
		'label': 'Outstanding : Rp. ' + addDecimalPoints(outstanding) + ',00',
		'color': '#f4eff7',
		'data': outstanding
	}, {
		'label': 'Pending : Rp. ' + addDecimalPoints(pending) + ',00',
		'color': '#FFC107',
		'data': pending
	}, {
		'label': 'Paid : Rp. ' + addDecimalPoints(paid) + ',00',
		'color': '#FF7043',
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
				if (arrData[i].state == "used") {
					used++;
				}

				if (arrData[i].state == "paid") {
					used++;
					paid++;
				}

				var dateValid = new Date(arrData[i].valid_at);
				var dateExpired = new Date(arrData[i].expired_at);
				dataSet[i] = [
					arrData[i].voucher_code
					, arrData[i].holder_description.toUpperCase()
					, dateValid.toDateString().toUpperCase()
					, dateExpired.toDateString().toUpperCase()
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
				"order": [[4, "desc"]],
				columns: [
					{title: "VOUCHER"},
					{title: "HOLDER"},
					{title: "VALID"},
					{title: "EXPIRED"},
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

			getProgram(id, arrData.length, used, paid);
		},
		error: function (data) {
			getProgram(id, 0, 0, 0);
		}
	});
}

function getPartner(id) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/partner/program?program_id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var i;
			var arrData = data.data;
			var limit = arrData.length;

			for (i = 0; i < limit; i++) {
				var sn = arrData[i].serial_number.String;
				if (sn == "") {
					sn = "-";
				}

				var html = "<div class='mda-list-item'>"
					+ "<div class='mda-list-item-icon mt0'><div class='text-lg initial64 bg-blue-500'>" + arrData[i].name.charAt(0).toUpperCase() + "</div></div>"
					+ "<div class='mda-list-item-text'>"
					+ "<h3>" + arrData[i].name + "</h3>"
					+ "<div class='text-muted text-ellipsis'>Serial Number : " + sn + "</div>"
					+ "</div></div>";
				var li = $("<div class='mda-list col-lg-6'></div>").html(html);
				li.appendTo('#list-partner');
			}
		},
		error: function (data) {
			$("<div class='card-body text-center'>No Partner Found</div>").appendTo('#card-partner');
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
					$("#button-voucher").attr("style", "display:none");
					break;
				case 'gift':
					programType = "Gift Voucher"
					$("#button-voucher").attr("style", "display:none");
					break;
			}

			// Program
			$('#program-name').html(result.name);
			$('#program-description').html(result.description);
			$('#program-type').html(programType);
			$('#voucher-type').html(result.voucher_type);
			$('#conversion-rate').html(result.voucher_price);
			$('#voucher-value').html('Rp. ' + addDecimalPoints(result.voucher_value) + ',00');
			$('#period').html(period);
			$('#program-tnc').html(result.tnc);

			// Voucher
			$('#max-quantity-voucher').html(result.max_quantity_voucher);
			$('#remaining-voucher').html(remainingVoucher);
			$('#created-voucher').html(voucher);
			$('#used-voucher').html(used);
			$('#paid-voucher').html(paid);

			// Cashflow
			totalMon = parseInt(result.voucher_value) * parseInt(result.max_quantity_voucher);
			$('#total-money').html('Rp. ' + addDecimalPoints(totalMon) + ',00');

			remainingMon = (parseInt(result.max_quantity_voucher) - parseInt(voucher)) * parseInt(result.voucher_value);
			$('#remaining-money').html('Rp. ' + addDecimalPoints(remainingMon) + ',00');

			outstandingMon = (parseInt(voucher) - parseInt(used)) * parseInt(result.voucher_value);
			$('#remaining-money').html('Rp. ' + addDecimalPoints(remainingMon) + ',00');

			pendingPay = parseInt(result.voucher_value) * parseInt(used);
			$('#pending-payment').html('Rp. ' + addDecimalPoints(pendingPay) + ',00');

			paidPay = parseInt(result.voucher_value) * parseInt(paid);
			$('#paid-payment').html('Rp. ' + addDecimalPoints(paidPay) + ',00');

			remainingPercent = remainingMon / totalMon * 100;
			$('#remaining-bar').attr('style', 'width:' + remainingPercent + '%');
			$('#remaining-bar').attr('data-original-title', 'Remaining : ' + remainingPercent + '%');

			pendingPercent = pendingPay / totalMon * 100;
			$('#pending-bar').attr('style', 'width:' + pendingPercent + '%');
			$('#pending-bar').attr('data-original-title', 'Pending : ' + pendingPercent + '%');

			paidPercent = paidPay / totalMon * 100;
			$('#paid-bar').attr('style', 'width:' + paidPercent + '%');
			$('#paid-bar').attr('data-original-title', 'Paid : ' + paidPercent + '%');

			chart(remainingMon, outstandingMon, pendingPay, paidPay);
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function programCampaign(){
	var id = findGetParameter("id");
	window.location = "/program/campaign?id=" + id;
}

function editProgram() {
	var id = findGetParameter("id");
	window.location = "/program/update?id=" + id;
}
