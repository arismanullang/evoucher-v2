$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	getPartner();
	getTag();
	getType();
	initForm();
	onChangeElem();

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

	$('#create-program').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			"program-name": {
				required: true
			},
			'program-valid-from': {
				required: true
			},
			'program-valid-to': {
				required: true,
				greaterThan: "#program-valid-from"
			},
			'voucher-price': {
				required: true,
				digits: true
			},
			'voucher-value': {
				required: true,
				digits: true,
				min: 5000
			},
			'voucher-quantity': {
				required: true,
				digits: true,
				min: 1
			},
			'generate-voucher': {
				required: true,
				digits: true,
				min: 1,
				lowerThan: "voucher-quantity"
			},
			'start-hour': {
				required: true
			},
			'end-hour': {
				required: true
			},
			partner: {
				required: true
			}
		}
	});
});

function initForm(){
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

	if ($("#redeem-validity-type").val() == "all") {
		$("#validity-day").attr("style", "display:none");
	} else if ($("#redeem-validity-type").val() == "selected") {
		$("#validity-day").attr("style", "display:block");
	} else {
		$("#validity-day").attr("style", "display:none");
	}
}

function onChangeElem(){
	$("#all-tenant").change(function () {
		var _this = $(this);
		_this.closest('#partner-list').find("input.partner").prop('checked', _this.prop('checked'));
	});
	$("#tag").change(function () {
		var li = $("input[class=partner]:not(:checked)");

		$.each( li, function (i, val) {
			$(val).parent().closest('.col-sm-4').remove();
		});
		getPartnerByTag($(this).val());
	});


// function removeElem(elem){
//   console.log("remove");
//   $(elem).parent().closest('tr').remove();
// }

	$("#voucher-validity-type").change(function () {
		if (this.value == "lifetime") {
			$("#validity-tifetime").attr("style", "display:block");
			$("#validity-date").attr("style", "display:none");
			$("#voucher-valid-from").val("");
			$("#voucher-valid-to").val("");
		} else if (this.value == "period") {
			$("#validity-tifetime").attr("style", "display:none");
			$("#validity-date").attr("style", "display:block");
			$("#voucher-lifetime").val("");
		} else {
			$("#validity-tifetime").attr("style", "display:none");
			$("#validity-date").attr("style", "display:none");
			$("#voucher-valid-from").val("");
			$("#voucher-valid-to").val("");
			$("#voucher-lifetime").val("");
		}
	});
	$("#redeem-validity-type").change(function () {
		if (this.value == "all") {
			$("#validity-day").attr("style", "display:none");
		} else if (this.value == "selected") {
			$("#validity-day").attr("style", "display:block");
		} else {
			$("#validity-day").attr("style", "display:none");
		}
	});
	$("#allow-accumulative").change(function () {
		if (this.checked == true) {
			$("#div-max-redeem-voucher").attr("style", "display:block");
		} else {
			$("#div-max-redeem-voucher").attr("style", "display:none");
		}
	});
	$("#image-url").change(function () {
		readURL(this);
	});
}

function readURL(input) {
	if (input.files && input.files[0]) {
		var reader = new FileReader();
		reader.onload = function (e) {
			$('#image-preview').attr('src', e.target.result);
		}

		reader.readAsDataURL(input.files[0]);
	}
}

function download() {
	window.location = "/v1/ui/sample/link?token="+token;
}

