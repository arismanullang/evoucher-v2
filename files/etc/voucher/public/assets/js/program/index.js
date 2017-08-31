$(document).ready(function () {
	getProgram();
});

function getProgram() {
	$.ajax({
		url: '/v1/ui/program/all?token=' + token,
		type: 'get',
		success: function (data) {
			var programList = [];
			var result = data.data;
			var limit = 5;

			var totalVoucher = 0;
			var totalUsedVoucher = 0;
			var totalGeneratedVoucher = 0;
			var totalProgram = 0;
			for (i = 0; i < result.length; i++) {
				if (totalProgram < 5) {

					var date = result[i].end_date.substring(0, 10).split("-");
					var dateEnd = new Date(date[0], date[1] - 1, date[2]);
					var dateEnd_ms = dateEnd.getTime();
					var dateNow_ms = Date.now();
					var one_day = 1000 * 60 * 60 * 24;
					var diffNow = Math.round((dateEnd_ms - dateNow_ms) / one_day);

					var html = "<h5 class='mb-sm'><a href='/program/check?id=" + result[i].id + "'>" + result[i].name + "</a></h5>"
						+ "<p class='text-muted'>End in " + diffNow + " days</p>";
					if (result[i].vouchers == null) {
						html += "<p>No voucher generated</p>";
					} else {
						var voucher = 0;
						for (var y = 0; y < result[i].vouchers.length; y++) {
							voucher += parseInt(result[i].vouchers[y].voucher);

							if (result[i].vouchers[y].state != "created") {
								totalUsedVoucher += parseInt(result[i].vouchers[y].voucher);
							} else {
								totalGeneratedVoucher += parseInt(result[i].vouchers[y].voucher);
							}
						}
						html += "<p>" + voucher + " vouchers have distributed. " + (result[i].max_quantity_voucher - voucher) + " vouchers left.</p>";
					}
					var li = $("<li class='list-group-item'></li>").html(html);
					if (dateEnd_ms > dateNow_ms) {
						li.appendTo('#upcomming-program');
						totalProgram++;
					}
				}

				if (!(programList.includes(result[i].id))) {
					programList.push(result[i].id);
					totalVoucher += parseInt(result[i].max_quantity_voucher);
				}

			}

			$("#total-program").html(result.length);
			$("#total-voucher").html(totalVoucher);
			$("#total-generated").html(totalGeneratedVoucher);
			$("#total-used").html(totalUsedVoucher);
		},
		error: function (data) {
			swal("Error", "Program Not Found.");
		}
	});
}

function addProgram() {
	window.location = "/program/create";
}
