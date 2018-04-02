$(document).ready(function () {
	getUser();
});

function getUser() {
	$.ajax({
		url: '/v1/ui/user/email/all?token=' + token,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			var dataSet = [];
			for (i = 0; i < arrData.length; i++) {
				var button = "<button value='" + arrData[i].id + "' type='button' class='btn btn-flat btn-sm btn-danger swal-demo-delete'><em class='ion-trash-a'></em></button>";
				dataSet[i] = [
					arrData[i].name
					, arrData[i].email
					, button
				];
			}

			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'rtip',
				"order": [[0, "asc"]],
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

			columnInputs
				.keyup(function () {
					table.fnFilter(this.value, columnInputs.index(this));
				});
		}
	});
}

function addUser() {
	window.location = "/user/email/register";
}

function deleteUser(id) {
	var user = {
		id: id
	};
	$.ajax({
		url: '/v1/ui/user/email/delete?token=' + token,
		type: 'POST',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			swal({
				title: 'Success',
				text: 'Delete Success',
				type: 'success',
				showCancelButton: false,
				confirmButtonText: 'Ok',
				closeOnConfirm: false
			},
			function() {
				window.location.reload();
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
		$(document).on('click', '.swal-demo-delete', function (e) {
			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want delete this user?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Delete',
					closeOnConfirm: false
				},
				function () {
						deleteUser(e.target.value);
				});

		});
		$(document).on('click', '.swal-demo-reset', function (e) {
			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want reset user\'s password?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Reset!',
					closeOnConfirm: false
				},
				function () {
					resetPassword(e.target.value);
				});

		});
	}

})();
