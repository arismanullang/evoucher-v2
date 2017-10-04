$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	getPartner();
	initForm();
	onChangeElem();


});

function initForm(){
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

	if ($("#redeemValidityType").val() == "all") {
		$("#validityDay").attr("style", "display:none");
	} else if ($("#redeemValidityType").val() == "selected") {
		$("#validityDay").attr("style", "display:block");
	} else {
		$("#validityDay").attr("style", "display:none");
	}

	if ($("#programType").val() == "bulk") {
		$("#target").attr("style", "display:block");
		$("#conversionRow").attr("style", "display:none");
		$(".distribution").attr("style", "display:none");
		$("#voucherPrice").val(0);
		$("#voucherPrice").attr("disabled", true);
	} else {
		$("#target").attr("style", "display:none");
		$("#conversionRow").attr("style", "display:block");
		$(".distribution").attr("style", "display:block");
		$("#voucherPrice").attr("disabled", false);
	}
	if ($("#allowAccumulative").is(":checked")) {
		$("#useVoucher").attr("style", "display:block");
	} else {
		$("#useVoucher").attr("style", "display:none");
	}
}

function onChangeElem(){
	$("#allTenant").change(function () {
		var _this = $(this);
		_this.closest('#partnerList').find("input.partner").prop('checked', _this.prop('checked'));
	});
	$("#voucherValidityType").change(function () {
		if (this.value == "lifetime") {
			$("#validityLifetime").attr("style", "display:block");
			$("#validityDate").attr("style", "display:none");
			$("#voucherValidFrom").val("");
			$("#voucherValidTo").val("");
		} else if (this.value == "period") {
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
	});
	$("#redeemValidityType").change(function () {
		if (this.value == "all") {
			$("#validityDay").attr("style", "display:none");
		} else if (this.value == "selected") {
			$("#validityDay").attr("style", "display:block");
		} else {
			$("#validityDay").attr("style", "display:none");
		}
	});
	$("#programType").change(function () {
		if (this.value == "bulk") {
			$("#target").attr("style", "display:block");
			$("#conversionRow").attr("style", "display:none");
			$(".distribution").attr("style", "display:none");
			$("#voucherPrice").val(0);
			$("#voucherPrice").attr("disabled", true);
		} else {
			$("#target").attr("style", "display:none");
			$("#conversionRow").attr("style", "display:block");
			$(".distribution").attr("style", "display:block");
			$("#voucherPrice").val("");
			$("#voucherPrice").attr("disabled", false);
		}
	});
	$("#allowAccumulative").change(function () {
		if (this.checked == true) {
			$("#divMaxRedeemVoucher").attr("style", "display:block");
		} else {
			$("#divMaxRedeemVoucher").attr("style", "display:none");
		}
	});
	$("#imageUrl").change(function () {
		readURL(this);
	});
}

function readURL(input) {
	if (input.files && input.files[0]) {
		var reader = new FileReader();
		reader.onload = function (e) {
			$('#imagePreview').attr('src', e.target.result);
		}

		reader.readAsDataURL(input.files[0]);
	}
}

function download() {
	window.location = "/v1/ui/sample/link?token="+token;
}

function send() {
	var programName = $("#programName").val();
	var programType = $("#programType").find(":selected").val();
	var voucherPrice = parseInt($("#voucherPrice").val());
	var maxQuantityVoucher = parseInt($("#voucherQuantity").val());
	var redemptionMethod = $("#redemptionMethod").find(":selected").val();
	var programValidFrom = $("#programValidFrom").val();
	var programValidTo = $("#programValidTo").val();
	var startHour = $("#startHour").val();
	var endHour = $("#endHour").val();
	var voucherValue = parseInt($("#voucherValue").val());
	var programDescription = $("#programDescription").val();

	// valid days
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

	// partner
	var listPartner = [];
	var li = $("input[class=partner]:checked");

	for (i = 0; i < li.length; i++) {
		if (li[i].value != "on") {
			listPartner[i] = li[i].value;
		}
	}

	// expired
	var lifetime = 0;
	var periodStart = "";
	var periodEnd = "";

	if ($("#voucherValidityType").val() == "period") {
		lifetime = 0;
		periodStart = $("#voucherValidFrom").val();
		periodEnd = $("#voucherValidTo").val();
	} else if ($("#voucherValidityType").val() == "lifetime") {
		$( "#voucherLifetime" ).rules( "add", {
			required: true,
			min: 1,
			max: 1800,
			digits: true
		});

		lifetime = parseInt($("#voucherLifetime").val());

		periodStart = "01/01/0001";
		periodEnd = "01/01/0001";
	}

	// voucher format
	var voucherFormat = {
		prefix: $("#prefix").val(),
		postfix: "",
		body: "",
		format_type: $("#voucherFormat").find(":selected").val(),
		length: 5
	};


	// tnc
	var str = $("#listRule").summernote('code');
	var tnc = str.replace(/^\s+|\s+$|(\r?\n|\r)/g, '');
	tnc = str.replace('<p>', '');
	tnc = str.replace('</p>', '');
	tnc = '<p>' + tnc + '</p>';

	// max generate and redeem
	var maxGenerate = parseInt($("#generateVoucher").val());
	var maxRedeem = 1;

	// voucher type
	var voucherType = "cash";

	// allow accumulative
	var allowAccumulative = $("#allowAccumulative").is(":checked");
	if (allowAccumulative) {
		maxRedeem = parseInt($("#maxRedeemVoucher").val());

		$( "#voucherLifetime" ).rules( "add", {
			required: true,
			digits: true,
			min: 1
		});
	}

	if(!$("#createProgram").valid()) {
		return
	}

	// image
	var formData = new FormData();
	var img = "https://storage.googleapis.com/e-voucher/L1LXN5bpMphnvG6Ce8eUbBSYDW5G3MaH.jpg";
	if ($('#imageUrl')[0].files[0] != null) {
		$('#modalLoader').modal({backdrop: 'static',keyboard: false}, 'show');
		formData.append('image-url', $('#imageUrl')[0].files[0]);

		jQuery.ajax({
			url: '/file/upload',
			type: "POST",
			processData: false,
			contentType: false,
			data: formData,
			success: function (data) {
				img = data.data;
				var program = {
					name: programName,
					type: programType,
					voucher_format: voucherFormat,
					voucher_type: voucherType,
					voucher_price: voucherPrice,
					max_quantity_voucher: maxQuantityVoucher,
					max_redeem_voucher: maxRedeem,
					max_generate_voucher: maxGenerate,
					allow_accumulative: allowAccumulative,
					redemption_method: redemptionMethod,
					start_date: programValidFrom,
					end_date: programValidTo,
					start_hour: startHour,
					end_hour: endHour,
					voucher_value: voucherValue,
					image_url: img,
					tnc: tnc,
					description: programDescription,
					validity_days: listDay,
					valid_partners: listPartner,
					valid_voucher_start: periodStart,
					valid_voucher_end: periodEnd,
					voucher_lifetime: lifetime
				};

				$.ajax({
					url: '/v1/ui/program/create?token=' + token,
					type: 'post',
					dataType: 'json',
					contentType: "application/json",
					data: JSON.stringify(program),
					success: function (data) {
						if ($("#programType").find(":selected").val() == "bulk") {

							var targets = new FormData();
							targets.append('list-target', $("#listTarget")[0].files[0]);

							jQuery.ajax({
								url: '/v1/ui/user/create/broadcast?token=' + token + '&program-id=' + data.data,
								type: "POST",
								processData: false,
								contentType: false,
								data: targets,
								success: function (data) {
									$('#modal-loader').modal('hide');
									swal({
											title: 'Success',
											text: 'Program Created',
											type: 'success',
											showCancelButton: false,
											confirmButtonText: 'Ok',
											closeOnConfirm: false
										},
										function() {
											window.location = "/program/search";
										});
								},
								error: function (data) {
									var a = JSON.parse(data.responseText);
									swal("Error", a.errors.detail);
								}
							});

						} else {
							$('#modal-loader').modal('hide');
							swal({
									title: 'Success',
									text: 'Program Created',
									type: 'success',
									showCancelButton: false,
									confirmButtonText: 'Ok',
									closeOnConfirm: false
								},
								function() {
									window.location = "/program/search";
								});
						}
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
			voucher_format: voucherFormat,
			voucher_type: voucherType,
			voucher_price: voucherPrice,
			max_quantity_voucher: maxQuantityVoucher,
			max_redeem_voucher: maxRedeem,
			max_generate_voucher: maxGenerate,
			allow_accumulative: allowAccumulative,
			redemption_method: redemptionMethod,
			start_date: programValidFrom,
			end_date: programValidTo,
			start_hour: startHour,
			end_hour: endHour,
			voucher_value: voucherValue,
			image_url: img,
			tnc: tnc,
			description: programDescription,
			validity_days: listDay,
			valid_partners: listPartner,
			valid_voucher_start: periodStart,
			valid_voucher_end: periodEnd,
			voucher_lifetime: lifetime
		};

		$.ajax({
			url: '/v1/ui/program/create?token=' + token,
			type: 'post',
			dataType: 'json',
			contentType: "application/json",
			data: JSON.stringify(program),
			success: function (data) {
				if ($("#programType").find(":selected").val() == "bulk") {

					var targets = new FormData();
					targets.append('list-target', $("#listTarget")[0].files[0]);

					jQuery.ajax({
						url: '/v1/ui/user/create/broadcast?token=' + token + '&program-id=' + data.data,
						type: "POST",
						processData: false,
						contentType: false,
						data: targets,
						success: function (data) {
							$('#modal-loader').modal('hide');
							swal({
									title: 'Success',
									text: 'Program Created',
									type: 'success',
									showCancelButton: false,
									confirmButtonText: 'Ok',
									closeOnConfirm: false
								},
								function() {
									window.location = "/program/search";
								});
						},
						error: function (data) {
							var a = JSON.parse(data.responseText);
							swal("Error", a.errors.detail);
						}
					});
				} else {
					$('#modal-loader').modal('hide');
					swal({
							title: 'Success',
							text: 'Program Created',
							type: 'success',
							showCancelButton: false,
							confirmButtonText: 'Ok',
							closeOnConfirm: false
						},
						function() {
							window.location = "/program/search";
						});
				}
			},
			error: function (data) {
				var a = JSON.parse(data.responseText);
				swal("Error", a.errors.detail);
			}
		});
	}
}

function getPartner() {
	$.ajax({
		url: '/v1/ui/partner/all?token=' + token,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<div class='col-sm-4'></div>");
				var html = "<label class='checkbox-inline c-checkbox'>"
					+ "<input type='checkbox' name='partner' class='partner' value='" + arrData[i].id + "'>"
					+ "<span class='ion-checkmark-round'></span>" + arrData[i].name
					+ "</label>";
				li.html(html);
				li.appendTo('#partnerList');
			}
		}
	});
}

// function removeElem(elem){
//   console.log("remove");
//   $(elem).parent().closest('tr').remove();
// }

(function () {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$(".select2").select2();
		$('.datepickerProgramFrom').datepicker({
			container: '#datepickerProgramFrom',
			autoclose: true,
			startDate: 'd',
			setDate: new Date()
		}).on('changeDate', function (selected) {
			var minDate = new Date(selected.date.valueOf());
			$('.datepickerProgramTo').datepicker('setStartDate', minDate);
			$('.datepickerVoucherFrom').datepicker('setStartDate', minDate);
		});
		$('.datepickerProgramTo').datepicker({
			container: '#datepickerProgramTo',
			autoclose: true,
			startDate: '+1d',
			setDate: new Date()
		});
		$('.datepickerVoucherFrom').datepicker({
			container: '#datepickerVoucherFrom',
			autoclose: true,
			startDate: 'd',
			setDate: new Date()
		}).on('changeDate', function (selected) {
			var minDate = new Date(selected.date.valueOf());
			$('.datepickerVoucherTo').datepicker('setStartDate', minDate);
		});
		$('.datepickerVoucherTo').datepicker({
			container: '#datepickerVoucherTo',
			autoclose: true,
			startDate: '+1d',
			setDate: new Date()
		});

		var cpInput = $('.clockpicker').clockpicker();
		// auto close picker on scroll
		$('main').scroll(function () {
			cpInput.clockpicker('hide');
		});

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

		// $('#form-example').validate({
		// 	errorPlacement: errorPlacementInput,
		// 	// Form rules
		// 	rules: {
		// 		sometext: {
		// 			required: true,
		// 			digits: true,
		// 			email: true,
		// 			url: true,
		// 			min: 6,
		// 			max: 6,
		// 			minlength: 6,
		// 			maxlength: 10,
		// 			range: [6,10],
		// 			equalTo: '#id-source'
		// 		}
		// 	}
		// });

		jQuery.validator.addMethod("greaterThan",
			function(value, element, params) {

				if (!/Invalid|NaN/.test(new Date(value))) {
					return new Date(value) > new Date($(params).val());
				}

				return isNaN(value) && isNaN($(params).val())
					|| (Number(value) > Number($(params).val()));
			},'Must be greater than {0}.');

		jQuery.validator.addMethod("lowerThan",
			function(value, element, params) {

				var ele = "#"+params;

				if(params.includes(" ")){
					var tempEle = params.split(" ");
					ele = "#"+tempEle[0].toLowerCase()+tempEle[1];
				}

				return isNaN(value) && isNaN($(ele).val())
					|| (Number(value) <= Number($(ele).val()));
			}, 'Must be lower than {0}.');

		$('#createProgram').validate({
			errorPlacement: errorPlacementInput,
			// Form rules
			rules: {
				programName: {
					required: true
				},
				programValidFrom: {
					required: true
				},
				programValidTo: {
					required: true,
					greaterThan: "#programValidFrom"
				},
				voucherPrice: {
					required: true,
					digits: true
				},
				voucherValue: {
					required: true,
					digits: true,
					min: 5000
				},
				voucherQuantity: {
					required: true,
					digits: true,
					min: 1
				},
				generateVoucher: {
					required: true,
					digits: true,
					min: 1,
					lowerThan: "Voucher Quantity"
				},
				startHour: {
					required: true
				},
				endHour: {
					required: true
				},
				partner: {
					required: true
				}
			}
		});
	}

})();

// Necessary to place dyncamic error messages
// without breaking the expected markup for custom input
function errorPlacementInput(error, element) {
	if( element.parent().parent().is('.mda-input-group') ) {
		error.insertAfter(element.parent().parent()); // insert at the end of group
		element.focus();
	}
	else if( element.parent().is('.mda-form-control') ) {
		error.insertAfter(element.parent()); // insert after .mda-form-control
		element.focus();
	}
	else if( element.parent().is('.input-group') ) {
		error.insertAfter(element.parent()); // insert after .mda-form-control
		element.focus();
	}
	else if ( element.is(':radio') || element.is(':checkbox')) {
		error.insertAfter(element.parent().parent().parent().parent().parent().find(".control-label"));
		$("input[name=partner]").removeClass('error');
		element.focus();
	}
	else {
		error.insertAfter(element);
		element.focus();
	}
}
