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
            var html;
            if($("#serial-number").val() == null){
              html = 'Do you want create partner '+$("#partner-name").val()+' with no serial number?';
            }
            else{
              html = 'Do you want create partner '+$("#partner-name").val()+' with serial number '+$("#serial-number").val()+'?';
            }

            swal({
                    title: 'Are you sure?',
                    text: html,
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Yes',
                    closeOnConfirm: false
                },
                function() {
                  error = false;
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

                  swal('Success', 'Partner '+$("#partner-name").val()+' created.', send());
                });

        });
    }

})();
