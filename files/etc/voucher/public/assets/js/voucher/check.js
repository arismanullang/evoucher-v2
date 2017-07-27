$( document ).ready(function() {
  getVoucher();
});

function getVoucher() {
    console.log("Get Voucher Data");
    var id = findGetParameter('id');
    $.ajax({
        url: '/v1/ui/vouchers/'+id+'?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data.data);
          var data = data.data;

          var date1 = data.valid_at.substring(0, 19).replace("T", " ");
          var date2 = data.expired_at.substring(0, 19).replace("T", " ");

          $("#program-name").html(data.Variant_name);
          $("#voucher-code").html(data.voucher_code);
          $("#voucher-type").html(data.Voucher_type);
          $("#voucher-value").html("Rp " + addDecimalPoints(data.discount_value) + ",00");
          $("#reference-no").html(data.reference_no);
          $("#period").html(date1 + "</br></br>To</br></br>" + date2)

          var email = data.holder_email;
          if(data.holder_email == ""){
              email = "Unknown";
          }
          var phone = data.holder_phone;
          if(data.holder_phone == ""){
              phone = "Unknown";
          }

          $("#holder-name").html(toTitleCase(data.holder));
          $("#holder-email").html(email);
          $("#holder-phone").html(phone);

          var dateCreated = data.created_at.substring(0, 19).replace("T", " ");
          $("#issued-state").html(dateCreated);
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
                    text: 'Do you want delete program?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Yes, delete it!',
                    closeOnConfirm: false
                },
                function() {
                    swal('Deleted!', 'Delete success.', 'delete');
                });

        });
    }

})();
