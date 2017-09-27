$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	getTag();
});

function getTag() {
	$.ajax({
		url: '/v1/ui/tag/all',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;

			var i;
			for (i = 0; i < arrData.length; i++) {
				var li = $("<option></option>").html(arrData[i]);
				li.appendTo('#tags');
			}
		}
	});
}

function send() {
	$("#createPartner").validate();
	if(!$("#createPartner").valid()){
		$(".error").focus();
		return;
	}

	var listTag = "";
	var li = $("ul.select2-selection__rendered").find("li");
	if (li.length > 0) {
		for (i = 0; i < li.length - 1; i++) {
			var text = li[i].getAttribute("title");

			listTag = listTag + "#" + text;
		}
	}

	var partner = {
		name: $("#partnerName").val(),
		serial_number: $("#serialNumber").val(),
		tag: listTag,
		description: $("#description").val()
	};

	$.ajax({
		url: '/v1/ui/partner/create?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(partner),
		success: function () {
			swal({
					title: 'Success',
					text: 'Partner Created',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/partner/search";
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

(function () {
	'use strict';

	$(runSweetAlert);

	function runSweetAlert() {
		$('.select2').select2({
			tags: true
		});
		$(document).on('click', '.swal-demo4', function (e) {
			e.preventDefault();
			var html;
			if ($("#serialNumber").val() == null) {
				html = 'Do you want create partner ' + $("#partnerName").val() + ' with no serial number?';
			}
			else {
				html = 'Do you want create partner ' + $("#partnerName").val() + ' with serial number ' + $("#serialNumber").val() + '?';
			}

			swal({
					title: 'Are you sure?',
					text: html,
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Yes',
					closeOnConfirm: false
				},
				function () {
					send();
				});

		});

		jQuery.validator.addMethod("greaterThan",
			function(value, element, params) {

				if (!/Invalid|NaN/.test(new Date(value))) {
					return new Date(value) > new Date($(params).val());
				}

				return isNaN(value) && isNaN($(params).val())
					|| (Number(value) > Number($(params).val()));
			},'Must be greater than {0}.');

		$('#createPartner').validate({
			errorPlacement: errorPlacementInput,
			// Form rules
			rules: {
				partnerName: {
					required: true
				},
				serialNumber:{
					digits: true
				}
			}
		});
	}

})();

function errorPlacementInput(error, element) {
	if (element.parent().parent().is('.mda-input-group')) {
		error.insertAfter(element.parent().parent()); // insert at the end of group
	}
	else if (element.parent().is('.mda-form-control')) {
		error.insertAfter(element.parent()); // insert after .mda-form-control
	}
	else if (element.parent().is('.input-group')) {
		error.insertAfter(element.parent()); // insert after .mda-form-control
	}
	else if (element.is(':radio') || element.is(':checkbox')) {
		error.insertAfter(element.parent().parent().parent().parent().parent().find(".control-label"));
	}
	else {
		error.insertAfter(element);
	}
}
