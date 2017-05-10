$( document ).ready(function() {
  getPartner();
});

function getPartner() {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/get/tag',
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;
        console.log(arrData);
        var i;
        for (i = 0; i < arrData.length; i++){
          var li = $("<li></li>").html(arrData[i]);
          li.appendTo('#listTag');
        }

        if ($.fn.DataTable.isDataTable("#datatable1")) {
          $('#datatable1').DataTable().clear().destroy();
        }

        var table = $('#datatable1').dataTable({
            data: dataSet,
            dom: 'lBrtip',
            buttons: [
                'copy', 'csv', 'excel', 'pdf', 'print'
            ],
            "order": [[ 4, "asc" ]],
            columns: [
                { title: "Program Name" },
                { title: "Action"}
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

          columnInputs.keyup(function() {
              table.fnFilter(this.value, columnInputs.index(this));
          });
      }
  });
}

function add(param) {

  var tag = {
    tag: param
  };

  $.ajax({
    url: '/v1/create/tag?token='+token,
    type: 'post',
    dataType: 'json',
    contentType: "application/json",
    data: JSON.stringify(tag),
    success: function (data) {
      location.reload();
    }
  });
}

(function() {
    'use strict';

    $(runSweetAlert);
    //onclick='deleteVariant(\""+arrData[i].Id+"\")'
    function runSweetAlert() {
        $(document).on('click', '.swal-demo4', function(e) {
            e.preventDefault();
            console.log(e.target.value);
            swal({
                    title: 'Are you sure?',
                    text: 'Do you want insert a new tag "'+$("#tag-value").val()+'"?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Insert',
                    closeOnConfirm: false
                },
                function() {
                    swal('Success!', 'Add success.', add($("#tag-value").val()));
                });

        });
    }

})();
