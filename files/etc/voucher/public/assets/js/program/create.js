$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	getPartner();

	$("#token").val(token);

	$("#all-tenant").change(function () {
		var lis = $("input[class=partner]");
		if ($("#all-tenant").is(':checked')) {
			for (var i = 0; i < lis.length; i++) {
				lis[i].checked = true;
			}
		} else {
			for (var i = 0; i < lis.length; i++) {
				lis[i].checked = false;
			}
		}
	});

	$("#voucher-validity-type").change(function () {
		if (this.value == "lifetime") {
			$("#validity-lifetime").attr("style", "display:block");
			$("#validity-date").attr("style", "display:none");
			$("#voucher-valid-from").val("");
			$("#voucher-valid-to").val("");
		} else if (this.value == "period") {
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
	$("#program-type").change(function () {
		if (this.value == "bulk") {
			$("#target").attr("style", "display:block");
			$("#conversion-row").attr("style", "display:none");
			$("#voucher-price").val(0);
			$("#voucher-price").attr("disabled", true);
		} else {
			$("#target").attr("style", "display:none");
			$("#conversion-row").attr("style", "display:block");
			$("#voucher-price").val("");
			$("#voucher-price").attr("disabled", false);
		}
	});
	$("#allow-accumulative").change(function () {
		if (this.checked == true) {
			$("#use-voucher").attr("style", "display:block");
		} else {
			$("#use-voucher").attr("style", "display:none");
		}
	});
	$("#image-url").change(function () {
		readURL(this);
	});

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

	if ($("#program-type").val() == "bulk") {
		$("#target").attr("style", "display:block");
		$("#conversion-row").attr("style", "display:none");
		$("#voucher-price").val(0);
		$("#voucher-price").attr("disabled", true);
	} else {
		$("#target").attr("style", "display:none");
		$("#conversion-row").attr("style", "display:block");
		$("#voucher-price").attr("disabled", false);
	}
	if ($("#allow-accumulative").is(":checked")) {
		$("#use-voucher").attr("style", "display:block");
	} else {
		$("#use-voucher").attr("style", "display:none");
	}
});

function readURL(input) {
	if (input.files && input.files[0]) {
		var reader = new FileReader();
		reader.onload = function (e) {
			$('#image-preview').attr('src', e.target.result);
		}

		reader.readAsDataURL(input.files[0]);
	}
}

function send() {
	error = false;
	var errorMessage = "";
	var i;

	// valid days
	var listDay = "";
	if ($("#redeem-validity-type").val() == "all") {
		listDay = "all";
	} else if ($("#redeem-validity-type").val() == "selected") {
		var li = $("ul.select2-selection__rendered").find("li");

		if (li.length == 0 || parseInt($("#length").val()) < 8) {
			error = true;
			errorMessage = "Valid days. ";
		}

		for (i = 0; i < li.length - 1; i++) {
			var text = li[i].getAttribute("title");
			var value = $("option").filter(function () {
				return $(this).text() === text;
			}).first().attr("value");

			listDay = listDay + value + ";";
		}
	}

	// partner
	var listPartner = [];
	var li = $("input[class=partner]:checked");

	if (li.length == 0 || parseInt($("#length").val()) < 8) {
		error = true;
		errorMessage = "Valid partners. ";
	}

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
		lifetime = $("#voucher-lifetime").val();

		if (lifetime == "" || lifetime == "0") {
			errorMessage = "Lifetime can't be empty or 0. "
			error = true;
		}
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
	}


	// tnc
	var str = $("#list-rule").summernote('code');
	var tnc = str.replace(/^\s+|\s+$|(\r?\n|\r)/g, '');
	tnc = str.replace('<p>', '');
	tnc = str.replace('</p>', '');
	tnc = '<p>' + tnc + '</p>';

	// max generate and redeem
	var maxGenerate = parseInt($("#max-generate-voucher").val());
	var maxRedeem = parseInt($("#max-redeem-voucher").val());

	// voucher type
	var voucherType = "cash";

	// allow accumulative
	var allowAccumulative = $("#allow-accumulative").is(":checked");
	if (allowAccumulative) {
		maxRedeem = 1;
	}

	$('input[check="true"]').each(function () {
		if ($("#program-type").val() == "bulk") {
			if (this.getAttribute("id") == "max-quantity-voucher" || this.getAttribute("id") == "max-generate-voucher" || this.getAttribute("id") == "voucher-price" || this.getAttribute("id") == "max-redeem-voucher") {
				allowAccumulative = false;
				maxGenerate = 1;
				maxRedeem = 1;
				return true;
			}
		}

		if ($(this).val() == "") {
			$(this).addClass("error");
			$(this).parent().closest('div').addClass("input-error");
			error = true;
			errorMessage = "Empty. ";
		}

		if (this.getAttribute("id") == "max-quantity-voucher" || this.getAttribute("id") == "max-generate-voucher" || this.getAttribute("id") == "voucher-price" || this.getAttribute("id") == "max-redeem-voucher") {
			if ($(this).val() < 1) {
				$(this).addClass("error");
				$(this).parent().closest('div').addClass("input-error");
				error = true;
				errorMessage = "0 Value. ";
			}
		}

		if ($(this).attr("id") == "length") {
			if (parseInt($(this).val()) < 8) {
				error = true;
				errorMessage = "Id must be 8 digit. ";
			}
		}
	});

	if (error) {
		$("#error-message").html(errorMessage + "Please check your input.");
		return
	}

	// image
	var formData = new FormData();
	var img = "https://storage.googleapis.com/e-voucher/Nd3QxH8El2Zuy12QhXs5Y305vPL4VZJJ.jpg";
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
					name: $("#program-name").val(),
					type: $("#program-type").find(":selected").val(),
					voucher_format: voucherFormat,
					voucher_type: voucherType,
					voucher_price: parseInt($("#voucher-price").val()),
					max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
					max_redeem_voucher: maxRedeem,
					max_generate_voucher: maxGenerate,
					allow_accumulative: allowAccumulative,
					redemption_method: $("#redemption-method").find(":selected").val(),
					start_date: $("#program-valid-from").val(),
					end_date: $("#program-valid-to").val(),
					start_hour: $("#start-hour").val(),
					end_hour: $("#end-hour").val(),
					voucher_value: parseInt($("#voucher-value").val()),
					image_url: img,
					tnc: tnc,
					description: $("#program-description").val(),
					validity_days: listDay,
					valid_partners: listPartner,
					valid_voucher_start: periodStart,
					valid_voucher_end: periodEnd,
					voucher_lifetime: parseInt(lifetime)
				};

				$.ajax({
					url: '/v1/ui/program/create?token=' + token,
					type: 'post',
					dataType: 'json',
					contentType: "application/json",
					data: JSON.stringify(program),
					success: function (data) {
						if ($("#program-type").find(":selected").val() == "bulk") {

							var targets = new FormData();
							targets.append('list-target', $("#list-target")[0].files[0]);

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
			name: $("#program-name").val(),
			type: $("#program-type").find(":selected").val(),
			voucher_format: voucherFormat,
			voucher_type: voucherType,
			voucher_price: parseInt($("#voucher-price").val()),
			max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
			max_redeem_voucher: maxRedeem,
			max_generate_voucher: maxGenerate,
			allow_accumulative: allowAccumulative,
			redemption_method: $("#redemption-method").find(":selected").val(),
			start_date: $("#program-valid-from").val(),
			end_date: $("#program-valid-to").val(),
			start_hour: $("#start-hour").val(),
			end_hour: $("#end-hour").val(),
			voucher_value: parseInt($("#voucher-value").val()),
			image_url: img,
			tnc: tnc,
			description: $("#program-description").val(),
			validity_days: listDay,
			valid_partners: listPartner,
			valid_voucher_start: periodStart,
			valid_voucher_end: periodEnd,
			voucher_lifetime: parseInt(lifetime)
		};

		$.ajax({
			url: '/v1/ui/program/create?token=' + token,
			type: 'post',
			dataType: 'json',
			contentType: "application/json",
			data: JSON.stringify(program),
			success: function (data) {
				if ($("#program-type").find(":selected").val() == "bulk") {

					var targets = new FormData();
					targets.append('list-target', $("#list-target")[0].files[0]);

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
					+ "<input type='checkbox' class='partner' value='" + arrData[i].id + "'>"
					+ "<span class='ion-checkmark-round'></span>" + arrData[i].name
					+ "</label>";
				li.html(html);
				li.appendTo('#partner-list');
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
		$('.select2').select2();
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
	}

})();
