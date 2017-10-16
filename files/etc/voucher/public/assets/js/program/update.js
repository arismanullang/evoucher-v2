$(window).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	var id = findGetParameter("id");
	searchById(id);
	getPartner(id);
	$("#imageUrl").change(function () {
		readURL(this);
	});

	$("#allTenant").change(function () {
		var _this = $(this);
		_this.closest('#partnerList').find("input[class=partner]").prop('checked', _this.prop('checked'));
	});

	$("#allowAccumulative").change(function () {
		if (this.checked == true) {
			$("#divMaxRedeemVoucher").attr("style", "display:block");
		} else {
			$("#divMaxRedeemVoucher").attr("style", "display:none");
		}
	});

	$('#updateProgram').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			programName: {
				required: true
			},
			maxGenerateVoucher: {
				required: true,
				digits: true,
				min: 1
			},
			'partner[]': {
				required: true
			}
		}
	});
});

function readURL(input) {
	if (input.files && input.files[0]) {
		var reader = new FileReader();
		reader.onload = function (e) {
			$('#imagePreview').attr('src', e.target.result);
		};

		reader.readAsDataURL(input.files[0]);
	}
}

function searchById(id) {

	var arrData = [];

	$.ajax({
		url: '/v1/ui/program/detail?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			console.log(data);
			var program = data.data[0];
			$("#programId").val(id);
			$("#programName").val(program.name);
			$("#programType").val(program.type);
			$("#voucherPrice").val(program.voucher_price);
			$("#maxQuantityVoucher").val(program.max_quantity_voucher);
			$("#maxGenerateVoucher").val(program.max_generate_voucher);
			$("#maxRedeemVoucher").val(program.max_redeem_voucher);
			$("#redemptionMethod").val(program.redeem_method);
			$("#programValidFrom").val(convertToDate(program.start_date));
			$("#programValidTo").val(convertToDate(program.end_date));
			$("#voucherValue").val(program.voucher_value);
			$("#listRule").html(program.tnc);
			$("#programDescription").val(program.description);
			$("#startHour").val(program.start_hour);
			$("#endHour").val(program.end_hour);
			$("#imageUrlDefault").val(program.image_url);
			$("#voucherValidFrom").val(program.valid_voucher_start);
			$("#voucherValidTo").val(program.valid_voucher_end);
			$("#allTenant").prop("checked", false);
			$("#imagePreview").attr("src", program.image_url);

			$("#programType").attr("disabled", "");
			$("#voucherPrice").attr("disabled", "");
			$("#maxQuantityVoucher").attr("disabled", "");
			$("#voucherValue").attr("disabled", "");
			$("#startHour").attr("disabled", "");
			$("#endHour").attr("disabled", "");
			$("#programValidFrom").attr("disabled", "");
			$("#programValidTo").attr("disabled", "");
			$("#voucherValidityType").attr("disabled", "");
			$("#voucherValidFrom").attr("disabled", "");
			$("#voucherValidTo").attr("disabled", "");

			if (program.voucher_lifetime != 0) {
				$("#voucherLifetime").attr("disabled", "");
				$("#voucherLifetime").val(program.voucher_lifetime);
				$("#validityLifetime").attr("style", "display:block");
				$("#validityDate").attr("style", "display:none");
				$("#voucherValidTrom").val("");
				$("#voucherValidTo").val("");
				$("#voucherValidityType").selectedIndex = 1;
				$("#voucherValidityType").val("lifetime");
			}
			if ($("#voucherValidityType").val() == "lifetime") {
				$("#validityLifetime").attr("style", "display:block");
				$("#validityDate").attr("style", "display:none");
				$("#voucherValidFrom").val("");
				$("#voucherValidTo").val("");
			} else if ($("#voucherValidityType").val() == "period") {
				$("#validityLifetime").attr("style", "display:none");
				$("#validityDate").attr("style", "display:block");
				$("#voucherLifetime").val("");
			} else {
				$("#validityLifetime").attr("style", "display:none");
				$("#validityDate").attr("style", "display:none");
				$("#voucherValidFrom").val("");
				$("#voucherValidTo").val("");
				$("#voucherLifetime").val("");
			}

			$("#redeemValidityType").attr("disabled", "");
			if (program.validity_days == "all") {
				$("#validityDay").attr("style", "display:none");
			} else {
				$("#redeemValidityType").val("selected");
				$("#validityDay").attr("style", "display:block");var i;

				var y;
				var li = $("input[class=days]");
				var arrData = program.validity_days.split(";");

				for (i = 0; i < li.length; i++) {
					var tempElem = li[i];
					var limit = arrData.length;
					for (y = 0; y < limit; y++) {
						if (tempElem.getAttribute("value") == arrData[y]) {
							console.log(tempElem.getAttribute("value") + " " + arrData[y]);
							console.log(tempElem);

							tempElem.checked = true;
						}
					}
				}
			}

			$("#programType").attr("disabled", "");
			if ($("#programType").val() == "bulk") {
				$("#target").attr("style", "display:block");
				$("#conversionRow").attr("style", "display:none");
				$("#generateRow").attr("style", "display:none");
				$("#maxQuantityVoucher").attr("disabled", "");
				$("#voucherPrice").attr("disabled", "");
			} else {
				$("#target").attr("style", "display:none");
				$("#conversionRow").attr("style", "display:block");
				$("#generateRow").attr("style", "display:block");
			}
			if (program.allow_accumulative) {
				$("#allowAccumulative").attr("checked", true);
				$("#useVoucher").attr("style", "display:block");
			}

			$(".select2").select2();
			$('.summernote').each(function () {
				$(this).summernote({
					height: 380,
					placeholder: 'Any Message...',
					callbacks: {
						onPaste: function (e) {
							var bufferText = ((e.originalEvent || e).clipboardData || window.clipboardData).getData('Text');

							e.preventDefault();

							// Firefox fix
							setTimeout(function () {
								document.execCommand('insertText', false, bufferText);
							}, 10);
						}
					}
				});
			});
		}
	});
}

