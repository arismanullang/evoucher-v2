$( document ).ready(function() {
  getFeature();
});

function getFeature() {
    $.ajax({
      url: '/v1/ui/feature/all',
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;
        console.log(arrData);
        var i;
        var uiFeature = [];
        var apiFeature = [];
        for (i = 0; i < arrData.length; i++){
	  var html = "<div class='card'><div class='card-body pt0 pb0'><div class='row'>"
		+ "<div class='col-sm-12'>"
		+ "<div class='checkbox c-checkbox'>"
		+ "<label class='text-thin font-size-12px'>"
		+ "<input name='agreements' value='"+arrData[i].id+"' type='checkbox'><span class='ion-checkmark-round'></span>"+arrData[i].category+"-"+arrData[i].detail
		+ "</label>"
		+ "</div>"
		+ "</div>"
		+ "</div></div></div>";
          var li = $("<div class='col-md-6'></div>").html(html);

	  if(arrData[i].type == 'ui'){
		  li.appendTo('#list-ui-feature');
	  }else{
		  li.appendTo('#list-api-feature');
	  }
        }
      }
  });
}

function add(param) {

  var tag = {
    tag: param
  };

  $.ajax({
    url: '/v1/ui/tag/create?token='+token,
    type: 'post',
    dataType: 'json',
    contentType: "application/json",
    data: JSON.stringify(tag),
    success: function (data) {
      $("#tag-value").val("");
      location.reload();
    }
  });
}

function deleteTag(param) {
  var tag = {
    tag: param
  };

  $.ajax({
    url: '/v1/ui/tag/delete?id='+param+'&token='+token,
    type: 'get',
    success: function (data) {
      location.reload();
    }
  });
}

function deleteTagBulk(param) {

  var tag = {
    tags: param
  };

  $.ajax({
    url: '/v1/ui/tag/delete?token='+token,
    type: 'post',
    dataType: 'json',
    contentType: "application/json",
    data: JSON.stringify(tag),
    success: function (data) {
      location.reload();
    }
  });
}

(function() {
    'use strict';

    $(runSweetAlert);
    function runSweetAlert() {
	$(document).on('click', '.swal-demo4', function(e) {
            	e.preventDefault();
            	console.log(e.target.value);
            	swal({
                    title: 'Are you sure?',
                    text: 'Do you want insert a new tag "'+$("#tag-value").val()+'"?',
                    type: 'warning',
                    showCancelButton: true,
                    confirmButtonColor: '#DD6B55',
                    confirmButtonText: 'Insert',
                    closeOnConfirm: false
                },
                function() {
                    swal('Success!', 'Add success.', add($("#tag-value").val()));
                });
	});
    	$(document).on('click', '.swal-demo2', function(e) {
                	e.preventDefault();
                	console.log(e.target.value);
                	swal({
                        title: 'Are you sure?',
                        text: 'Do you want delete tag "'+e.target.value+'"?',
                        type: 'warning',
                        showCancelButton: true,
                        confirmButtonColor: '#DD6B55',
                        confirmButtonText: 'Delete',
                        closeOnConfirm: false
                    },
                    function() {
                        swal('Success!', 'Delete success.', deleteTag(e.target.value));
                    });
    	});
    	$(document).on('click', '.swal-demo3', function(e) {
		var checkbox = $("input[type=checkbox]:checked");
		var data = [];

		for( var i = 0; i < checkbox.length; i++){
			data[i] = checkbox[i].value;
		}

        	e.preventDefault();
        	console.log(data);
        	swal({
                title: 'Are you sure?',
                text: 'Do you want delete all these tags?',
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#DD6B55',
                confirmButtonText: 'Delete',
                closeOnConfirm: false
            },
            function() {
                swal('Success!', 'Delete success.', deleteTagBulk(data));
            });
    	});
    }

})();
