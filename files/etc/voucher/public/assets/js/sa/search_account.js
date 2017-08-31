$(document).ready(function () {
	getAccount();
});

function getAccount() {
	$.ajax({
		url: '/v1/ui/sa/account?token=' + token,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			var dataSet = [];
			for (i = 0; i < arrData.length; i++) {
				var button = "<button value='" + arrData[i].id + "' status=" + arrData[i].status + " type='button' class='btn btn-flat btn-sm btn-danger swal-demo-delete'><em class='ion-trash-a'></em></button>";

				var status = "ACTIVE";
				if (arrData[i].status == "deleted") {
					status = "INACTIVE";
				}

				var tempArray = [
					arrData[i].name.toUpperCase()
					, arrData[i].alias.toUpperCase()
					, status
					, button
				];

				dataSet.push(tempArray);
			}

			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'rtip',
				"order": [[1, "desc"]],
				columns: [
					{title: "NAME"},
					{title: "ALIAS"},
					{title: "STATUS"},
					{title: "ACTION"}
				],
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

function detail(id) {
	window.location = "/sa/a-check?id=" + id;
}

function create() {
	var account = {
		name: $("#input-name").val(),
		alias: $("#input-alias").val(),
	};

	$.ajax({
		url: '/v1/ui/sa/a-create?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(account),
		success: function () {
			$("#modal-account").attr("style", "display : none");
			swal('Created!', 'Create success.', window.location.reload());
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function deleteAccount(id) {
	var user = {
		id: id
	};
	$.ajax({
		url: '/v1/ui/sa/a-block?token=' + token,
		type: 'POST',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			setTimeout(
				swal({
						title: 'Blocked!',
						text: 'Block Account Success',
						type: 'success'
					},
					function () {
						location.reload();
					}
				), 1000);
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function activateAccount(id) {
	var user = {
		id: id
	};
	$.ajax({
		url: '/v1/ui/sa/a-activate?token=' + token,
		type: 'POST',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			setTimeout(
				swal({
						title: 'Activated!',
						text: 'Activate Account Success',
						type: 'success'
					},
					function () {
						location.reload();
					}
				), 1000);
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

	//onclick='deleteProgram(\""+arrData[i].Id+"\")'
	function runSweetAlert() {
		$(document).on('click', '.swal-demo-delete', function (e) {
			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want change account\'s status?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Change!',
					closeOnConfirm: false
				},
				function () {
					if (e.target.getAttribute("status") == "created") {
						deleteAccount(e.target.value);
					} else {
						activateAccount(e.target.value);
					}
				});

		});
	}

})();
