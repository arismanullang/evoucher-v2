$( document ).ready(function() {
  getTag();
});

function getTag() {
    console.log("Get Tag List");

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
          var li = $("<option></option>").html(arrData[i]);
          li.appendTo('#tags');
        }
      }
  });
}

function send() {

  var listTag = "";
  var li = $( "ul.select2-selection__rendered" ).find( "li" );
  if(li.length == 0 || parseInt($("#length").val()) < 8){
    error = true;
  }

  for (i = 0; i < li.length-1; i++) {
    var text = li[i].getAttribute("title");

    listTag = listTag+"#"+text;
  }

  var partner = {
    partner_name: $("#partner-name").val(),
    serial_number: $("#serial-number").val(),
    tag: listTag,
    description: $("#description").val(),
  };

  console.log(partner);
  $.ajax({
     url: '/v1/create/partner?token='+token,
     type: 'post',
     dataType: 'json',
     contentType: "application/json",
     data: JSON.stringify(partner),
     success: function () {
         window.location = "/partner/search?token="+token;
     }
 });
}

(function() {
    'use strict';

    $(runSweetAlert);
    //onclick='deleteVariant(\""+arrData[i].Id+"\")'
    function runSweetAlert() {
    	$('.select2').select2();
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
