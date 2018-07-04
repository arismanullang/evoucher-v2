$(document).ready(function () {
	getTransactionByPartner("");
  getPartner();
  $(document).ajaxStart(function(){
    // Show image container
    $(".cssload-loader").show();
   });
   $(document).ajaxComplete(function(){
    // Hide image container
    $(".cssload-loader").hide();
   });
});

function getPartner() {
	$.ajax({
		url: '/v1/ui/partner/all?token=' + token,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<option value='" + arrData[i].name + "'>" + arrData[i].name + "</div>");
				li.appendTo('#partner-id');
			}
		}
	});
}

function getTransactionByPartner(partnerId) {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/transaction/partner?token=' + token + '&partner=' + partnerId,
		type: 'get',
		success: function (data) {
			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var arrData = data.data;
			var i;
			var dataSet = [];
			var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
			var username = [];
			var usernameExist = false;

			for (i = 0; i < arrData.length; i++) {
				usernameExist = false;
				var date1 = arrData[i].issued.substring(0, 10).split("-");
				var date2 = arrData[i].redeemed.substring(0, 10).split("-");
				var date3 = arrData[i].cashout.String.substring(0, 10).split("-");
				console.log(date3);
				console.log(arrData[i].cashout.String);
				var cashoutDate = date3[2] + " " + months[parseInt(date3[1]) - 1] + " " + date3[0];
				var cashoutCashier = arrData[i].username.String;

				if (username.length == 0) {
					usernameExist = true;
				} else {
					for (var y = 0; y < username.length; y++) {
						if (username[y] == arrData[i].username.String.toUpperCase()) {
							usernameExist = false;
							break;
						}
						usernameExist = true;
					}
				}

				if (usernameExist) {
					username.push(arrData[i].username.String.toUpperCase());
				}
				var tempVoucherLength = arrData[i].vouchers.length;
				for (y = 0; y < tempVoucherLength; y++) {
					var status = "Paid"
					if (arrData[i].vouchers[y].state == "used") {
						cashoutDate = "-";
						cashoutCashier = "-";
						status = "Pending";
					}else if (arrData[i].vouchers[y].state == "created"){
						cashoutDate = "-";
						cashoutCashier = "-";
						status = "Issued";
					}

					var tempArray = [
						arrData[i].partner_name.toUpperCase()
						, arrData[i].transaction_code
						, arrData[i].program_name.toUpperCase()
						, arrData[i].vouchers[y].voucher_code
						, "Rp." + addDecimalPoints(arrData[i].voucher_value).trim() + ",00"
						, date1[2] + " " + months[parseInt(date1[1]) - 1] + " " + date1[0]
						, date2[2] + " " + months[parseInt(date2[1]) - 1] + " " + date2[0]
						, cashoutDate
						, cashoutCashier.toUpperCase()
						, status.toUpperCase()
					];

					dataSet.push(tempArray);
				}
				i += tempVoucherLength - 1;
			}

			for (y = 0; y < username.length; y++) {
				var li = $("<option value='" + username[y] + "'>" + username[y] + "</div>");
				li.appendTo('#username');
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'lBrtip',
				buttons: [
					'copy', 'csv', 'excel', 'pdf', 'print'
				],
				"order": [[9, "desc"]],
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