function send() {
	var programName = $("#programName").val();
	var programType = $("#programType").find(":selected").val();
	var voucherPrice = parseInt($("#voucherPrice").val());
	var maxQuantityVoucher = parseInt($("#maxQuantityVoucher").val());
	var programValidFrom = $("#programValidFrom").val();
	var programValidTo = $("#programValidTo").val();
	var startHour = $("#startHour").val();
	var endHour = $("#endHour").val();
	var voucherValue = parseInt($("#voucherValue").val());
	var programDescription = $("#programDescription").val();

	var listDay = "";
	if ($("#redeemValidityType").val() == "all") {
		listDay = "all";
	} else if ($("#redeemValidityType").val() == "weekend") {
		listDay = "sunday;saturday";
	} else if ($("#redeemValidityType").val() == "weekday") {
		listDay = "monday;tuesday;wednesday;thursday;friday";
	} else if ($("#redeemValidityType").val() == "selected") {
		var li = $("input[class=days]:checked");

		$( ".days" ).rules( "add", {
			required: true
		});


		if (li.length > 0) {
			for (i = 0; i < li.length; i++) {
				if (li[i].value != "on") {
					listDay = listDay + li[i].value + ";";
				}
			}
		}
	}

	var listPartner = [];
	var li = $("input[class=partner]:checked");

	for (i = 0; i < li.length; i++) {
		listPartner[i] = li[i].value;
	}

	var lifetime = 0;
	var periodStart = "";
	var periodEnd = "";

	if ($("#voucherValidityType").val() == "period") {
		lifetime = 0;
		periodStart = $("#voucherValidFrom").val();
		periodEnd = $("#voucherValidTo").val();
	} else if ($("#voucherValidityType").val() == "lifetime") {
		lifetime = parseInt($("#voucherLifetime").val());
		periodStart = "1001-01-01T00:00:00Z";
		periodEnd = "1001-01-01T00:00:00Z";
	}

	var maxRedeem = parseInt($("#maxRedeemVoucher").val());
	var maxGenerate = parseInt($("#maxGenerateVoucher").val());

	if ($("#programType").val() == "bulk") {
		maxGenerate = 1;
		maxRedeem = 1;
	}

	var str = $("#listRule").summernote('code');
	var tnc = str.replace(/^\s+|\s+$|(\r?\n|\r)/g, '');

	if (!str.includes("<p>")) {
		tnc = '<p>' + tnc + '</p>';
	}

	//allow accumulative
	var allowAccumulative = $("#allowAccumulative").is(":checked");
	if (allowAccumulative) {
		maxRedeem = parseInt($("#maxRedeemVoucher").val());

		$( "#voucherLifetime" ).rules( "add", {
			required: true,
			digits: true,
			min: 1
		});
	}

	var formData = new FormData();
	var img = $('#imageUrlDefault').val();
	var redeem = $("#redemptionMethod").val();
	var voucherType = "cash";
	var id = $("#programId").val();
	if ($('#imageUrl')[0].files[0] != null) {
		$('#modal-loader').modal({backdrop: 'static', keyboard: false}, 'show');

		formData.append('image-url', $('#imageUrl')[0].files[0]);
		jQuery.ajax({
			url: '/file/upload',
			type: "POST",
			processData: false,
			contentType: false,
			data: formData,
			success: function (data) {
				console.log(data.data);
				img = data.data;
				var voucherType = "cash";
				var program = {
					name: programName,
					type: programType,
					voucher_type: voucherType,
					voucher_price: voucherPrice,
					max_quantity_voucher: maxQuantityVoucher,
					max_redeem_voucher: maxRedeem,
					max_generate_voucher: maxGenerate,
					allow_accumulative: allowAccumulative,
					redemption_method: redeem,
					start_date: programValidFrom,
					end_date: programValidTo,
					start_hour: startHour,
					end_hour: endHour,
					voucher_value: voucherValue,
					image_url: img,
					tnc: tnc,
					description: programDescription,
					validity_days: listDay,
					valid_voucher_start: periodStart,
					valid_voucher_end: periodEnd,
					voucher_lifetime: lifetime
				};

				console.log(program);

				$.ajax({
					url: '/v1/ui/program/update?id=' + id + '&type=detail&token=' + token,
					type: 'post',
					dataType: 'json',
					contentType: "application/json",
					data: JSON.stringify(program),
					success: function () {
						var partner = {
							user: "user",
							data: listPartner
						};

						$.ajax({
							url: '/v1/ui/program/update?id=' + id + '&type=tenant&token=' + token,
							type: 'post',
							dataType: 'json',
							contentType: "application/json",
							data: JSON.stringify(partner),
							success: function () {
								swal({
										title: 'Success',
										text: 'Program Updated',
										type: 'success',
										showCancelButton: false,
										confirmButtonText: 'Ok',
										closeOnConfirm: false
									},
									function() {
										var id = findGetParameter("id");
										window.location = "/program/check?id=" + id;
									});
							}
						});
					},
					error: function (data) {
						var a = JSON.parse(data.responseText);
						swal("Error", a.errors.detail);
					}
				});
			},
			error: function (data) {
				var a = JSON.parse(data.responseText);
				swal("Error", a.errors.detail);
			}
		});
	} else {
		var program = {
			name: programName,
			type: programType,
			voucher_type: voucherType,
			voucher_price: voucherPrice,
			max_quantity_voucher: maxQuantityVoucher,
			max_redeem_voucher: maxRedeem,
			max_generate_voucher: maxGenerate,
			allow_accumulative: allowAccumulative,
			redemption_method: redeem,
			start_date: programValidFrom,
			end_date: programValidTo,
			start_hour: startHour,
			end_hour: endHour,
			voucher_value: voucherValue,
			image_url: img,
			tnc: tnc,
			description: programDescription,
			validity_days: listDay,
			valid_voucher_start: periodStart,
			valid_voucher_end: periodEnd,
			voucher_lifetime: lifetime
		};

		$.ajax({
			url: '/v1/ui/program/update?id=' + id + '&type=detail&token=' + token,
			type: 'post',
			dataType: 'json',
			contentType: "application/json",
			data: JSON.stringify(program),
			success: function () {
				var partner = {
					user: "user",
					data: listPartner
				};

				$.ajax({
					url: '/v1/ui/program/update?id=' + id + '&type=tenant&token=' + token,
					type: 'post',
					dataType: 'json',
					contentType: "application/json",
					data: JSON.stringify(partner),
					success: function () {
						swal({
								title: 'Success',
								text: 'Program Updated',
								type: 'success',
								showCancelButton: false,
								confirmButtonText: 'Ok',
								closeOnConfirm: false
							},
							function() {
								var id = findGetParameter("id");
								window.location = "/program/check?id=" + id;
							});
					}
				});
			},
			error: function (data) {
				var a = JSON.parse(data.responseText);
				swal("Error", a.errors.detail);
			}
		});
	}
}