function send() {
	var programName = $("#program-name").val();
	var programType = $("#program-type").find(":selected").val();
	var voucherPrice = parseInt($("#voucher-price").val());
	var maxQuantityVoucher = parseInt($("#voucher-quantity").val());
	var redemptionMethod = $("#redemption-method").find(":selected").val();
	var programValidFrom = $("#program-valid-from").val();
	var programValidTo = $("#program-valid-to").val();
	var startHour = $("#start-hour").val();
	var endHour = $("#end-hour").val();
	var voucherValue = parseInt($("#voucher-value").val());
	var programDescription = $("#program-description").val();

	if(startHour == '00:00'){
		startHour = '00:01';
	}
	if(endHour == '00:00'){
		endHour = '23:59';
	}

	// valid days
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

	if ($("#voucher-validity-type").val() == "period") {
		lifetime = 0;
		periodStart = $("#voucher-valid-from").val();
		periodEnd = $("#voucher-valid-to").val();
	} else if ($("#voucher-validity-type").val() == "lifetime") {
		$( "#voucher-lifetime" ).rules( "add", {
			required: true,
			min: 1,
			max: 1800,
			digits: true
		});

		lifetime = parseInt($("#voucher-lifetime").val());

		periodStart = "01/01/0001";
		periodEnd = "01/01/0001";
	}

	// voucher format
	var voucherFormat = {
		prefix: $("#prefix").val(),
		postfix: "",
		body: "",
		format_type: $("#voucher-format").find(":selected").val(),
		length: 5
	};


	// tnc
	var str = $("#list-rule").summernote('code');
	var tnc = str.replace(/^\s+|\s+$|(\r?\n|\r)/g, '');
	tnc = str.replace('<p>', '');
	tnc = str.replace('</p>', '');
	tnc = '<p>' + tnc + '</p>';

	// max generate and redeem
	var maxGenerate = parseInt($("#generate-voucher").val());
	var maxRedeem = 1;

	// voucher type
	var voucherType = "cash";

	// allow accumulative
	var allowAccumulative = $("#allow-accumulative").is(":checked");
	if (allowAccumulative) {
		maxRedeem = parseInt($("#max-redeem-voucher").val());

		$( "#voucher-lifetime" ).rules( "add", {
			required: true,
			digits: true,
			min: 1
		});
	}

	if(!$("#create-program").valid()) {
		return
	}

	// image
	var formData = new FormData();
	var img = "https://storage.googleapis.com/e-voucher/L1LXN5bpMphnvG6Ce8eUbBSYDW5G3MaH.jpg";
	if ($('#image-url')[0].files[0] != null) {
		$('#modal-loader').modal({backdrop: 'static',keyboard: false}, 'show');
		formData.append('image-url', $('#image-url')[0].files[0]);

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
						$('#modal-loader').modal('hide');
						swal({
								title: 'Success',
								text: 'Program Created',
								type: 'success',
								showCancelButton: false,
								confirmButtonText: 'Ok',
								// closeOnConfirm: false
								closeOnConfirm: true
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
				$('#modal-loader').modal('hide');
				swal({
						title: 'Success',
						text: 'Program Created',
						type: 'success',
						showCancelButton: false,
						confirmButtonText: 'Ok',
						// closeOnConfirm: false
						closeOnConfirm: true
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
				li.appendTo('#partner-list');
			}
		}
	});
}

function getPartnerByTag(param) {
	if(param == "All"){
		getPartner();
		return;
	}
	$.ajax({
		url: '/v1/ui/partner?tag='+param+'&token=' + token,
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
				li.appendTo('#partner-list');
			}
		}
	});
}

function getTag() {
	$.ajax({
		url: '/v1/ui/tag/all',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<option value='"+arrData[i]+"'></option>").html(arrData[i]);
				li.appendTo('#tag');
			}
		}
	});
}

function getType() {
	$.ajax({
		url: '/v1/ui/program/type',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var tempLabel = 'Stock Voucher';
				if(arrData[i] == 'on-demand'){
					tempLabel = 'Mobile Application';
				}else if(arrData[i] == 'bulk'){
					tempLabel = 'Email Blast';
				}else if(arrData[i] == 'gift'){
					tempLabel = 'Gift Voucher';
				}

				var li = $("<option value='"+arrData[i]+"'></option>").html(tempLabel);
				li.appendTo('#program-type');
			}
		}
	});
}

(function () {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$(".select2").select2();
		$('.datepicker-program-from').datepicker({
			container: '#datepicker-program-from',
			autoclose: true,
			startDate: 'd',
			setDate: new Date()
		}).on('changeDate', function (selected) {
			var minDate = new Date(selected.date.valueOf());
			$('.datepicker-program-to').datepicker('setStartDate', minDate);
			$('.datepicker-voucher-from').datepicker('setStartDate', minDate);
		});
		$('.datepicker-program-to').datepicker({
			container: '#datepicker-program-to',
			autoclose: true,
			startDate: '+1d',
			setDate: new Date()
		});
		$('.datepicker-voucher-from').datepicker({
			container: '#datepicker-voucher-from',
			autoclose: true,
			startDate: 'd',
			setDate: new Date()
		}).on('changeDate', function (selected) {
			var minDate = new Date(selected.date.valueOf());
			$('.datepicker-voucher-to').datepicker('setStartDate', minDate);
		});
		$('.datepicker-voucher-to').datepicker({
			container: '#datepicker-voucher-to',
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
	}

})();

// Necessary to place dyncamic error messages
// without breaking the expected markup for custom input
// function errorPlacementInput(error, element) {
// 	if( element.parent().parent().is('.mda-input-group') ) {
// 		error.insertAfter(element.parent().parent()); // insert at the end of group
// 		element.focus();
// 	}
// 	else if( element.parent().is('.mda-form-control') ) {
// 		error.insertAfter(element.parent()); // insert after .mda-form-control
// 		element.focus();
// 	}
// 	else if( element.parent().is('.input-group') ) {
// 		error.insertAfter(element.parent()); // insert after .mda-form-control
// 		element.focus();
// 	}
// 	else if ( element.is(':radio') || element.is(':checkbox')) {
// 		error.insertAfter(element.parent().parent().parent().parent().parent().find(".control-label"));
// 		$("input[name=partner]").removeClass('error');
// 		element.focus();
// 	}
// 	else {
// 		error.insertAfter(element);
// 		element.focus();
// 	}
// }
