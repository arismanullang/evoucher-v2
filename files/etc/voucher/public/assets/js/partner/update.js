$( document ).ready(function() {
  var id = findGetParameter("id");
  getPartner(id);
});

function getPartner(id) {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/ui/partner?id='+id+"&token="+token,
      type: 'get',
      success: function (data) {
        console.log(data.data);
        var arrData = data.data[0];
        $("#partner-name").html(arrData.name);
        $("#serial-number").val(arrData.serial_number.String);
        $("#description").val(arrData.description.String);
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

  var partner = {
      serial_number: $("#serial-number").val(),
      description: $("#description").val()
    };

    console.log(partner);
    $.ajax({
       url: '/v1/ui/partner/update?id='+id+'&token='+token,
       type: 'post',
       dataType: 'json',
       contentType: "application/json",
       data: JSON.stringify(partner),
       success: function () {
           alert("Partner Updated.");
           window.location = "/partner/search";
       }
   });
}
