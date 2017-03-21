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
          var html = "<div class='mda-list-item-icon'><em class='ion-home icon-2x'></em></div>"
          +  "<div class='mda-list-item-text'>"
          +  "<h3>"+arrData[i].PartnerName+"</h3>"
          +  "<p class='text-muted'> Serial Number : "+arrData[i].SerialNumber.String+"</p>"
          +"</div>";
          var li = $("<div class='mda-list-item'></div>").html(html);
          li.appendTo('#listPartner');
        }
      }
  });
}

function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    location.search
    .substr(1)
        .split("&")
        .forEach(function (item) {
        tmp = item.split("=");
        if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
    });
    return result;
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
                    text: 'Do you want delete variant?',
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
