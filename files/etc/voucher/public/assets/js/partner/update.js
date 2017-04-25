$( document ).ready(function() {
  var id = findGetParameter("id");
  getPartner(id);
});

function getPartner(id) {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/get/partner/'+id+"?token="+token,
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;
        console.log(arrData);
        $("#partner-name").html(arrData.partner_name);
        $("#serial-number").val(arrData.serial_number.String);
      }
  });
}

function update() {
  var i;

  var id = findGetParameter("id");
  var error = false;
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

  var userReq = {
      serial_number: $("#serial-number").val()
    };

    console.log(userReq);
    $.ajax({
       url: '/v1/update/partner/'+id+'?token='+token,
       type: 'post',
       dataType: 'json',
       contentType: "application/json",
       data: JSON.stringify(userReq),
       success: function () {
           alert("Partner Updated.");
           window.location = "/partner/search";
       }
   });
}
