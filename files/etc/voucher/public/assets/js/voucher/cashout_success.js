$( document ).ready(function() {
	var id = findGetParameter('id');
	$('#cashout-id').val(id);
});

function print(){
	window.location = "/voucher/print?id="+$('#cashout-id').val();
}

function next(){
	swal({
			title: 'Are you already print the invoice?',
			text: 'You will not be able to recover the last details',
			type: 'warning',
			showCancelButton: true,
			confirmButtonColor: '#4CAF50',
			confirmButtonText: 'Yes',
			closeOnConfirm: true
		},
		function() {
			window.location = "/voucher/cashout";
		});
}
