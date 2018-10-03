$(window).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	var id = findGetParameter("id");
	searchById(id);
	getPartner(id);
	$("#image-url").change(function () {
		readURL(this);
	});

	$("#all-tenant").change(function () {
		var _this = $(this);
		_this.closest('#partner-list').find("input[class=partner]").prop('checked', _this.prop('checked'));
	});

	$("#allow-accumulative").change(function () {
		if (this.checked == true) {
			$("#div-max-redeem-voucher").attr("style", "display:block");
		} else {
			$("#div-max-redeem-voucher").attr("style", "display:none");
		}
	});

	$('#updateProgram').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			'program-name': {
				required: true
			},
			'max-generate-voucher': {
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
			$("#program-id").val(id);
			$("#program-name").val(program.name);
			$("#program-type").val(program.type);
			$("#voucher-price").val(program.voucher_price);
			$("#voucher-quantity").val(program.max_quantity_voucher);
			$("#generate-voucher").val(program.max_generate_voucher);
      $("#max-redeem-voucher").val(program.max_redeem_voucher);
			$("#limit-redeem-by").val(program.limit_redeem_by);
			$("#redemption-method").val(program.redeem_method);
			$("#program-valid-from").val(convertToDate(program.start_date));
			$("#program-valid-to").val(convertToDate(program.end_date));
			$("#voucher-value").val(program.voucher_value);
			$("#list-rule").html(program.tnc);
			$("#program-description").val(program.description);
			$("#start-hour").val(program.start_hour);
			$("#end-hour").val(program.end_hour);
			$("#image-url-default").val(program.image_url);
			$("#voucher-valid-from").val(program.valid_voucher_start);
			$("#voucher-valid-to").val(program.valid_voucher_end);
			$("#all-tenant").prop("checked", false);
			$("#image-preview").attr("src", program.image_url);
			$("#format").val(program.voucher_format);
			$("#visibility").val(program.visibility);

			$("#program-type").attr("disabled", "");
			$("#voucher-price").attr("disabled", "");
			$("#voucher-quantity").attr("disabled", "");
			$("#voucher-value").attr("disabled", "");
			$("#start-hour").attr("disabled", "");
			$("#end-hour").attr("disabled", "");
			$("#program-valid-from").attr("disabled", "");
			$("#program-valid-to").attr("disabled", "");
			$("#voucher-validity-type").attr("disabled", "");
			$("#voucher-valid-from").attr("disabled", "");
			$("#voucher-valid-to").attr("disabled", "");

			if (program.voucher_lifetime != 0) {
				$("#voucher-lifetime").attr("disabled", "");
				$("#voucher-lifetime").val(program.voucher_lifetime);
				$("#validity-lifetime").attr("style", "display:block");
				$("#validity-date").attr("style", "display:none");
				$("#voucher-validity-type").selectedIndex = 1;
				$("#voucher-validity-type").val("lifetime");
			} else {
				$("#voucher-lifetime").attr("style", "display:none");
				$("#validity-date").attr("disabled", "disabled");
				$("#voucher-valid-from").val(convertToDate(program.valid_voucher_start));
				$("#voucher-valid-to").val(convertToDate(program.valid_voucher_end));
				$("#voucher-validity-type").selectedIndex = 2;
				$("#voucher-validity-type").val("period");
			}

			if ($("#voucher-validity-type").val() == "lifetime") {
				$("#validity-lifetime").attr("style", "display:block");
				$("#validity-date").attr("style", "display:none");
				$("#voucher-valid-from").val("");
				$("#voucher-valid-to").val("");
			} else if ($("#voucher-validity-type").val() == "period") {
				$("#validity-lifetime").attr("style", "display:none");
				$("#validity-date").attr("style", "display:block");
				$("#voucher-lifetime").val("");
			} else {
				$("#validity-lifetime").attr("style", "display:none");
				$("#validity-date").attr("style", "display:none");
				$("#voucher-valid-from").val("");
				$("#voucher-valid-to").val("");
				$("#voucher-lifetime").val("");
			}

			$("#redeem-validity-type").attr("disabled", "");
			if (program.validity_days == "all") {
				$("#validity-day").attr("style", "display:none");
			} else {
				$("#redeem-validity-type").val("selected");
				$("#validity-day").attr("style", "display:block");var i;

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

			$("#program-type").attr("disabled", "");
			if ($("#program-type").val() == "bulk") {
				$("#target").attr("style", "display:block");
				$("#conversion-row").attr("style", "display:none");
				$("#generate-row").attr("style", "display:none");
				$("#voucher-quantity").attr("disabled", "");
				$("#voucher-price").attr("disabled", "");
			} else {
				$("#target").attr("style", "display:none");
				$("#conversion-row").attr("style", "display:block");
				$("#generate-row").attr("style", "display:block");
			}
			if (program.allow_accumulative) {
				$("#allow-accumulative").attr("checked", true);
				$("#use-voucher").attr("style", "display:block");
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
	var programName = $("#program-name").val();
	var programType = $("#program-type").find(":selected").val();
	var voucherPrice = parseInt($("#voucher-price").val());
	var maxQuantityVoucher = parseInt($("#voucher-quantity").val());
	var programValidFrom = $("#program-valid-from").val();
	var programValidTo = $("#program-valid-to").val();
	var startHour = $("#start-hour").val();
	var endHour = $("#end-hour").val();
	var voucherValue = parseInt($("#voucher-value").val());
	var programDescription = $("#program-description").val();
	var visibility = $("#visibility").val();
	var format = $("#format").val();

	var listDay = "";
	if ($("#redeem-validity-type").val() == "all") {
		listDay = "all";
	} else if ($("#redeem-validity-type").val() == "weekend") {
		listDay = "sunday;saturday";
	} else if ($("#redeem-validity-type").val() == "weekday") {
		listDay = "monday;tuesday;wednesday;thursday;friday";
	} else if ($("#redeem-validity-type").val() == "selected") {
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

	if ($("#voucher-validity-type").val() == "period") {
		lifetime = 0;
		periodStart = $("#voucher-valid-from").val();
		periodEnd = $("#voucher-valid-to").val();
	} else if ($("#voucher-validity-type").val() == "lifetime") {
		lifetime = parseInt($("#voucher-lifetime").val());
		periodStart = "01/01/1001";
		periodEnd = "01/01/1001";
	}

	var maxRedeem = parseInt($("#max-redeem-voucher").val());
  var maxGenerate = parseInt($("#generate-voucher").val());

  var limitRedeemBy = $("#limit-redeem-by").find(":selected").val();

	if ($("#program-type").val() == "bulk") {
		maxGenerate = 1;
		maxRedeem = 1;
	}

	var str = $("#list-rule").summernote('code');
	var tnc = str.replace(/^\s+|\s+$|(\r?\n|\r)/g, '');

	if (!str.includes("<p>")) {
		tnc = '<p>' + tnc + '</p>';
	}

	//allow accumulative
	var allowAccumulative = $("#allow-accumulative").is(":checked");
	if (allowAccumulative) {
		maxRedeem = parseInt($("#max-redeem-voucher").val());

		$( "#voucher-lifetime" ).rules( "add", {
			required: true,
			digits: true,
			min: 1
		});
  }

  var programEndDate = new Date(programValidTo);
  programEndDate.setHours(23);
  programEndDate.setMinutes(59);
  programEndDate.setSeconds(59);

  var voucherEndDate = new Date(periodEnd);
  if(periodEnd != "01/01/1970"){
    voucherEndDate.setHours(23);
    voucherEndDate.setMinutes(59);
    voucherEndDate.setSeconds(59);
  }
  var today = dateFormat(new Date(), 'isoUtcDateTime');
  startHour = startHour + ':00' + today.substr(19);
  endHour = endHour + ':00' + today.substr(19);

  //2018-06-29T09:07:51Z
  console.log(startHour);
  console.log(endHour);

	var formData = new FormData();
	var img = $('#image-url-default').val();
	var redeem = $("#redemption-method").val();
	var voucherType = "cash";
	var id = $("#program-id").val();
	if ($('#image-url')[0].files[0] != null) {
		$('#modal-loader').modal({backdrop: 'static', keyboard: false}, 'show');

		formData.append('image-url', $('#image-url')[0].files[0]);
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
          limit_redeem_by: limitRedeemBy,
					allow_accumulative: allowAccumulative,
					redemption_method: redeem,
					start_date: dateFormat(new Date(programValidFrom), 'isoUtcDateTime'),
					end_date: dateFormat(programEndDate, 'isoUtcDateTime'),
					start_hour: startHour,
					end_hour: endHour,
					voucher_value: voucherValue,
					image_url: img,
					tnc: tnc,
					description: programDescription,
					validity_days: listDay,
					valid_voucher_start: dateFormat(new Date(periodStart), 'isoUtcDateTime'),
					valid_voucher_end: dateFormat(voucherEndDate, 'isoUtcDateTime'),
					voucher_lifetime: lifetime,
					visibility: visibility,
					voucher_format: format
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
      limit_redeem_by: limitRedeemBy,
			allow_accumulative: allowAccumulative,
			redemption_method: redeem,
			start_date: dateFormat(new Date(programValidFrom), 'isoUtcDateTime'),
			end_date: dateFormat(programEndDate, 'isoUtcDateTime'),
			start_hour: startHour,
			end_hour: endHour,
			voucher_value: voucherValue,
			image_url: img,
			tnc: tnc,
			description: programDescription,
			validity_days: listDay,
			valid_voucher_start: dateFormat(new Date(periodStart), 'isoUtcDateTime'),
			valid_voucher_end: dateFormat(voucherEndDate, 'isoUtcDateTime'),
			voucher_lifetime: lifetime,
			visibility: visibility,
			voucher_format: format
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
					+ "<input type='checkbox' id='chk_partner" + i +"' class='partner'>"
					+ "<span class='ion-checkmark-round'></span>" + arrData[i].name
					+ "</label>";
				li.html(html);
        li.appendTo('#partner-list');
        $('#chk_partner'+i).val(arrData[i].id);
      }

			$.ajax({
				url: '/v1/ui/partner/program?program_id=' + id + '&token=' + token,
				type: 'get',
				success: function (data) {
					var i;
					var y;
					var inputPartner = $(".partner");

					for (i = 0; i < inputPartner.length; i++) {
						var tempElem = inputPartner[i];
						var arrData = data.data;
						var limit = arrData.length;
						for (y = 0; y < limit; y++) {
							if (tempElem.value == arrData[y].id) {
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
  // var string1 = date.split("T")[0];
	// var string2 = string1.split("-");
	// var result = string2[1] + "/" + string2[2] + "/" + string2[0];
  var newDate = new Date(date)

  var dd = newDate.getDate();
  var mm = newDate.getMonth()+1; //January is 0!

  var yyyy = newDate.getFullYear();
  if(dd<10){
      dd='0'+dd;
  }
  if(mm<10){
      mm='0'+mm;
  }
	var result = mm + "/" + dd + "/" + yyyy;

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
	}

})();
