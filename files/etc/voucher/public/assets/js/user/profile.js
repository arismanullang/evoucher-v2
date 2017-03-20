var id = findGetParameter('id')
$( document ).ready(function() {
  getUser(id);
  getAccount(id);

  $('#profileForm').submit(function(e) {
       e.preventDefault();
       e.returnValue = false;
  });
});

function getUser(id) {
    console.log("Get Voucher Data");

    var arrData = [];
    $.ajax({
        url: '/v1/vouchers?variant_id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.length;
          if (arrData.length > 4){
            limit = 4;
            $("<div class='card-body pv0 text-right'><a href='/voucher/search?variant_id="+id+"' class='btn btn-flat btn-info'>View all</a></div>").appendTo('#cardVoucher');
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
        },
        error: function (data) {
          console.log(data.data);
          $("<div class='card-body text-center'>No Voucher Yet</div>").appendTo('#cardVoucher');
        }
    });
}

function getPartner(id) {
    console.log("Get Voucher Data");

    var arrData = [];
    $.ajax({
        url: '/v1/api/get/partner?variant_id='+id+'&token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var i;
          var arrData = data.data;
          var limit = arrData.length;

          for ( i = 0; i < limit; i++){
            var html = "<div class='mda-list-item-icon'><em class='ion-ios-person icon-2x'></em></div>"
            +  "<div class='mda-list-item-text'>"
            +  "<h3><a href='#'>"+arrData[i].partner_name+"</a></h3>"
            +  "<div class='text-muted text-ellipsis'>"+arrData[i].serial_number+"</div>"
            +"</div>";
            var li = $("<div class='mda-list-item'></div>").html(html);
            li.appendTo('#listPartner');
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

          var remainingVoucher = result.MaxQuantityVoucher;
          if( result.Voucher != null)
            remainingVoucher = esult.MaxQuantityVoucher - result.Voucher.length;

          $('#variantName').html(result.VariantName);
          $('#variantDescription').html(result.VariantDescription);
          $('#variantType').html(result.VariantType);
          $('#voucherType').html(result.VoucherType);
          $('#conversionRate').html(result.VoucherPrice);
          $('#maxQuantityVoucher').html(result.MaxQuantityVoucher);
          $('#voucherValue').html(result.DiscountValue);
          $('#period').html(period);
          $('#variantTnc').html(result.VariantTnc);
          $('#remainingVoucher').html(remainingVoucher);
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
