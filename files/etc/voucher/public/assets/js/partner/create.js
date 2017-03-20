function send() {
  var partner = {
      partner_name: $("#partnerName").val(),
      serial_number: $("#serialNumber").val(),
    };

    console.log(partner);
    $.ajax({
       url: '/v1/create/partner?token='+token,
       type: 'post',
       dataType: 'json',
       contentType: "application/json",
       data: JSON.stringify(partner),
       success: function () {
           alert("Partner created.");
       }
   });
}
