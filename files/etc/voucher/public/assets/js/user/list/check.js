$(document).ready(function () {
	var id = findGetParameter('id');
	$('#list-id').val(id);
	getUser(id);
});

function getUser(id) {
	$.ajax({
		url: '/v1/ui/user/list?token=' + token + '&id=' + id,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			var dataSet = [];
			$("#list-name").html(arrData.name);
			for (i = 0; i < arrData.email_users.length; i++) {
				var button = "<button value='" + arrData.email_users[i].id + "' type='button' class='btn btn-flat btn-sm btn-danger swal-demo-delete'><em class='ion-trash-a'></em></button>";
				dataSet[i] = [
					arrData.email_users[i].name
					, arrData.email_users[i].email
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

function addNew() {
	window.location = "/user/list/add-new?id="+$('#list-id').val();;
}

function addExisting() {
	window.location = "/user/list/add-exist?id="+$('#list-id').val();;
}

function deleteUser(id) {
	var user = {
		email_user_id: id,
		list_id: $("#list-id").val()
	};
	$.ajax({
		url: '/v1/ui/user/list/remove?token=' + token,
		type: 'POST',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			window.location.reload();
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
					text: 'Do you want remove user from list?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Remove',
					closeOnConfirm: true
				},
				function () {
						deleteUser(e.target.value);
				});

		});
	}
})();
