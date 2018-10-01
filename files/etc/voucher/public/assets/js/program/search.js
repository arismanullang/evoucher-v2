$(document).ready(function () {
  getProgram();
  $(document).ajaxStart(function(){
    // Show image container
    $(".cssload-loader").show();
   });
   $(document).ajaxComplete(function(){
    // Hide image container
    $(".cssload-loader").hide();
   });
});

function getProgram() {
	var arrData = [];
	$.ajax({
		url: '/v1/ui/program/all?token=' + token,
    type: 'get',
    processing: true,
		success: function (data) {
			arrData = data.data;
			var i;
			var dataSet = [];
			var dataId = [];
			var dataType = [];
			var dataStart = [];
			var dataEnd = [];
			var dataModified = [];
			var dataName = [];
			var dataPrice = [];
			var dataValue = [];
			var dataMax = [];
			var dataVoucher = [];
			var dataRedeem = [];
			var dataStatus = [];
			var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
			for (i = 0; i < arrData.length; i++) {
        var startDate = new Date(arrData[i].start_date);
        var endDate = new Date(arrData[i].end_date);
        var createdAt = new Date(arrData[i].created_at);
        var updatedAt = new Date(arrData[i].updated_at.String);

				var date1 = startDate.toString().substring(4, 15).split(" ");
				var date2 = endDate.toString().substring(4, 15).split(" ");
				var date3 = createdAt.toString().substring(4, 15).split(" ");
        var date4 = updatedAt.toString().substring(4, 15).split(" ");

				dataId.push(arrData[i].id);
				if (arrData[i].type == "on-demand") {
					dataType.push("mobile app");
				} else if (arrData[i].type == "gift") {
					dataType.push("gift voucher");
				} else if (arrData[i].type == "privilege") {
          dataType.push("privilege");
				} else {
					dataType.push("Email Blast");
				}
        dataStart.push(date1[1] + " " + date1[0] + " " + date1[2]);
				dataEnd.push(date2[1] + " " + date2[0] + " " + date2[2]);
				dataName.push(arrData[i].name);
				dataPrice.push(arrData[i].voucher_price);
				dataValue.push(arrData[i].voucher_value);
				dataMax.push(arrData[i].max_quantity_voucher);

				var created = 0;
				var redeem = 0;

				if (arrData[i].vouchers != null) {
					for (y = 0; y < arrData[i].vouchers.length; y++) {
						created += parseInt(arrData[i].vouchers[y].voucher);
						if (arrData[i].vouchers[y].state != 'created') {
							redeem += parseInt(arrData[i].vouchers[y].voucher);
						}
					}
				}

				dataVoucher.push(created);
				dataRedeem.push(redeem);

				if (arrData[i].status = 'created') {
          // var dateStart = new Date(date1[0], date1[1] - 1, date1[2]);
					// var dateEnd = new Date(date2[0], date2[1] - 1, date2[2], 23, 59, 59);
          var dateStart = startDate;
          var dateEnd = endDate;
					if (Date.now() < dateStart.getTime()) {
						dataStatus.push("Not Active");
					} else if (Date.now() > dateStart.getTime() && Date.now() < dateEnd.getTime()) {
						dataStatus.push("Active");
					} else if (Date.now() > dateEnd.getTime()) {
						dataStatus.push("End");
					}
				} else if (arrData[i].status = "deleted") {
					dataStatus.push("Disabled");
				}

				if (arrData[i].updated_at.String != "") {
					dataModified.push(date4[1] + " " + date4[0] + " " + date4[2]);
				} else {
					dataModified.push(date3[1] + " " + date3[0] + " " + date3[2]);
				}

			}

			for (i = 0; i < dataId.length; i++) {
				var button = "<button type='button' onclick='detail(\"" + dataId[i] + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>" +
					"<button type='button' onclick='edit(\"" + dataId[i] + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-edit'></em></button>" +
          "<button type='button' value=\"" + dataId[i] + "\" class='btn btn-flat btn-sm btn-danger swal-demo4'><em class='ion-trash-a'></em></button>"

          if(dataType[i] == "privilege"){
            button = "<button type='button' onclick='detail(\"" + dataId[i] + "\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"
          }

				var avail = 0;
				var redemptionRate = 0;
				//var distributionRate = 0;
				if (dataMax[i] != 0) {
					avail = dataMax[i] - dataVoucher[i];
					if (dataVoucher[i] != 0) {
						redemptionRate = dataRedeem[i] / dataVoucher[i] * 100;
						//distributionRate = dataVoucher[i] / dataMax[i] * 100;
					}
				}

				var tempArray = [
					dataName[i].toUpperCase()
					, dataType[i].toUpperCase()
					, dataPrice[i] + " /<br> Rp. " + addDecimalPoints(dataValue[i]) + ",00"
					, dataStatus[i].toUpperCase()
					, dataStart[i].toUpperCase()
					, dataEnd[i].toUpperCase()
					, dataModified[i].toUpperCase()
					, dataMax[i]
					, avail
					, Math.round(redemptionRate) + "%"
					, button
        ];

        var privilegeArray = [
          dataName[i].toUpperCase()
					, dataType[i].toUpperCase()
					, "-"
					, dataStatus[i].toUpperCase()
					, dataStart[i].toUpperCase()
					, dataEnd[i].toUpperCase()
					, dataModified[i].toUpperCase()
					, "-"
					, "-"
					, "-"
					, button
        ];

        if(dataType[i] == "privilege"){
          dataSet.push(privilegeArray);
        } else {
          dataSet.push(tempArray);
        }
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
				"order": [[6, "desc"]],
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

function edit(url) {
	window.location = "/program/update?id=" + url;
}

function detail(url) {
	window.location = "/program/check?id=" + url;
}

function addProgram(url) {
	window.location = "/program/create";
}

function deleteProgram(id) {
	var status = "";
	$.ajax({
		url: '/v1/ui/program/delete?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			getProgram();
			swal('Delete Success!');
		},
		error: function (xhr, ajaxOptions, thrownError) {
			swal('Program Cannot Be Deleted', xhr.responseJSON.errors.detail);
		}
	});
	return status;
}

(function () {
	'use strict';

	$(runSweetAlert);

	//onclick='deleteProgram(\""+arrData[i].Id+"\")'
	function runSweetAlert() {
		$(document).on('click', '.swal-demo4', function (e) {
			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want delete program?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Yes, delete it!',
					closeOnConfirm: false
				},
				function () {
					deleteProgram(e.target.value);
				});

		});
	}

})();
