var id = findGetParameter('id')
$( document ).ready(function() {
  getVoucher(id);
  getVariant(id);

  $('#profileForm').submit(function(e) {
       e.preventDefault();
       e.returnValue = false;
  });
});

function getVoucher(id) {
    console.log("Get Voucher Data");

    var arrData = [];
    $.ajax({
        url: '/v1/voucher?variant_id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.length;
          if (arrData.length > 4){
            limit = 4;
            $("<div class='card-body pv0 text-right'><a href='#' class='btn btn-flat btn-info'>View all</a></div>").appendTo('#listVoucher');
          }

          for ( i = 0; i < limit; i++){
            var html = "<div class='mda-list-item-icon'><em class='ion-pricetag icon-2x'></em></div>"
            +  "<div class='mda-list-item-text'>"
            +  "<h3><a href='#'>"+arrData[i].voucher_code+"</a></h3>"
            +  "<div class='text-muted text-ellipsis'>Status "+arrData[i].state+"</div>"
            +"</div>";
            var li = $("<div class='mda-list-item'></div>").html(html);
            li.appendTo('#listVoucher');
          }
        }
    });
}

function getVariant(id) {
    console.log("Get Variant Data");

    var arrData = [];
    $.ajax({
        url: '/v1/api/get/variant/'+id+'?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var result = data.data;

          var startDate = result.StartDate.substr(0,10);
          var endDate = result.EndDate.substr(0,10);
          var period = startDate + " to " + endDate;

          $('#variantName').html(result.VariantName);
          $('#variantDescription').html(result.VariantDescription);
          $('#variantType').html(result.VariantType);
          $('#voucherType').html(result.VoucherType);
          $('#conversionRate').html(result.VoucherPrice);
          $('#maxQuantityVoucher').html(result.MaxQuantityVoucher);
          $('#voucherValue').html(result.DiscountValue);
          $('#period').html(period);
          $('#variantTnc').html(result.VariantTnc);
          $('#remainingVoucher').html(result.MaxQuantityVoucher - result.Voucher.length);

          var i;
          var arrData = data.data.ValidPartners;
          var limit = arrData.length;
          if (arrData.length > 4){
            limit = 4;
            $("<div class='card-body pv0 text-right'><a href='#' class='btn btn-flat btn-info'>View all</a></div>").appendTo('#listPartner');
          }

          for ( i = 0; i < limit; i++){
            var html = "<div class='mda-list-item-icon'><em class='ion-ios-person icon-2x'></em></div>"
            +  "<div class='mda-list-item-text'>"
            +  "<p>"+arrData[i]+"</p>"
            +"</div>";
            var li = $("<div class='mda-list-item'></div>").html(html);
            li.appendTo('#listPartner');
          }


        }
    });
}

function editPartner(){
  window.location = "/variant/update?id="+id;
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
