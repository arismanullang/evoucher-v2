$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	jQuery.validator.addMethod("greaterThan",
		function(value, element, params) {

			if (!/Invalid|NaN/.test(new Date(value))) {
				return new Date(value) > new Date($(params).val());
			}

			return isNaN(value) && isNaN($(params).val())
				|| (Number(value) > Number($(params).val()));
		},'Must be greater than {0}.');

	$('#create-bank-account').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			'company-name': {
				required: true
			},
			'company-pic':{
				required: true
			},
			'company-telp': {
				required: true
			},
			'company-email':{
				required: true
			},
			'bank-name':{
				required: true
			},
			'bank-branch':{
				required: true
			},
			'bank-account-holder':{
				required: true
			},
			'bank-account-number':{
				digits: true
			}
		}
	});
});

function send() {
	if(!$("#create-bank-account").valid()){
		$(".error").focus();
		return;
	}

	var partner = {
		company_name: $("#company-name").val(),
		company_pic: $("#company-pic").val(),
		company_telp: $("#company-telp").val(),
		company_email: $("#company-email").val(),
		bank_name: $("#bank-name").val(),
		bank_branch: $("#bank-branch").val(),
		bank_account_number: $("#bank-account-number").val(),
		bank_account_holder: $("#bank-account-holder").val()
	};

	$.ajax({
		url: '/v1/ui/bank_account/create?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(partner),
		success: function () {
			swal({
					title: 'Success',
					text: 'Bank Account Created',
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
		$(document).on('click', '.swal-demo4', function (e) {
			e.preventDefault();
			var html;
			html = 'Do you want create bank account ' + $("#bank-account-number").val() + ' on behalf of ' + $("#bank-account-holder").val() + '?';

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
	}
})();
