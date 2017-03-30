function send() {
  var partner = {
    partner_name: $("#partner-name").val(),
    serial_number: $("#serial-number").val(),
  };

  console.log(partner);
  $.ajax({
     url: '/v1/create/partner?token='+token,
     type: 'post',
     dataType: 'json',
     contentType: "application/json",
     data: JSON.stringify(partner),
     success: function () {
         window.location = "/partner/search";
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
                    text: 'Do you want create partner '+$("#partnerName").val()+' with serial number '+$("#serialNumber").val()+'?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Yes',
                    closeOnConfirm: false
                },
                function() {
                  $('input[check="true"]').each(function() {
                    if($(this).val() == ""){
                      $(this).addClass("error");
                      $(this).parent().closest('div').addClass("input-error");
                      error = true;
                    }
                  });

                  if(error){
                    alert("Please check your input.");
                    return
                  }

                  swal('Success', 'Partner '+$("#partnerName").val()+' created.', send());
                });

        });
    }

})();
