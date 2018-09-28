$(document).ready(function () {
	getProgram();
});

function getProgram() {
	$.ajax({
		url: '/v1/ui/program/ongoing?token=' + token,
		type: 'get',
		success: function (data) {
			var programList = [];
			var result = data.data;
			console.log(result);
			var limit = 5;

			var totalVoucher = 0;
			var totalIssuedVoucher = 0;
			var totalRedeemedVoucher = 0;
			var totalPaidVoucher = 0;
			var totalPayableVoucher = 0;

			var totalProgram = result.length;
			for (i = 0; i < result.length; i++) {
				var date = result[i].end_date.substring(0, 10).split("-");
				var dateEnd = new Date(date[0], date[1] - 1, date[2]);
				var dateEnd_ms = dateEnd.getTime();
				var dateNow_ms = Date.now();
				var one_day = 1000 * 60 * 60 * 24;
				var diffNow = Math.round((dateEnd_ms - dateNow_ms) / one_day);

				var voucher = result[i].max_quantity_voucher;
				var tempProgramIssued = 0;
				var tempProgramRedeemed = 0;
				var tempProgramPaid = 0;
				var tempProgramVoucherValue = parseInt(result[i].voucher_value);
				// var html = "<h5 class='mb-sm'><a href='/program/check?id=" + result[i].id + "'>" + result[i].name + "</a></h5>"
				// 	+ "<p class='text-muted'>End in " + diffNow + " days</p>";

				var eleVoucher = "";
				if (result[i].partners == null) {
					eleVoucher += "<p>No voucher generated</p>";
				} else {
					for (var y = 0; y < result[i].partners.length; y++) {
						var tempPartner = result[i].partners[y];
						var tempIssued = 0;
						var tempRedeemed = 0;
						var tempPaid = 0;
						var tempUnpaid = 0;

						for (var yy = 0; yy < tempPartner.vouchers.length; yy++) {
							if (tempPartner.vouchers[yy].state == "created") {
								tempProgramIssued += parseInt(tempPartner.vouchers[yy].voucher);

								tempIssued += parseInt(tempPartner.vouchers[yy].voucher);
							} else if (tempPartner.vouchers[yy].state == "used") {
								tempProgramIssued += parseInt(tempPartner.vouchers[yy].voucher);
								tempProgramRedeemed += parseInt(tempPartner.vouchers[yy].voucher);

								tempIssued += parseInt(tempPartner.vouchers[yy].voucher);
								tempRedeemed += parseInt(tempPartner.vouchers[yy].voucher);
							} else {
								tempProgramIssued += parseInt(tempPartner.vouchers[yy].voucher);
								tempProgramRedeemed += parseInt(tempPartner.vouchers[yy].voucher);
								tempProgramPaid += parseInt(tempPartner.vouchers[yy].voucher)

								tempIssued += parseInt(tempPartner.vouchers[yy].voucher);
								tempRedeemed += parseInt(tempPartner.vouchers[yy].voucher);
								tempPaid += parseInt(tempPartner.vouchers[yy].voucher)
							}
						}

						var tempRedeemValue = tempRedeemed * tempProgramVoucherValue;
						var tempUnpaidValue = 0;
						if(tempRedeemed > tempPaid){
							tempUnpaid = tempRedeemed - tempPaid;
							tempUnpaidValue = tempUnpaid * tempProgramVoucherValue;
							totalPayableVoucher += tempUnpaidValue;
						}
						if(tempPartner.name != ""){
							eleVoucher += "<div class='row'>"
								+ "<div class='col-md-4'>Partner : "+tempPartner.name+"</div>"
								+ "<div class='col-md-4'>Used : "+tempRedeemed+" / Rp. "+addDecimalPoints(tempRedeemValue)+",00 </div>"
								+ "<div class='col-md-4'>Unpaid : "+tempUnpaid+" / Rp. "+addDecimalPoints(tempUnpaidValue)+",00 </div>"
								+ "</div>";
						}
					}
				}
				var header = "<div id='heading"+i+"'role='tab' class='panel-heading'>"
					+ "<h4 class='panel-title'>"
					+ "<a role='button' data-toggle='collapse' data-parent='#accordion' href='#collapse"+i+"' aria-expanded='false' aria-controls='collapse"+i+"' class='collapsed'>"
					+ "<div class='row'><div class='col-lg-2'>" +result[i].name+"</div>"
					+ "<div class='col-lg-2'>End in : "+diffNow+" Days</div>"
					+ "<div class='col-lg-2'>Stock : "+(voucher-tempProgramIssued)+"</div>"
					+ "<div class='col-lg-2'>Redeemed : "+tempProgramIssued+"</div>"
					+ "<div class='col-lg-2'>Used : "+tempProgramRedeemed+"</div>"
					+ "<div class='col-lg-2'>Paid : "+tempProgramPaid+"</div>"
					+ "</div></div></a></h4></div>"
					+ "<div id='collapse"+i+"' role='tabpanel' aria-labelledby='heading"+i+"' class='panel-collapse collapse' aria-expanded='false' style='height: 0px;'>"
          + "<div class='panel-body'>";

          if(result[i].type == "privilege"){
              header = "<div id='heading"+i+"'role='tab' class='panel-heading'>"
            + "<h4 class='panel-title'>"
            + "<a role='button' data-toggle='collapse' data-parent='#accordion' href='#collapse"+i+"' aria-expanded='false' aria-controls='collapse"+i+"' class='collapsed'>"
            + "<div class='row'><div class='col-lg-2'>" +result[i].name+"</div>"
            + "<div class='col-lg-2'>End in : -</div>"
            + "<div class='col-lg-2'>Stock : -</div>"
            + "<div class='col-lg-2'>Redeemed : -</div>"
            + "<div class='col-lg-2'>Used : "+tempProgramRedeemed+"</div>"
            + "<div class='col-lg-2'>Paid : -</div>"
            + "</div></div></a></h4></div>"
            + "<div id='collapse"+i+"' role='tabpanel' aria-labelledby='heading"+i+"' class='panel-collapse collapse' aria-expanded='false' style='height: 0px;'>"
            + "<div class='panel-body'>";
          }

				var html = header + eleVoucher;
				html += "</div></div>";

				var li = $("<div class='panel panel-default'></div>").html(html);
				li.appendTo('#accordion');
				totalProgram++;

				totalVoucher += voucher;
				totalIssuedVoucher += tempProgramIssued;
				totalRedeemedVoucher += tempProgramRedeemed;
				totalPaidVoucher += tempProgramPaid;
			}
			$("#totalProgram").html(result.length);
			$("#totalIssued").html(totalIssuedVoucher);
			$("#totalUnpaid").html(totalRedeemedVoucher-totalPaidVoucher);
			$("#payable").html(addDecimalPoints(totalPayableVoucher));

			var stock = totalVoucher-totalIssuedVoucher;
			var outstanding = totalIssuedVoucher-totalRedeemedVoucher;
			var pending = totalRedeemedVoucher-totalPaidVoucher;
			chart(stock,outstanding,pending,totalPaidVoucher);
		},
		error: function (data) {
		}
	});
}

function addProgram() {
	window.location = "/program/create";
}

function chart(remaining, outstanding, pending, paid) {
	var total = remaining + outstanding + pending + paid;
	var pieData = [{
		'label': "<div class='col-md-6'>Stock : </div><div class='col-md-6'>"+remaining+"</div>",
		'color': '#e4eff7',
		'data': remaining
	}, {
		'label': "<div class='col-md-6'>Outstanding : </div><div class='col-md-6'>"+outstanding+"</div>",
		'color': '#f4eff7',
		'data': outstanding
	}, {
		'label': "<div class='col-md-6'>Pending : </div><div class='col-md-6'>"+pending+"</div>",
		'color': '#AFEFEF',
		'data': pending
	}, {
		'label': "<div class='col-md-6'>Paid : </div><div class='col-md-6'>"+paid+"</div>",
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