function getPartner(id) {
	console.log("Get Partner Data");

	$.ajax({
		url: '/v1/ui/partner/all?token=' + token,
		type: 'get',
		success: function (data) {
			console.log("Render Data");
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<div class='col-sm-4'></div>");
				var html = "<label class='checkbox-inline c-checkbox'>"
					+ "<input type='checkbox' class='partner' value='" + arrData[i].id + "' text='" + arrData[i].name + "'>"
					+ "<span class='ion-checkmark-round'></span>" + arrData[i].name
					+ "</label>";
				li.html(html);
				li.appendTo('#partnerList');
			}

			$.ajax({
				url: '/v1/ui/partner/program?program_id=' + id + '&token=' + token,
				type: 'get',
				success: function (data) {
					var i;
					var y;
					var li = $("input[type=checkbox]");

					for (i = 0; i < li.length; i++) {
						var tempElem = li[i];
						var arrData = data.data;
						var limit = arrData.length;
						for (y = 0; y < limit; y++) {
							if (tempElem.getAttribute("text") == arrData[y].name) {
								tempElem.checked = true;
							}
						}
					}
				},
				error: function (data) {
					console.log(data.data);
					$("<div class='card-body text-center'>No Partner Found</div>").appendTo('#cardPartner');
				}
			});
		}
	});
}

function convertToDate(date) {
	var string1 = date.split("T")[0];
	var string2 = string1.split("-");
	var result = string2[1] + "/" + string2[2] + "/" + string2[0];

	return result;
}

function convertToUpperCase(upperCase) {
	var result = "";
	var firstChar = upperCase.charAt(0);
	upperCase = upperCase.replace(firstChar, firstChar.toUpperCase());
	result = upperCase;

	return result;
}

(function () {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$('.select2').select2();
		$("#collapseThree").removeClass("in");
		$("#collapseTwo").removeClass("in");
		$("#collapseFour").removeClass("in");
		$('.datepicker4').datepicker({
			container: '#example-datepicker-container-4',
			autoclose: true,
			startDate: 'd',
			setDate: new Date()
		});

		$('.datepicker3').datepicker({
			container: '#example-datepicker-container-3',
			autoclose: true,
			startDate: 'd',
			setDate: new Date()
		});
		$('#startDate').datepicker('update', new Date());
		$('#endDate').datepicker('update', '+1d');

		var cpInput = $('.clockpicker').clockpicker();
		// auto close picker on scroll
		$('main').scroll(function () {
			cpInput.clockpicker('hide');
		});
	}

})();
