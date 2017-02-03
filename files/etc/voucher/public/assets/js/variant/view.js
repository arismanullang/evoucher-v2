$( window ).load(function() {
  searchByUser();
});

function searchByUser() {
    var request = {
        user: "IzKyd9yX" //$("#variantName").val()
      };

    $.ajax({
        url: 'http://127.0.0.1:8080/variant/print/',
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify(request),
        success: function (data){
          alert(data);
        }
    });
}
