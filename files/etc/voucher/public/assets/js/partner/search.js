$( document ).ready(function() {
  getPartner();
});

function getPartner() {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/get/partner',
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;
        console.log(arrData);
        var i;
        for (i = 0; i < arrData.length; i++){
          var html = "<div class='mda-list-item-icon bg-info'><em class='ion-home icon-2x'></em></div>"
          + "<div class='mda-list-item-text'>"
          + "<h3>"+arrData[i].partner_name+"</h3>"
          + "<h4 class='text-muted'> Serial Number : "+arrData[i].serial_number.String+"</h4>"
          + "<h4 class='text-muted'> Serial Number : "+arrData[i].serial_number.String+"</h4>"
          + "</div>"
          + "<div class='pull-right dropdown dropdown-partner'>"
          + "<button type='button' data-toggle='dropdown' class='btn btn-default btn-flat btn-flat-icon'><em class='ion-android-more-vertical'></em></button>"
          + "<ul role='menu' class='dropdown-menu dropdown-menu-partner md-dropdown-menu dropdown-menu-right'>"
          + "<li><button type='button' class='btn btn-flat btn-sm btn-info'><em class='ion-edit'></em> Edit</button>"
          + "<li><button type='button' class='btn btn-flat btn-sm btn-danger swal-demo4'><em class='ion-trash-a'></em> Delete</button></li>"
          + "</ul>"
          + "</div>";
          var li = $("<div class='mda-list-item col-sm-6'></div>").html(html);
          li.appendTo('#listPartner');
        }
      }
  });
}

function edit(url){
  window.location = "/partner/update?id="+url;
}

function addPartner() {
  window.location = "/partner/create";
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
                    text: 'Do you want insert new partner?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Yes, delete it!',
                    closeOnConfirm: false
                },
                function() {
                    swal('Deleted!', 'Delete success.', deleteVariant(e.target.value));
                });

        });
    }

})();
