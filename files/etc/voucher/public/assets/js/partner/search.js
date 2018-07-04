$( document ).ready(function() {
  getPartner();
  $(document).ajaxStart(function(){
    // Show image container
    $(".cssload-loader").show();
   });
   $(document).ajaxComplete(function(){
    // Hide image container
    $(".cssload-loader").hide();
   });
});

function getPartner() {
    $.ajax({
      url: '/v1/ui/partner/all?token='+token,
      type: 'get',
      success: function (data) {
        var arrData = [];
        arrData = data.data;

        var i;
        var dataSet = [];
        for (i = 0; i < arrData.length; i++){
	  var button = "<button type='button' onclick='detail(\""+arrData[i].id+"\")' class='btn btn-flat btn-sm btn-info'><em class='ion-search'></em></button>"+
		  "<button type='button' class='btn btn-flat btn-sm btn-info' onclick='edit(\""+arrData[i].id+"\")'><em class='ion-edit'></em></button>"+
		"<button value='"+arrData[i].id+"' type='button' class='btn btn-flat btn-sm btn-danger swal-demo4'><em class='ion-trash-a'></em></button>";

	  var tempArray = [
		arrData[i].name
		, arrData[i].company_name + ", " + arrData[i].bank_name + " - " + arrData[i].bank_account_number
		, arrData[i].tag.String
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
	      "order": [[ 1, "desc" ]],
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
	      .keyup(function() {
		      table.fnFilter(this.value, columnInputs.index(this));
	      });
      	}
  });
}

function detail(url){
	window.location = "/partner/check?id="+url;
}

function edit(url){
  window.location = "/partner/update?id="+url;
}

function addPartner() {
  window.location = "/partner/create";
}

function deletePartner(id) {
	$.ajax({
		url: '/v1/ui/partner/delete?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			swal('Deleted!', 'Delete success.');

			setTimeout(function () {
				getPartner();
			}, 1000);
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

(function() {
    'use strict';

    $(runSweetAlert);
    //onclick='deleteProgram(\""+arrData[i].Id+"\")'
    function runSweetAlert() {
        $(document).on('click', '.swal-demo4', function(e) {
            e.preventDefault();
            swal({
                    title: 'Are you sure?',
                    text: 'Do you want delete partner?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Yes, delete it!',
                    closeOnConfirm: false
                },
                function() {
					deletePartner(e.target.value);
                });

        });
    }

})();
